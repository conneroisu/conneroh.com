package email

import (
	"bytes"
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

// Sender defines the interface for sending emails
type Sender interface {
	SendContactEmail(msg ContactMessage) (string, error)
}

// Config holds configuration for the email service
type Config struct {
	FromEmail string
	ToEmail   string
}

// Service implements the Sender interface using Resend
type Service struct {
	client *resend.Client
	config Config
}

type ContactMessage struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func NewService(apiKey string) *Service {
	// Use default config for backward compatibility
	defaultConfig := Config{
		FromEmail: "noreply@conneroh.com",
		ToEmail:   "conneroisu@outlook.com",
	}
	return NewServiceWithConfig(apiKey, defaultConfig)
}

func NewServiceWithConfig(apiKey string, config Config) *Service {
	return &Service{
		client: resend.NewClient(apiKey),
		config: config,
	}
}

// MockEmailSender is a mock implementation of the Sender interface for testing
type MockEmailSender struct {
	SentEmails []ContactMessage
	ShouldFail bool
	EmailID    string
}

func (m *MockEmailSender) SendContactEmail(msg ContactMessage) (string, error) {
	if m.ShouldFail {
		return "", fmt.Errorf("mock email send failed")
	}
	
	m.SentEmails = append(m.SentEmails, msg)
	
	if m.EmailID != "" {
		return m.EmailID, nil
	}
	
	return "mock-email-id-123", nil
}

func (s *Service) SendContactEmail(msg ContactMessage) (string, error) {
	// Render the email template using templ
	buf := &bytes.Buffer{}
	component := ContactEmailTemplate(msg)
	if err := component.Render(context.Background(), buf); err != nil {
		return "", fmt.Errorf("failed to render email template: %w", err)
	}
	htmlContent := buf.String()

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("Contact Form <%s>", s.config.FromEmail),
		To:      []string{s.config.ToEmail},
		Html:    htmlContent,
		Subject: fmt.Sprintf("Contact Form: %s", msg.Subject),
		ReplyTo: msg.Email,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return sent.Id, nil
}