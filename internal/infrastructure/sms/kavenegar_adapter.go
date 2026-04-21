package sms

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

// KavenegarAdapter implements SMSProvider using the Kavenegar REST API.
// Docs: https://app.kavenegar.com/docs
type KavenegarAdapter struct {
	apiKey      string
	otpTemplate string
	httpClient  *http.Client
}

func NewKavenegarAdapter(apiKey, otpTemplate string) *KavenegarAdapter {
	return &KavenegarAdapter{
		apiKey:      apiKey,
		otpTemplate: otpTemplate,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendOTP sends an OTP SMS via Kavenegar VerifyLookup API.
func (k *KavenegarAdapter) SendOTP(ctx context.Context, mobile, otp string) error {
	endpoint := fmt.Sprintf(
		"https://api.kavenegar.com/v1/%s/verify/lookup.json",
		k.apiKey,
	)

	params := url.Values{}
	params.Set("receptor", mobile)
	params.Set("template", k.otpTemplate)
	params.Set("token", otp)
	params.Set("type", "sms")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create kavenegar request: %w", err)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := k.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("kavenegar request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(body)).
			Str("mobile", mobile).
			Msg("kavenegar send failed")
		return fmt.Errorf("kavenegar returned status %d", resp.StatusCode)
	}

	log.Debug().
		Str("mobile", mobile).
		Msg("OTP SMS sent via Kavenegar")

	return nil
}
