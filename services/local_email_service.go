package services

import (
	"context"
	"fmt"
	"os"
	servcfg "p2p-management-service/config"
	"path/filepath"
	"time"
)

type LocalEmailService struct {
	SavePath string
}

func NewLocalEmailService() *LocalEmailService {
	path := servcfg.AppConfig.EMAIL_SAVE_PATH
	if path == "" {
		path = "./logs/emails"
	}
	// ensure directory exists
	_ = os.MkdirAll(path, 0o755)

	return &LocalEmailService{SavePath: path}
}

// SendEmailAsync implements EmailService and writes the email to a file for local testing.
func (s *LocalEmailService) SendEmailAsync(ctx context.Context, req EmailRequest) {
	go func() {
		timestamp := time.Now().UTC().Format("20060102_150405")
		filename := fmt.Sprintf("email_%s_%s.txt", timestamp, req.To)
		fullpath := filepath.Join(s.SavePath, filename)

		content := fmt.Sprintf("To: %s\nSubject: %s\n\n%s\n", req.To, req.Subject, req.Body)
		if err := os.WriteFile(fullpath, []byte(content), 0o644); err != nil {
			fmt.Printf("❌ Failed to save email to %s: %v\n", fullpath, err)
			return
		}
		fmt.Printf("✅ Email saved locally to %s\n", fullpath)
	}()
}
