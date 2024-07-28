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
	TmplHeader    *template.Template
	TmplFooter    *template.Template
	TmplBody      *template.Template
	css           string
	signatureLink string
}

type TemplateData struct {
	CSS       template.CSS
	Signature template.URL
	Data      map[string]string
}

func NewEmailTemplate(templateDir, bodyFile, signatureLink string) (EmailTemplate, error) {
	headerFile := path.Join(templateDir, "header.html")
	footerFile := path.Join(templateDir, "footer.html")
	bodyFilePath := path.Join(templateDir, "bodies", bodyFile)
	cssFilePath := path.Join(templateDir, "styles.css")

	tmplHeader, err := template.ParseFiles(headerFile)
	if err != nil {
		slog.Error("could not parse header template file: %v", slog.Any("error", err))
		return EmailTemplate{}, err
	}

	templateContent, err := template.ParseFiles(bodyFilePath)
	if err != nil {
		slog.Error("could not parse template file: %v", slog.Any("error", err))
		return EmailTemplate{}, err
	}

	tmplFooter, err := template.ParseFiles(footerFile)
	if err != nil {
		slog.Error("could not parse footer template file: %v", slog.Any("error", err))
		return EmailTemplate{}, err
	}

	css, err := os.ReadFile(cssFilePath)
	if err != nil {
		slog.Warn("could not read CSS file: %v", slog.Any("error", err))
	}

	return EmailTemplate{
		TmplHeader:    tmplHeader,
		TmplFooter:    tmplFooter,
		TmplBody:      templateContent,
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

	err := t.TmplHeader.Execute(&body, templateData)
	if err != nil {
		return "", fmt.Errorf("could not execute header template: %v", err)
	}

	err = t.TmplBody.Execute(&body, templateData)
	if err != nil {
		return "", fmt.Errorf("could not execute template: %v", err)
	}

	err = t.TmplFooter.Execute(&body, templateData)
	if err != nil {
		return "", fmt.Errorf("could not execute footer template: %v", err)
	}

	return body.String(), nil
}
