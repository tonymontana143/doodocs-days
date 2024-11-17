package repository

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
)

// CheckZipMagicNumber verifies the first 4 bytes for ZIP's magic number.
func CheckZipMagicNumber(file multipart.File) (bool, error) {
	file.Seek(0, io.SeekStart)

	var header [4]byte
	if _, err := file.Read(header[:]); err != nil {
		return false, fmt.Errorf("failed to read file header: %v", err)
	}

	if header != [4]byte{0x50, 0x4B, 0x03, 0x04} {
		return false, errors.New("file is not a valid zip archive")
	}
	return true, nil
}

// GetMimeTypeFromExtension returns the MIME type based on the file extension.
func GetMimeTypeFromExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream" // Fallback for unknown types
	}
	return mimeType
}
