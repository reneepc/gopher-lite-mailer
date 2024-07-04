package mailer

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

type Mailer struct {
	server        string
	from          string
	auth          smtp.Auth
	customHeaders map[string]string
	attachments   []Attachment
}

type Attachment struct {
	FileName     string
	ContentType  string
	Base64Encode bool
	ContentID    string
}

func (m Mailer) SendMail(to, subject string, data string) error {
	recipient, err := mail.ParseAddress(to)
	if err != nil {
		return fmt.Errorf("invalid recipient: %v", err)
	}

	headers := map[string]string{
		"From":         m.from,
		"To":           recipient.Address,
		"Subject":      subject,
		"MIME-Version": "1.0",
	}

	if len(m.attachments) > 0 {
		headers["Content-Type"] = "multipart/related; boundary=boundary"
	} else {
		headers["Content-Type"] = "text/html; charset=\"UTF-8\""
	}

	for k, v := range m.customHeaders {
		headers[k] = v
	}

	var msg string
	if len(m.attachments) > 0 {
		msg, err = m.buildMultipartEmail(data, headers)
		if err != nil {
			return fmt.Errorf("error building multipart email: %v", err)
		}
	} else {
		msg = m.buildSimpleEmail(data, headers)
	}

	err = smtp.SendMail(
		m.server,
		m.auth,
		m.from,
		[]string{recipient.Address},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Errorf("error sending mail: %v", err)
	}

	return nil
}

func (m Mailer) buildMultipartEmail(data string, headers map[string]string) (string, error) {
	var msg strings.Builder
	msg.WriteString(buildHeaders(headers))

	msg.WriteString("\r\n--boundary\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString(data)

	for _, attachment := range m.attachments {
		attachmentPart, err := attachment.buildPart()
		if err != nil {
			return "", fmt.Errorf("could not build attachment part: %v", err)
		}

		msg.WriteString(attachmentPart)
	}
	msg.WriteString("\r\n--boundary--")

	return msg.String(), nil
}

func (m Mailer) buildSimpleEmail(data string, headers map[string]string) string {
	var msg strings.Builder
	msg.WriteString(buildHeaders(headers))

	msg.WriteString(data)
	return msg.String()
}

func (a Attachment) buildPart() (string, error) {
	var part strings.Builder
	content, err := os.ReadFile(a.FileName)
	if err != nil {
		return "", fmt.Errorf("could not read attachment file: %v", err)
	}

	part.WriteString("\r\n--boundary\r\n")
	part.WriteString(fmt.Sprintf("Content-Type: %s\r\n", a.ContentType))
	part.WriteString(fmt.Sprintf("Content-ID: <%s>\r\n\r\n", a.ContentID))

	if a.Base64Encode {
		part.WriteString("Content-Transfer-Encoding: base64\r\n")
		part.WriteString(base64.StdEncoding.EncodeToString(content))
	}

	return part.String(), nil
}

func buildHeaders(headers map[string]string) string {
	var headerBuilder strings.Builder
	for k, v := range headers {
		headerBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	headerBuilder.WriteString("\r\n")
	return headerBuilder.String()
}
