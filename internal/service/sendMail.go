package service

import (
	"doodocs-days/internal/config"
	"doodocs-days/internal/repository"
	"errors"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

type SendMailService interface {
	ValidateFile(file *multipart.FileHeader) (bool, error)
	SendMails(emails []string, files []*multipart.FileHeader) error
}

type SendMail struct {
	conf config.MailConfig
}

func NewSendMailService(conf config.MailConfig) SendMailService {
	return &SendMail{conf: conf}
}

func (s *SendMail) ValidateFile(file *multipart.FileHeader) (bool, error) {
	return repository.IsValidMimeTypeForMail(file)
}

func (s *SendMail) SendMails(emailsOne []string, files []*multipart.FileHeader) error {
	if len(emailsOne) == 0 {
		return errors.New("no recipients specified")
	}
	emails := strings.Split(emailsOne[0], ",")
	e := email.NewEmail()
	e.From = s.conf.EmailSenderAddress
	e.To = emails
	e.Subject = "Your Subject Here"
	e.Text = []byte("Please find the attached files.")

	for _, file := range files {
		fileContent, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file.Filename, err)
		}
		defer fileContent.Close()

		_, err = e.Attach(fileContent, file.Filename, file.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", file.Filename, err)
		}
	}

	// Set up SMTP authentication
	auth := smtp.PlainAuth(
		"",
		s.conf.EmailSenderAddress,
		s.conf.EmailSenderPassword,
		s.conf.Host,
	)

	// Send the email
	serverAddress := fmt.Sprintf("%s:%s", s.conf.Host, s.conf.Port)

	err := e.Send(serverAddress, auth)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
