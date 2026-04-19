package sms

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// KavenegarSender sends OTP messages via the Kavenegar REST API.
// API key is loaded from environment and never logged or exposed in responses (T-03-04).
type KavenegarSender struct {
	apiKey   string
	template string
	client   *http.Client
}

// NewKavenegarSender creates a new KavenegarSender with the given API key and template name.
func NewKavenegarSender(apiKey, template string) *KavenegarSender {
	return &KavenegarSender{
		apiKey:   apiKey,
		template: template,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// SendOTP sends an OTP code to the given phone number via Kavenegar's verify/lookup API.
func (k *KavenegarSender) SendOTP(phone, code string) error {
	url := fmt.Sprintf(
		"https://api.kavenegar.com/v1/%s/verify/lookup.json?receptor=%s&token=%s&template=%s",
		k.apiKey, phone, code, k.template,
	)

	resp, err := k.client.Get(url)
	if err != nil {
		return fmt.Errorf("kavenegar request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("kavenegar returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
