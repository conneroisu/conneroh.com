package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/resend/resend-go/v2"
)

// Sender defines the interface for sending emails
type Sender interface {
	SendContactEmail(msg ContactMessage) (string, error)
}

// Service implements the Sender interface using Resend
type Service struct {
	client *resend.Client
}

type ContactMessage struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func NewService(apiKey string) *Service {
	return &Service{
		client: resend.NewClient(apiKey),
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
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Contact Form Submission</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px;">
            New Contact Form Submission
        </h2>
        
        <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <p><strong>From:</strong> {{.Name}} ({{.Email}})</p>
            <p><strong>Subject:</strong> {{.Subject}}</p>
        </div>
        
        <div style="background-color: #fff; padding: 20px; border-left: 4px solid #3498db; margin: 20px 0;">
            <h3 style="margin-top: 0; color: #2c3e50;">Message:</h3>
            <p style="white-space: pre-wrap;">{{.Message}}</p>
        </div>
        
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #666; text-align: center;">
            This email was sent from the contact form on conneroh.com
        </p>
    </div>
</body>
</html>`

	tmpl, err := template.New("contact").Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	var htmlContent string
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, msg); err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}
	htmlContent = buf.String()

	params := &resend.SendEmailRequest{
		From:    "Contact Form <noreply@conneroh.com>",
		To:      []string{"conneroisu@outlook.com"},
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