package repository

import (
	"errors"
	"mime/multipart"
)

var allowedMimeTypesForMail = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/pdf": true,
}

func IsValidMimeTypeForMail(file *multipart.FileHeader) (bool, error) {
	mimeType := file.Header.Get("Content-Type")
	if !allowedMimeTypesForMail[mimeType] {
		return false, errors.New("unsupported file type")
	}
	return true, nil
}
