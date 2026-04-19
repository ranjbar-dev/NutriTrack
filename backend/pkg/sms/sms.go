package sms

// Sender defines the interface for sending OTP messages.
// Implementations include MockSender (development) and KavenegarSender (production).
type Sender interface {
	SendOTP(phone string, code string) error
}
