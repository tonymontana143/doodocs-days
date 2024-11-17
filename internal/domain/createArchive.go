package domain

import (
	"doodocs-days/internal/service"
	"net/http"
)

type CreateArchiveHandler struct {
	CreateArchiveService service.CreateArchiveService
}

func NewCreateArchiveHandler(archiveService service.CreateArchiveService) *CreateArchiveHandler {
	return &CreateArchiveHandler{CreateArchiveService: archiveService}
}

func (h *CreateArchiveHandler) CreateArchive(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}
	zipBuffer, err := h.CreateArchiveService.ValidateAndZipFiles(files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")

	w.Write(zipBuffer.Bytes())
}
