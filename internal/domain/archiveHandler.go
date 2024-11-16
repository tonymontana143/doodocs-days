package domain

import (
	"archive/zip"
	"doodocs-days/internal/models"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type ZipHeader struct {
	Signature uint32
}

func ArchiveInfoHandlers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	file, handler, err := r.FormFile("zipfile")
	if err != nil {
		http.Error(w, "Failed to get file from form-data", http.StatusBadRequest)
		return
	}
	defer file.Close()
	data, err := getZipFileInfo(file, handler)
	if err != nil {
		http.Error(w, "Invalid zip file: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data)
}
func getZipFileInfo(file multipart.File, handler *multipart.FileHeader) (*models.ArchiveInfo, error) {
	isZip, err := isZipFile(file)
	if err != nil || !isZip {
		return nil, fmt.Errorf("error: %v", err)
	}
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
			MimeType: getMimeType(file.Name),
		}
		filesIncluded = append(filesIncluded, *tempFile)
		totalSize += float64(file.UncompressedSize64)

	}
	newItem := &models.ArchiveInfo{
		Filename:     handler.Filename,
		Archive_size: float64(handler.Size),
		Total_size:   totalSize,
		Total_files:  float64(len(filesIncluded)),
		Files:        filesIncluded,
	}
	return newItem, nil
}
func getMimeType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".txt":
		return "text/plain"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}
func isZipFile(file multipart.File) (bool, error) {
	var header ZipHeader
	err := binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return false, fmt.Errorf("failed to read file header: %v", err)
	}
	if header.Signature != 0x04034b50 {
		return false, err
	}
	return true, nil
}
