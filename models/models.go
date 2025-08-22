package models

// All request/response structs used across packages

type EmailAttachment struct {
	Filename string `json:"filename" validate:"required"`
	Content  string `json:"content"  validate:"required"` // base64 encoded
}

type EmailRequest struct {
	Sender      string            `json:"sender"      validate:"required,email"`
	Password    string            `json:"password"    validate:"required"`
	SMTPServer  string            `json:"smtp_server" validate:"required"`
	SMTPPort    int               `json:"smtp_port"   validate:"required,min=1"`
	UseTLS      bool              `json:"use_tls"`
	Recipient   string            `json:"recipient"   validate:"required,email"`
	Subject     string            `json:"subject"     validate:"required"`
	Body        string            `json:"body"        validate:"required"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
}

type APIResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
