package services

// PostfixEmailService implements the EmailService interface using the local sendmail command.

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	servcfg "p2p-management-service/config"
)

type PostfixEmailService struct {
	sendmailPath string
	sender       string
	maxRetries   int
}

func NewPostfixEmailService() *PostfixEmailService {
	// Default to standard sendmail path if not configured
	path := servcfg.AppConfig.SENDMAIL_PATH
	if path == "" {
		path = "/usr/sbin/sendmail"
	}

	return &PostfixEmailService{
		sendmailPath: path,
		sender:       servcfg.AppConfig.SES_SENDER, // Using SES_SENDER as the FROM address
		maxRetries:   3,
	}
}

func (s *PostfixEmailService) SendEmailAsync(ctx context.Context, req EmailRequest) {
	go func() {
		var lastErr error

		// Construct crude email message with headers
		// Note: A more robust implementation might use net/mail or similar to format headers
		// But for simple HTML emails, this suffices.
		headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
			s.sender, req.To, req.Subject)
		body := headers + req.Body

		for attempt := 1; attempt <= s.maxRetries; attempt++ {
			// -t: Read message for recipients. -i: Ignore dots.
			cmd := exec.CommandContext(ctx, s.sendmailPath, "-t", "-i")
			cmd.Stdin = bytes.NewReader([]byte(body))

			// Capture output for debugging
			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			fmt.Printf("Sending email via %s to %s (attempt %d)\n", s.sendmailPath, req.To, attempt)

			err := cmd.Run()
			if err == nil {
				fmt.Printf("✅ Email handed off to sendmail for %s\n", req.To)
				return
			}

			lastErr = fmt.Errorf("%w (stderr: %s)", err, stderr.String())

			// Exponential backoff
			wait := time.Duration(attempt*2) * time.Second
			// Add jitter
			wait += time.Duration(rand.Intn(1000)) * time.Millisecond

			log.Printf("⚠️ Attempt %d failed to invoke sendmail for %s: %v. Retrying in %v...", attempt, req.To, err, wait)
			time.Sleep(wait)
		}

		// All retries failed
		log.Printf("❌ Failed to send email to %s after %d attempts: %v", req.To, s.maxRetries, lastErr)
	}()
}
