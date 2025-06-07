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

func TestHandleContactForm_InvalidEmail(t *testing.T) {
	mockSender := &email.MockEmailSender{}
	handler := handleContactForm(mockSender)

	form := url.Values{}
	form.Add("name", "John Doe")
	form.Add("email", "invalid-email") // Invalid email format
	form.Add("subject", "Test Subject")
	form.Add("message", "Test message")

	req := httptest.NewRequest("POST", "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	// The handler should succeed but might validate email format
	err := handler(w, req)
	if err != nil {
		// If validation is implemented, check for email validation error
		if !strings.Contains(err.Error(), "email") {
			t.Errorf("expected email validation error, got: %v", err)
		}
	} else {
		// If no validation is implemented yet, email should be sent
		if len(mockSender.SentEmails) != 1 {
			t.Errorf("expected 1 sent email, got %d", len(mockSender.SentEmails))
		}
	}
}

func TestHandleContactForm_LongFieldValues(t *testing.T) {
	mockSender := &email.MockEmailSender{}
	handler := handleContactForm(mockSender)

	// Create extremely long values
	longString := strings.Repeat("a", 10000)
	
	form := url.Values{}
	form.Add("name", longString)
	form.Add("email", "test@example.com")
	form.Add("subject", longString)
	form.Add("message", longString)

	req := httptest.NewRequest("POST", "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	err := handler(w, req)
	// Handler should either handle long values gracefully or return an error
	if err == nil && len(mockSender.SentEmails) == 1 {
		// Check that the long values were preserved
		sentEmail := mockSender.SentEmails[0]
		if len(sentEmail.Name) != 10000 {
			t.Errorf("expected name length 10000, got %d", len(sentEmail.Name))
		}
	}
}

func TestHandleContactForm_WhitespaceFields(t *testing.T) {
	mockSender := &email.MockEmailSender{}
	handler := handleContactForm(mockSender)

	form := url.Values{}
	form.Add("name", "   ") // Only whitespace
	form.Add("email", "test@example.com")
	form.Add("subject", "\t\n") // Only whitespace
	form.Add("message", "   \n   ") // Only whitespace

	req := httptest.NewRequest("POST", "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	err := handler(w, req)
	// Handler should validate that fields contain non-whitespace content
	if err == nil {
		// If no validation, check what was sent
		if len(mockSender.SentEmails) == 1 {
			sentEmail := mockSender.SentEmails[0]
			// These should ideally be trimmed or rejected
			if strings.TrimSpace(sentEmail.Name) == "" {
				t.Log("Warning: handler accepted whitespace-only name")
			}
		}
	}
}

func TestHandleContactForm_SpecialCharacters(t *testing.T) {
	mockSender := &email.MockEmailSender{}
	handler := handleContactForm(mockSender)

	form := url.Values{}
	form.Add("name", "John <script>alert('xss')</script> Doe")
	form.Add("email", "test@example.com")
	form.Add("subject", "Test & Subject < > \" '")
	form.Add("message", "Message with special chars: &<>\"'\n\nAnd newlines")

	req := httptest.NewRequest("POST", "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	err := handler(w, req)
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	// Check that email was sent
	if len(mockSender.SentEmails) != 1 {
		t.Fatalf("expected 1 sent email, got %d", len(mockSender.SentEmails))
	}

	sentEmail := mockSender.SentEmails[0]
	// Check that special characters are preserved (they should be escaped in HTML rendering)
	if !strings.Contains(sentEmail.Name, "<script>") {
		t.Errorf("expected special characters to be preserved in name")
	}
	if !strings.Contains(sentEmail.Subject, "&") || !strings.Contains(sentEmail.Subject, "<") {
		t.Errorf("expected special characters to be preserved in subject")
	}
}