// pkg/email/mock.go
package email

import (
	"log"
)

type MockEmailService struct{}

func (m *MockEmailService) SendVerificationCode(to string, code string) error {
	log.Printf("[MOCK EMAIL] Sending verification code '%s' to: %s\n", code, to)
	return nil
}
