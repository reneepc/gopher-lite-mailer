package mailer

import (
	"net/smtp"
	"testing"
)

func TestMailerBuilder(t *testing.T) {
	tests := map[string]struct {
		builder        MailerBuilder
		expectedMailer Mailer
	}{
		"Gmail Builder": {
			builder: NewGMailMailerBuilder("user@gmail.com", "password"),
			expectedMailer: Mailer{
				server: "smtp.gmail.com:587",
				from:   "user@gmail.com",
				auth:   smtp.PlainAuth("", "user@gmail.com", "password", "smtp.gmail.com"),
			},
		},
		"Outlook Builder": {
			builder: NewOutlookMailerBuilder("user@outlook.com", "password"),
			expectedMailer: Mailer{
				server: "smtp-mail.outlook.com:587",
				from:   "user@outlook.com",
				auth:   smtp.PlainAuth("", "user@outlook.com", "password", "smtp-mail.outlook.com"),
			},
		},
		"Custom Headers": {
			builder: NewGMailMailerBuilder("user@gmail.com", "password").
				WithHeader("X-Custom-Header", "CustomValue"),
			expectedMailer: Mailer{
				server: "smtp.gmail.com:587",
				from:   "user@gmail.com",
				auth:   smtp.PlainAuth("", "user@gmail.com", "password", "smtp.gmail.com"),
				customHeaders: map[string]string{
					"X-Custom-Header": "CustomValue",
				},
			},
		},
		"Attachments": {
			builder: NewGMailMailerBuilder("user@gmail.com", "password").
				WithAttachment("file.txt", "text/plain", "file123", true),
			expectedMailer: Mailer{
				server: "smtp.gmail.com:587",
				from:   "user@gmail.com",
				auth:   smtp.PlainAuth("", "user@gmail.com", "password", "smtp.gmail.com"),
				attachments: []Attachment{{
					FileName:     "file.txt",
					ContentType:  "text/plain",
					Base64Encode: true,
					ContentID:    "file123",
				}},
			},
		},
		"Custom Host and Port": {
			builder: NewGMailMailerBuilder("user@gmail.com", "password").
				WithHost("smtp.custom.com").
				WithPort(2525),
			expectedMailer: Mailer{
				server: "smtp.custom.com:2525",
				from:   "user@gmail.com",
				auth:   smtp.PlainAuth("", "user@gmail.com", "password", "smtp.gmail.com"),
			},
		},
		"Custom Builder": {
			builder: NewMailerBuilder("smtp.custom.com", 2525, "user@custom.com", "password"),
			expectedMailer: Mailer{
				server: "smtp.custom.com:2525",
				from:   "user@custom.com",
				auth:   smtp.PlainAuth("", "user@custom.com", "password", "smtp.custom.com"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mailer := tt.builder.Build()

			if !compareMailers(tt.expectedMailer, mailer) {
				t.Errorf("expected mailer %+v, got %+v", tt.expectedMailer, mailer)
			}
		})
	}
}

func compareMailers(a, b Mailer) bool {
	if a.server != b.server {
		return false
	}
	if a.from != b.from {
		return false
	}
	if !compareSMTPAuth(a.auth, b.auth) {
		return false
	}
	if !compareHeaders(a.customHeaders, b.customHeaders) {
		return false
	}
	if !compareAttachments(a.attachments, b.attachments) {
		return false
	}
	return true
}

func compareSMTPAuth(a, b smtp.Auth) bool {
	if a == nil && b == nil {
		return true
	}

	if a != nil && b != nil {
		return true
	}

	return false
}

func compareHeaders(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func compareAttachments(a, b []Attachment) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}
