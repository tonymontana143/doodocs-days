package repository

import (
	"errors"
	"mime/multipart"
)

var allowedMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func IsValidMimeType(file *multipart.FileHeader) (bool, error) {
	mimeType := file.Header.Get("Content-Type")
	if !allowedMimeTypes[mimeType] {
		return false, errors.New("unsupported file type")
	}
	return true, nil
}
