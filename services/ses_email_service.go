package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	servcfg "p2p-management-service/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESEmailService struct {
	client     *ses.Client
	sender     string
	maxRetries int
}

func NewSESEmailService() *SESEmailService {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(servcfg.AppConfig.AWS_REGION),
	)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	return &SESEmailService{
		client:     ses.NewFromConfig(cfg),
		sender:     servcfg.AppConfig.SES_SENDER,
		maxRetries: 3, // configurable retry count
	}
}

func (s *SESEmailService) SendEmailAsync(ctx context.Context, req EmailRequest) {
	go func() {
		var lastErr error

		for attempt := 1; attempt <= s.maxRetries; attempt++ {
			input := &ses.SendEmailInput{
				Destination: &types.Destination{
					ToAddresses: []string{req.To},
				},
				Message: &types.Message{
					Body: &types.Body{
						Html: &types.Content{
							Charset: aws.String("UTF-8"),
							Data:    aws.String(req.Body),
						},
					},
					Subject: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(req.Subject),
					},
				},
				Source: aws.String(s.sender),
			}

			fmt.Println("Sending email", input)

			_, err := s.client.SendEmail(ctx, input)
			if err == nil {
				fmt.Printf("✅ Email sent successfully to %s\n", req.To)
				return
			}

			lastErr = err
			wait := time.Duration(attempt*2) * time.Second
			// add jitter to avoid thundering herd in case of multiple retries
			wait += time.Duration(rand.Intn(1000)) * time.Millisecond
			log.Printf("⚠️ Attempt %d failed to send email to %s: %v. Retrying in %v...", attempt, req.To, err, wait)
			time.Sleep(wait)
		}

		// All retries failed
		log.Printf("❌ Failed to send email to %s after %d attempts: %v", req.To, s.maxRetries, lastErr)
	}()
}
