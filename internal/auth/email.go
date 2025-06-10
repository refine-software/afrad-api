package auth

import (
	"github.com/refine-software/afrad-api/config"
	"gopkg.in/gomail.v2"
)

func SendOtpEmail(email, OTP string, env *config.Env) error {
	htmlBody := generateTemplate(OTP)

	m := gomail.NewMessage()
	m.SetHeader("From", env.Email)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Afrad OTP Email")
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(
		"smtp.hostinger.com",
		465,
		env.Email,
		env.Password,
	)
	d.SSL = true
	return d.DialAndSend(m)
}
