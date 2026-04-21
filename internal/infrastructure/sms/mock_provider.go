package sms

import (
	"context"

	"github.com/rs/zerolog/log"
)

// MockSMSProvider is used in development/testing — logs OTP to stdout instead of sending SMS.
type MockSMSProvider struct{}

func NewMockSMSProvider() *MockSMSProvider {
	return &MockSMSProvider{}
}

func (m *MockSMSProvider) SendOTP(ctx context.Context, mobile, otp string) error {
	log.Info().
		Str("mobile", mobile).
		Str("otp", otp).
		Msg("MOCK SMS: OTP sent (development mode)")
	return nil
}
