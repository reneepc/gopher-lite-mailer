package mailer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/reneepc/gopher-lite-mailer/mailer"
)

func createTempFile(t *testing.T, dir, filename, content string) string {
	t.Helper()

	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	return filePath
}

func TestNewEmailTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid template files
	validHeaderContent := "<div class='header'>Header</div>"
	validFooterContent := "<div class='footer'>Footer</div>"
	validBodyContent := "<p>Hello, {{.Data.Name}}!</p>"
	validCssContent := "body { color: red; }"

	// Create the required directories and files
	createTempFile(t, tmpDir, "header.html", validHeaderContent)
	createTempFile(t, tmpDir, "footer.html", validFooterContent)
	createTempFile(t, filepath.Join(tmpDir, "bodies"), "body1.html", validBodyContent)
	createTempFile(t, tmpDir, "styles.css", validCssContent)

	// Create an invalid body file for testing
	invalidBodyContent := "{{.Invalid"
	createTempFile(t, filepath.Join(tmpDir, "bodies"), "body_invalid.html", invalidBodyContent)

	tests := map[string]struct {
		bodyFile      string
		signatureLink string
		expectError   bool
		beforeTest    func()
	}{
		"Valid Files": {
			bodyFile:      "body1.html",
			signatureLink: "http://golang.samba.br",
			expectError:   false,
			beforeTest:    func() {},
		},
		"Invalid Body File": {
			bodyFile:      "body_invalid.html",
			signatureLink: "http://golang.samba.br",
			expectError:   true,
			beforeTest:    func() {},
		},
		"Missing CSS File": {
			bodyFile:      "body1.html",
			signatureLink: "http://golang.samba.br",
			expectError:   false, // Should not fail completely, only log a warning
			beforeTest: func() {
				// Remove the CSS file for this test case
				err := os.Remove(filepath.Join(tmpDir, "styles.css"))
				if err != nil {
					t.Fatalf("Failed to remove styles.css: %v", err)
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.beforeTest()

			_, err := mailer.NewEmailTemplate(
				tmpDir,
				tt.bodyFile,
				tt.signatureLink,
			)
			if err != nil && !tt.expectError {
				t.Errorf("NewEmailTemplate() error = %v, expectError %v", err, tt.expectError)
			}
			if err == nil && tt.expectError {
				t.Errorf("Expected error but got none")
			}
		})
	}
}

func TestEmailTemplate_Execute(t *testing.T) {
	validHeaderContent := "<div class='header'>Header</div>"
	validFooterContent := "<div class='footer'>Footer</div>"
	validBodyContent := "<p>Hello, {{.Data.Name}}!</p>"
	validCssContent := "body { color: red; }"

	tests := map[string]struct {
		setupFunc   func(tmpDir, bodyDir string)
		data        map[string]string
		expected    string
		expectError bool
	}{
		"Valid Template and Data": {
			setupFunc: func(tmpDir, bodyDir string) {
				createTempFile(t, tmpDir, "header.html", validHeaderContent)
				createTempFile(t, tmpDir, "footer.html", validFooterContent)
				createTempFile(t, bodyDir, "body1.html", validBodyContent)
				createTempFile(t, tmpDir, "styles.css", validCssContent)
			},
			data:        map[string]string{"Name": "Renê Cardozo"},
			expected:    "<div class='header'>Header</div><p>Hello, Renê Cardozo!</p><div class='footer'>Footer</div>",
			expectError: false,
		},
		"Missing Header File": {
			setupFunc: func(tmpDir, bodyDir string) {
				createTempFile(t, tmpDir, "footer.html", validFooterContent)
				createTempFile(t, bodyDir, "body1.html", validBodyContent)
				createTempFile(t, tmpDir, "styles.css", validCssContent)
			},
			data:        map[string]string{"Name": "Renê Cardozo"},
			expected:    "",
			expectError: true,
		},
		"Missing Footer File": {
			setupFunc: func(tmpDir, bodyDir string) {
				createTempFile(t, tmpDir, "header.html", validHeaderContent)
				createTempFile(t, bodyDir, "body1.html", validBodyContent)
				createTempFile(t, tmpDir, "styles.css", validCssContent)
			},
			data:        map[string]string{"Name": "Renê Cardozo"},
			expected:    "",
			expectError: true,
		},
		"Missing Body File": {
			setupFunc: func(tmpDir, bodyDir string) {
				createTempFile(t, tmpDir, "header.html", validHeaderContent)
				createTempFile(t, tmpDir, "footer.html", validFooterContent)
				createTempFile(t, tmpDir, "styles.css", validCssContent)
			},
			data:        map[string]string{"Name": "Renê Cardozo"},
			expected:    "",
			expectError: true,
		},
		"Missing CSS File": {
			setupFunc: func(tmpDir, bodyDir string) {
				createTempFile(t, tmpDir, "header.html", validHeaderContent)
				createTempFile(t, tmpDir, "footer.html", validFooterContent)
				createTempFile(t, bodyDir, "body1.html", validBodyContent)
			},
			data:        map[string]string{"Name": "Renê Cardozo"},
			expected:    "<div class='header'>Header</div><p>Hello, Renê Cardozo!</p><div class='footer'>Footer</div>",
			expectError: false, // Only logs a warning, should not fail completely
		},
		"Invalid Template Data": {
			setupFunc: func(tmpDir, bodyDir string) {
				createTempFile(t, tmpDir, "header.html", validHeaderContent)
				createTempFile(t, tmpDir, "footer.html", validFooterContent)
				createTempFile(t, bodyDir, "body1.html", validBodyContent)
				createTempFile(t, tmpDir, "styles.css", validCssContent)
			},
			data:        map[string]string{},
			expected:    "<div class='header'>Header</div><p>Hello, !</p><div class='footer'>Footer</div>",
			expectError: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tmpDir := t.TempDir()
			bodyDir := filepath.Join(tmpDir, "bodies")

			tt.setupFunc(tmpDir, bodyDir)

			emailTemplate, err := mailer.NewEmailTemplate(tmpDir, "body1.html", "http://golang.samba.br")
			if err != nil && !tt.expectError {
				t.Fatalf("Failed to create EmailTemplate: %v", err)
			}
			if err == nil && tt.expectError {
				t.Fatalf("Expected error but got none")
			}
			if tt.expectError {
				return
			}

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
