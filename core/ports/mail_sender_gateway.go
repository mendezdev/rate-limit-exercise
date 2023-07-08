package ports

//go:generate mockgen -destination=mocks/mock_mail_sender_gateway.go -package=ports -source=mail_sender_gateway.go MailSenderGateway

type MailSenderGateway interface {
	Send(userID string, message string) error
}
