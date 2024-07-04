package mailer

import (
	"fmt"
	"net/smtp"
)

type MailerBuilder struct {
	smtpHost      string
	smtpPort      int
	from          string
	password      string
	auth          smtp.Auth
	customHeaders map[string]string
	attachments   []Attachment
}

func NewMailerBuilder(smtpHost string, SMTPPort int, from, password string) MailerBuilder {
	auth := smtp.PlainAuth("", from, password, smtpHost)
	return MailerBuilder{
		smtpHost:      smtpHost,
		smtpPort:      SMTPPort,
		from:          from,
		password:      password,
		auth:          auth,
		customHeaders: make(map[string]string),
		attachments:   make([]Attachment, 0),
	}
}

func NewGMailMailerBuilder(from, password string) MailerBuilder {
	return NewMailerBuilder("smtp.gmail.com", 587, from, password)
}

func NewOutlookMailerBuilder(from, password string) MailerBuilder {
	return NewMailerBuilder("smtp-mail.outlook.com", 587, from, password)
}

func (b MailerBuilder) WithHeader(key, value string) MailerBuilder {
	b.customHeaders[key] = value
	return b
}

func (b MailerBuilder) WithAttachment(fileName, contentType, contentID string, base64Encode bool) MailerBuilder {
	b.attachments = append(b.attachments, Attachment{
		FileName:     fileName,
		ContentType:  contentType,
		Base64Encode: base64Encode,
		ContentID:    contentID,
	})

	return b
}

func (b MailerBuilder) WithHost(host string) MailerBuilder {
	b.smtpHost = host
	return b
}

func (b MailerBuilder) WithPort(port int) MailerBuilder {
	b.smtpPort = port
	return b
}

func (b MailerBuilder) Build() Mailer {
	return Mailer{
		server:        fmt.Sprintf("%s:%d", b.smtpHost, b.smtpPort),
		from:          b.from,
		auth:          b.auth,
		customHeaders: b.customHeaders,
		attachments:   b.attachments,
	}
}
