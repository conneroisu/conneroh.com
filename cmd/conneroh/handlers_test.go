package conneroh

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/conneroisu/conneroh.com/internal/email"
)

func TestHandleContactForm_Success(t *testing.T) {
	// Create a mock email sender
	mockSender := &email.MockEmailSender{
		EmailID: "test-email-123",
	}

	// Create the handler
	handler := handleContactForm(mockSender)

	// Create form data
	form := url.Values{}
	form.Add("name", "John Doe")
	form.Add("email", "john@example.com")
	form.Add("subject", "Test Subject")
	form.Add("message", "This is a test message")

	// Create request
	req := httptest.NewRequest("POST", "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call the handler
	err := handler(w, req)
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	// Check that response is successful
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check that email was sent
	if len(mockSender.SentEmails) != 1 {
		t.Fatalf("expected 1 email to be sent, got %d", len(mockSender.SentEmails))
	}

	sentEmail := mockSender.SentEmails[0]
	if sentEmail.Name != "John Doe" {
		t.Errorf("expected name 'John Doe', got '%s'", sentEmail.Name)
	}
	if sentEmail.Email != "john@example.com" {
		t.Errorf("expected email 'john@example.com', got '%s'", sentEmail.Email)
	}
	if sentEmail.Subject != "Test Subject" {
		t.Errorf("expected subject 'Test Subject', got '%s'", sentEmail.Subject)
	}
	if sentEmail.Message != "This is a test message" {
		t.Errorf("expected message 'This is a test message', got '%s'", sentEmail.Message)
	}
}

func TestHandleContactForm_EmailError(t *testing.T) {
	// Create a mock email sender that fails
	mockSender := &email.MockEmailSender{
		ShouldFail: true,
	}

	// Create the handler
	handler := handleContactForm(mockSender)

	// Create form data
	form := url.Values{}
	form.Add("name", "John Doe")
	form.Add("email", "john@example.com")
	form.Add("subject", "Test Subject")
	form.Add("message", "This is a test message")

	// Create request
	req := httptest.NewRequest("POST", "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call the handler - should not return error even if email fails
	err := handler(w, req)
	if err != nil {
		t.Fatalf("handler should not return error even if email fails: %v", err)
	}

	// Check that response is still successful (shows thank you page)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 even with email error, got %d", w.Code)
	}

	// Check that no email was sent
	if len(mockSender.SentEmails) != 0 {
		t.Errorf("expected 0 emails to be sent, got %d", len(mockSender.SentEmails))
	}
}

func TestHandleContactForm_NoEmailService(t *testing.T) {
	// Create the handler with no email service
	handler := handleContactForm(nil)

	// Create form data
	form := url.Values{}
	form.Add("name", "John Doe")
	form.Add("email", "john@example.com")
	form.Add("subject", "Test Subject")
	form.Add("message", "This is a test message")

	// Create request
	req := httptest.NewRequest("POST", "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call the handler
	err := handler(w, req)
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	// Check that response is successful
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHandleContactForm_InvalidForm(t *testing.T) {
	mockSender := &email.MockEmailSender{}
	handler := handleContactForm(mockSender)

	// Create request with invalid form data (missing required fields)
	form := url.Values{}
	form.Add("name", "") // Empty name should cause validation error

	req := httptest.NewRequest("POST", "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	// This should return an error because schema validation requires all fields
	err := handler(w, req)
	if err == nil {
		t.Fatalf("handler should return error for missing required fields")
	}

	// Should contain validation error message
	if !strings.Contains(err.Error(), "failed to decode contact form") {
		t.Errorf("expected validation error, got: %v", err)
	}
}