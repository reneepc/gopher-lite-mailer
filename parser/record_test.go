package parser_test

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"

	"github.com/reneepc/gopher-lite-mailer/parser"
)

func createTempCSVFile(t *testing.T, content [][]string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)
	if err := writer.WriteAll(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	writer.Flush()

	return tmpFile.Name()
}

func TestParseRecords(t *testing.T) {
	tests := map[string]struct {
		content     [][]string
		expected    []parser.MailRecord
		expectError bool
	}{
		"Valid CSV": {
			content: [][]string{
				{"Email", "Name", "Age"},
				{"rene.epcrdz@gmail.com", "Renê Cardozo", "25"},
				{"jorge@example.com", "Jorge", "67"},
			},
			expected: []parser.MailRecord{
				{
					Email: "rene.epcrdz@gmail.com",
					Data: map[string]string{
						"Email": "rene.epcrdz@gmail.com",
						"Name":  "Renê Cardozo",
						"Age":   "25",
					},
				},
				{
					Email: "jorge@example.com",
					Data: map[string]string{
						"Email": "jorge@example.com",
						"Name":  "Jorge",
						"Age":   "67",
					},
				},
			},
			expectError: false,
		},
		"Empty CSV": {
			content:     [][]string{},
			expected:    nil,
			expectError: true,
		},
		"Only Header": {
			content:     [][]string{{"Email", "Name", "Age"}},
			expected:    nil,
			expectError: true,
		},
		"Missing Email": {
			content: [][]string{
				{"Email", "Name", "Age"},
				{"", "No Email", "40"},
			},
			expected: []parser.MailRecord{
				{
					Email: "",
					Data: map[string]string{
						"Email": "",
						"Name":  "No Email",
						"Age":   "40",
					},
				},
			},
			expectError: false,
		},
		"File Open Error": {
			content:     nil,
			expected:    nil,
			expectError: true,
		},
		"Read All Error": {
			content:     [][]string{{"Email", "Name"}, {"rene.epcrdz@gmail.com", "Renê Cardozo"}},
			expected:    nil,
			expectError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var tmpFilePath string
			if name == "File Open Error" {
				tmpFilePath = "non_existent_file.csv"
			} else {
				tmpFilePath = createTempCSVFile(t, tt.content)

				if name == "Read All Error" {
					corruptFileForReadAll(t, tmpFilePath)
				}
			}
			defer os.Remove(tmpFilePath)

			records, err := parser.ParseRecords(tmpFilePath)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error = %v, got %v", tt.expectError, err)
			}

			if !reflect.DeepEqual(records, tt.expected) {
				t.Errorf("expected records = %v, got %v", tt.expected, records)
			}
		})
	}
}

func corruptFileForReadAll(t *testing.T, tmpFilePath string) {
	t.Helper()

	file, err := os.OpenFile(tmpFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatalf("Failed to open file for corruption: %v", err)
	}
	file.WriteString("Email,Name\n\"rene.epcrdz@gmail.com\"\"Renê Cardozo\n")
	file.Close()
}
