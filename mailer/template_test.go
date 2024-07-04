package mailer

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempFile(t *testing.T, dir, pattern, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		t.Fatalf("Failed to write content to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	return tmpfile.Name()
}

func TestNewEmailTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	validTemplateContent := "Hello, {{.Data.Name}}!"
	validCSSContent := "body { color: red; }"

	validTemplateFile := createTempFile(t, tmpDir, "template-*.html", validTemplateContent)
	validCSSFile := createTempFile(t, tmpDir, "styles-*.css", validCSSContent)
	invalidTemplateFile := createTempFile(t, tmpDir, "template-*.html", "{{.Invalid")

	tests := map[string]struct {
		templateFile  string
		cssFile       string
		signatureLink string
		expectError   bool
	}{
		"Valid Files": {
			validTemplateFile,
			validCSSFile,
			"http://golang.samba.br",
			false,
		},
		"Invalid Template File": {
			invalidTemplateFile,
			validCSSFile,
			"http://golang.samba.br",
			true,
		},
		"Invalid CSS File": {
			validTemplateFile,
			filepath.Join(tmpDir, "nonexistent.css"),
			"http://golang.samba.br",
			false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewEmailTemplate(filepath.Dir(tt.templateFile), filepath.Base(tt.templateFile), tt.cssFile, tt.signatureLink)
			if err != nil && !tt.expectError {
				t.Errorf("NewEmailTemplate() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestEmailTemplate_Execute(t *testing.T) {
	tmpDir := t.TempDir()

	validTemplateContent := "Hello, {{.Data.Name}}!"
	validCSSContent := "body { color: red; }"

	validTemplateFile := createTempFile(t, tmpDir, "template-*.html", validTemplateContent)
	validCSSFile := createTempFile(t, tmpDir, "styles-*.css", validCSSContent)

	emailTemplate, err := NewEmailTemplate(filepath.Dir(validTemplateFile),
		filepath.Base(validTemplateFile), validCSSFile, "http://golang.samba.br")
	if err != nil {
		t.Fatalf("Failed to create EmailTemplate: %v", err)
	}

	tests := map[string]struct {
		data        map[string]string
		expected    string
		expectError bool
	}{
		"Valid Template and Data": {
			map[string]string{"Name": "Renê Cardozo"},
			"Hello, Renê Cardozo!",
			false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := emailTemplate.Execute(tt.data)
			if (err != nil) != tt.expectError {
				t.Errorf("Execute() error = %v, expectError %v", err, tt.expectError)
			}
			if result != tt.expected {
				t.Errorf("Execute() result = %v, expected %v", result, tt.expected)
			}
		})
	}
}
