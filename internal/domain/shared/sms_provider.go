package shared

import "context"

// SMSProvider is the domain port for sending SMS messages.
// The Kavenegar adapter implements this interface in infrastructure/sms/.
type SMSProvider interface {
	SendOTP(ctx context.Context, mobile, otp string) error
}
