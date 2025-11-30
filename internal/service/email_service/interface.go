package emailservice

import "context"

type EmailManager interface {
	SendEmail(ctx context.Context, to, head, body string) error
}
