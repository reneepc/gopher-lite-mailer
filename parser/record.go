package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type MailRecord struct {
	Email string
	Data  map[string]string
}

func ParseRecords(path string) ([]MailRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("file must have at least two rows: header and data")
	}

	headers := rows[0]
	records := make([]MailRecord, 0, len(rows)-1)
	for _, row := range rows[1:] {
		record := MailRecord{
			Email: row[0],
			Data:  make(map[string]string),
		}
		for i, value := range row {
			record.Data[strings.TrimSpace(headers[i])] = strings.TrimSpace(value)
		}
		records = append(records, record)
	}

	return records, nil
}
