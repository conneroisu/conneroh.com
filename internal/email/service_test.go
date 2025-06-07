package email

import (
	"testing"
)


func TestMockEmailSender_SendContactEmail_Success(t *testing.T) {
	mock := &MockEmailSender{
		EmailID: "test-email-id",
	}
	
	msg := ContactMessage{
		Name:    "John Doe",
		Email:   "john@example.com",
		Subject: "Test Subject",
		Message: "Test message content",
	}
	
	emailID, err := mock.SendContactEmail(msg)
	
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	
	if emailID != "test-email-id" {
		t.Fatalf("expected email ID 'test-email-id', got: %s", emailID)
	}
	
	if len(mock.SentEmails) != 1 {
		t.Fatalf("expected 1 sent email, got: %d", len(mock.SentEmails))
	}
	
	sentEmail := mock.SentEmails[0]
	if sentEmail.Name != msg.Name {
		t.Errorf("expected name %s, got %s", msg.Name, sentEmail.Name)
	}
	if sentEmail.Email != msg.Email {
		t.Errorf("expected email %s, got %s", msg.Email, sentEmail.Email)
	}
	if sentEmail.Subject != msg.Subject {
		t.Errorf("expected subject %s, got %s", msg.Subject, sentEmail.Subject)
	}
	if sentEmail.Message != msg.Message {
		t.Errorf("expected message %s, got %s", msg.Message, sentEmail.Message)
	}
}

func TestMockEmailSender_SendContactEmail_Failure(t *testing.T) {
	mock := &MockEmailSender{
		ShouldFail: true,
	}
	
	msg := ContactMessage{
		Name:    "John Doe",
		Email:   "john@example.com",
		Subject: "Test Subject",
		Message: "Test message content",
	}
	
	emailID, err := mock.SendContactEmail(msg)
	
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	
	if emailID != "" {
		t.Fatalf("expected empty email ID, got: %s", emailID)
	}
	
	if len(mock.SentEmails) != 0 {
		t.Fatalf("expected 0 sent emails, got: %d", len(mock.SentEmails))
	}
}

func TestContactMessage_Fields(t *testing.T) {
	msg := ContactMessage{
		Name:    "Jane Smith",
		Email:   "jane@example.com", 
		Subject: "Hello World",
		Message: "This is a test message with multiple lines\nand special characters: &<>\"'",
	}
	
	// Test that all fields are properly set
	if msg.Name == "" {
		t.Error("Name should not be empty")
	}
	if msg.Email == "" {
		t.Error("Email should not be empty")
	}
	if msg.Subject == "" {
		t.Error("Subject should not be empty")
	}
	if msg.Message == "" {
		t.Error("Message should not be empty")
	}
}

func TestNewServiceWithConfig(t *testing.T) {
	config := Config{
		FromEmail: "test@example.com",
		ToEmail:   "recipient@example.com",
	}
	
	service := NewServiceWithConfig("test-api-key", config)
	
	if service == nil {
		t.Fatal("expected service to be created")
	}
	
	if service.config.FromEmail != config.FromEmail {
		t.Errorf("expected FromEmail %s, got %s", config.FromEmail, service.config.FromEmail)
	}
	
	if service.config.ToEmail != config.ToEmail {
		t.Errorf("expected ToEmail %s, got %s", config.ToEmail, service.config.ToEmail)
	}
	
	if service.tmpl == nil {
		t.Error("expected template to be parsed")
	}
	
	if service.client == nil {
		t.Error("expected client to be initialized")
	}
}

func TestNewService_DefaultConfig(t *testing.T) {
	service := NewService("test-api-key")
	
	if service == nil {
		t.Fatal("expected service to be created")
	}
	
	if service.config.FromEmail != "noreply@conneroh.com" {
		t.Errorf("expected default FromEmail 'noreply@conneroh.com', got %s", service.config.FromEmail)
	}
	
	if service.config.ToEmail != "conneroisu@outlook.com" {
		t.Errorf("expected default ToEmail 'conneroisu@outlook.com', got %s", service.config.ToEmail)
	}
}