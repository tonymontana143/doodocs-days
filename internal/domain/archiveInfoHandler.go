package domain

import (
	"doodocs-days/internal/service"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type ArchiveInfoHandler struct {
	ArchiveService service.ArchiveInfoService
}

func NewFileHandler(fileService service.ArchiveInfoService) *ArchiveInfoHandler {
	return &ArchiveInfoHandler{ArchiveService: fileService}
}

func (h *ArchiveInfoHandler) ArchiveInfoHandle(w http.ResponseWriter, r *http.Request) {
	// Log request method and URL
	slog.Info("Handling ArchiveInfo request", "method", r.Method, "url", r.URL)

	if r.Method != "POST" {
		slog.Error("Invalid method", "method", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Failed to retrieve file from form data", "error", err)
		http.Error(w, "Failed to get file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	isZip, err := h.ArchiveService.IsZipFile(file)
	if err != nil {
		slog.Error("Error checking if file is a zip", "error", err)
		http.Error(w, "Failed to check file type", http.StatusInternalServerError)
		return
	}
	if !isZip {
		slog.Warn("File is not a valid zip file", "filename", handler.Filename)
		http.Error(w, "Invalid zip file", http.StatusUnsupportedMediaType)
		return
	}

	// Reset the file pointer before reading again
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		slog.Error("Failed to reset file pointer", "error", err)
		http.Error(w, "Failed to process ZIP file", http.StatusInternalServerError)
		return
	}

	archiveInfo, err := h.ArchiveService.GetZipFileInfo(file, handler)
	if err != nil {
		slog.Error("Failed to process ZIP file", "error", err)
		http.Error(w, "Failed to process ZIP file", http.StatusInternalServerError)
		return
	}

	// Log success in processing
	slog.Info("Successfully processed ZIP file", "filename", handler.Filename)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(archiveInfo); err != nil {
		slog.Error("Failed to encode archive info as JSON", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
