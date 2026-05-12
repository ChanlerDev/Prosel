package subscribe

import "context"

type MailMessage struct {
	To      string
	Subject string
	Body    string
}

type Mailer interface {
	Send(ctx context.Context, message MailMessage) error
}
