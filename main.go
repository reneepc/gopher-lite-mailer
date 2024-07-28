package main

import (
	"flag"
	"log/slog"
	"os"
	"path"
	"sync"

	"github.com/reneepc/gopher-lite-mailer/mailer"
	"github.com/reneepc/gopher-lite-mailer/parser"
)

func main() {
	templateSubDir := flag.String("dir", "standard", "Subdirectory containing the template files")
	bodyFile := flag.String("body", "workshop-confirmation.html", "Body template file to use")
	dataFile := flag.String("data", "data.csv", "Data file to use (should be in the data subdirectory of the template directory)")
	signatureLink := flag.String("signature", "https://golang.sampa.br/img/golangsp01.png", "Signature link to use for the email body")
	subject := flag.String("subject", "", "Subject of the email")

	flag.Parse()

	if len(flag.Args()) < 2 {
		slog.Error("email and password are required")
		slog.Error("Usage: gopher-lite-mailer [options] <email> <password>")
		slog.Error("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	args := flag.Args()
	email, password := args[0], args[1]

	templateDir := path.Join("templates", *templateSubDir)
	templateContent, err := mailer.NewEmailTemplate(templateDir, *bodyFile, *signatureLink)
	if err != nil {
		slog.Error("could not create email template: %v", slog.Any("error", err))
		return
	}

	dataFilePath := path.Join(templateDir, "data", *dataFile)
	mailContent, err := parser.ParseRecords(dataFilePath)
	if err != nil {
		slog.Error("could not parse CSV file: %v", slog.Any("error", err))
		return
	}

	emailMailer := mailer.NewGMailMailerBuilder(email, password).Build()

	sendEmails(emailMailer, *subject, templateContent, mailContent)
}

func sendEmails(mailer mailer.Mailer, subject string, template mailer.EmailTemplate, records []parser.MailRecord) {
	var wg sync.WaitGroup

	for _, mailRecord := range records {
		wg.Add(1)
		go func(record parser.MailRecord) {
			// time.Sleep(time.Duration(rand.IntN(30)) * time.Second)
			defer wg.Done()

			body, err := template.Execute(record.Data)
			if err != nil {
				slog.Error("could not execute template: %v", slog.Any("error", err))
				return
			}

			err = mailer.SendMail(record.Email, subject, body)
			if err != nil {
				slog.Error("could not send email to %s: %v", slog.String("email", record.Email), slog.Any("error", err))
			}
		}(mailRecord)
	}

	wg.Wait()
}
