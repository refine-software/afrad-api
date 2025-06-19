package auth

import (
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendOtpEmail(userEmail, OTP string) error
}

type EmailService struct {
	Email    string
	Password string
}

func NewEmailService(serviceEmail, servicePassword string) *EmailService {
	return &EmailService{
		Email:    serviceEmail,
		Password: servicePassword,
	}
}

func (e *EmailService) SendOtpEmail(userEmail, OTP string) error {
	htmlBody := generateTemplate(OTP)

	m := gomail.NewMessage()
	m.SetHeader("From", e.Email)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Afrad OTP Email")
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(
		"smtp.hostinger.com",
		465,
		e.Email,
		e.Password,
	)
	d.SSL = true
	return d.DialAndSend(m)
}
