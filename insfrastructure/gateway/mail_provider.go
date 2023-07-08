package gateway

import (
	"fmt"

	"github.com/mendezdev/rate-limit-example/core/ports"
)

type mailSenderGateway struct{}

func NewMailSender() ports.MailSenderGateway {
	return &mailSenderGateway{}
}

// Send will send email for the given user_id with the given message
func (ms *mailSenderGateway) Send(userID string, message string) error {
	fmt.Printf("sending message to user %s\n", userID)
	return nil
}
