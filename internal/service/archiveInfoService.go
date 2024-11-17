package service

import (
	"archive/zip"
	"doodocs-days/internal/models"
	"doodocs-days/internal/repository"
	"fmt"
	"io"
	"mime/multipart"
)

type ArchiveInfoService interface {
	GetZipFileInfo(file multipart.File, handler *multipart.FileHeader) (*models.ArchiveInfo, error)
	IsZipFile(file multipart.File) (bool, error)
}
type ArchiveService struct{}

func NewArchiveService() *ArchiveService {
	return &ArchiveService{}
}
func (s *ArchiveService) GetZipFileInfo(file multipart.File, handler *multipart.FileHeader) (*models.ArchiveInfo, error) {
	file.Seek(0, io.SeekStart)
	zipReader, err := zip.NewReader(file, handler.Size)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	filesIncluded := make([]models.ObjectFile, 0)
	var totalSize float64

	for _, file := range zipReader.File {
		tempFile := &models.ObjectFile{
			FilePath: file.Name,
			Size:     float64(file.CompressedSize64),
			MimeType: repository.GetMimeTypeFromExtension(file.Name),
		}
		filesIncluded = append(filesIncluded, *tempFile)
		totalSize += float64(file.UncompressedSize64)

	}
	newItem := &models.ArchiveInfo{
		Filename:     handler.Filename,
		Archive_size: float32(handler.Size),
		Total_size:   float32(totalSize),
		Total_files:  float32(len(filesIncluded)),
		Files:        filesIncluded,
	}
	return newItem, nil

}
func (s *ArchiveService) IsZipFile(file multipart.File) (bool, error) {
	return repository.CheckZipMagicNumber(file)
}
