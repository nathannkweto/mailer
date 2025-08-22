package email

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	gomail "gopkg.in/gomail.v2"

	"github.com/nathannkweto/mailer/applog"
	"github.com/nathannkweto/mailer/models"
)

// SendEmail sends the message synchronously; it is the caller's responsibility to call with timeout/cancellation.
func SendEmail(ctx context.Context, req models.EmailRequest, attachmentPaths []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", req.Sender)
	m.SetHeader("To", req.Recipient)
	m.SetHeader("Subject", req.Subject)

	if looksLikeHTML(req.Body) {
		m.SetBody("text/html", req.Body)
	} else {
		m.SetBody("text/plain", req.Body)
	}

	for _, p := range attachmentPaths {
		// attach and continue even if missing
		if _, err := os.Stat(p); err == nil {
			m.Attach(p)
		} else {
			applog.Log.WithField("path", p).Warn("attachment missing when attaching")
		}
	}

	dialer := gomail.NewDialer(req.SMTPServer, req.SMTPPort, req.Sender, req.Password)

	// quick TLS control: if UseTLS false, set TLSConfig to nil
	if !req.UseTLS {
		dialer.TLSConfig = nil
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- dialer.DialAndSend(m)
	}()

	select {
	case <-ctx.Done():
		// context cancelled or timed out
		// Note: gomail's DialAndSend can't be cancelled directly. We return a ctx error here.
		return ctx.Err()
	case err := <-errCh:
		if err != nil {
			// wrap the error with safe message
			return fmt.Errorf("smtp send error: %w", err)
		}
		return nil
	case <-time.After(5 * time.Minute):
		// safety fallback
		return fmt.Errorf("send timed out after fallback period")
	}
}

func looksLikeHTML(s string) bool {
	if s == "" {
		return false
	}
	l := strings.ToLower(s)
	return strings.Contains(l, "<html") || strings.Contains(l, "<body") || strings.Contains(l, "<div")
}
