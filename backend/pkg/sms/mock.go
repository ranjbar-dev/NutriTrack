package sms

import "github.com/rs/zerolog"

// MockSender logs OTP messages to stdout instead of sending real SMS.
// Used in development mode per D-04.
type MockSender struct {
	Logger zerolog.Logger
}

// NewMockSender creates a new MockSender with the given logger.
func NewMockSender(logger zerolog.Logger) *MockSender {
	return &MockSender{Logger: logger}
}

// SendOTP logs the OTP code and phone number instead of sending an SMS.
func (m *MockSender) SendOTP(phone, code string) error {
	m.Logger.Info().
		Str("phone", phone).
		Str("code", code).
		Msg("📱 OTP sent (mock mode)")
	return nil
}
