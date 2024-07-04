package mailer

import (
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path"
	"strings"
)

type EmailTemplate struct {
	Tmpl          *template.Template
	css           string
	signatureLink string
}

type TemplateData struct {
	CSS       template.CSS
	Signature template.URL
	Data      map[string]string
}

func NewEmailTemplate(templateDir, templateFile, cssFile, signatureLink string) (EmailTemplate, error) {
	templateContent, err := template.ParseFiles(path.Join(templateDir, templateFile))
	if err != nil {
		slog.Error("could not parse template file: %v", slog.Any("error", err))
		return EmailTemplate{}, err
	}

	css, err := os.ReadFile(cssFile)
	if err != nil {
		slog.Warn("could not read CSS file: %v", slog.Any("error", err))
	}

	return EmailTemplate{
		Tmpl:          templateContent,
		css:           string(css),
		signatureLink: signatureLink,
	}, nil
}

func (t *EmailTemplate) Execute(data map[string]string) (string, error) {
	templateData := TemplateData{
		CSS:       template.CSS(t.css),
		Signature: template.URL(t.signatureLink),
		Data:      data,
	}

	var body strings.Builder
	err := t.Tmpl.Execute(&body, templateData)
	if err != nil {
		return "", fmt.Errorf("could not execute template: %v", err)
	}

	return body.String(), nil
}
