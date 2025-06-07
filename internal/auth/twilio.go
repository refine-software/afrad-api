package auth

import (
	"fmt"

	"github.com/refine-software/afrad-api/config"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Twilio struct {
	client *twilio.RestClient
	form   string
}

func NewTwilio(env *config.Env) Twilio {
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioAccountSID,
		Password: env.TwilioAuthToken,
	})

	return Twilio{
		client: twilioClient,
		form:   env.TwilioSandboxFrom,
	}
}

func (t *Twilio) SendWhatsAppOTP(phoneNumber string, otp string) error {
	if t.client == nil {
		return fmt.Errorf("Twilio client is not initialized")
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo("whatsapp:" + phoneNumber)
	params.SetFrom(t.form)
	params.SetBody(fmt.Sprintf("Your OTP is %s\n\nPlease don't share this code to anyone.", otp))

	_, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %v", err)
	}

	return nil
}
