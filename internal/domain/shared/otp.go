package shared

import (
	"crypto/rand"
	"math/big"
)

// GenerateOTP generates a cryptographically secure numeric OTP of the given length.
func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[n.Int64()]
	}
	return string(otp), nil
}
