package service

import (
	"archive/zip"
	"bytes"
	"doodocs-days/internal/repository"
	"fmt"
	"io"
	"mime/multipart"
)

type CreateArchiveService interface {
	ValidateAndZipFiles(files []*multipart.FileHeader) (*bytes.Buffer, error)
}
type createArchive struct{}

func NewCreateArchiveService() CreateArchiveService {
	return &createArchive{}
}
func (s *createArchive) ValidateAndZipFiles(files []*multipart.FileHeader) (*bytes.Buffer, error) {
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	for _, fileHeader := range files {
		isValid, err := repository.IsValidMimeType(fileHeader)
		if err != nil || !isValid {
			return nil, fmt.Errorf("file %s has unsupported MIME type", fileHeader.Filename)
		}

		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("could not open file %s", fileHeader.Filename)
		}

		zipFileWriter, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			file.Close() // Close immediately if there's an error
			return nil, fmt.Errorf("failed to add file %s to archive", fileHeader.Filename)
		}

		_, err = io.Copy(zipFileWriter, file)
		file.Close() // Close immediately after copying
		if err != nil {
			return nil, fmt.Errorf("failed to write file %s to archive", fileHeader.Filename)
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to finalize ZIP archive")
	}

	return &zipBuffer, nil
}
