package services

import "context"

type EmailRequest struct {
	To      string
	Subject string
	Body    string
}

type EmailService interface {
	SendEmailAsync(ctx context.Context, req EmailRequest)
}
