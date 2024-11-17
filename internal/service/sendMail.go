package service

import (
	"doodocs-days/internal/repository"
	"mime/multipart"
)

type SendMailService interface {
	ValidateFile(file *multipart.FileHeader) (bool, error)
	SendMails(emails []string, file *multipart.File) error
}
type sendMail struct{}

func NewSendMailService() SendMailService {
	return &sendMail{}
}
func (s *sendMail) ValidateFile(file *multipart.FileHeader) (bool, error) {
	return repository.IsValidMimeTypeForMail(file)
}
func (s *sendMail) SendMails(emails []string, file *multipart.File) error {
	return nil
}
