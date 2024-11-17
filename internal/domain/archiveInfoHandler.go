package domain

import (
	"doodocs-days/internal/service"
	"encoding/json"
	"io"
	"net/http"
)

type ArchiveInfoHandler struct {
	ArchiveService service.ArchiveInfoService
}

func NewFileHandler(fileService service.ArchiveInfoService) *ArchiveInfoHandler {
	return &ArchiveInfoHandler{ArchiveService: fileService}
}
func (h *ArchiveInfoHandler) ArchiveInfoHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()
	isZip, err := h.ArchiveService.IsZipFile(file)
	if !isZip || err != nil {
		http.Error(w, "Invalid zip file", http.StatusUnsupportedMediaType)
		return
	}
	file.Seek(0, io.SeekStart)
	archiveInfo, err := h.ArchiveService.GetZipFileInfo(file, handler)
	if err != nil {
		http.Error(w, "Failed to process ZIP file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(archiveInfo)
}
