package domain

import (
	"doodocs-days/internal/service"
	"log/slog"
	"net/http"
)

type CreateArchiveHandler struct {
	CreateArchiveService service.CreateArchiveService
}

func NewCreateArchiveHandler(archiveService service.CreateArchiveService) *CreateArchiveHandler {
	return &CreateArchiveHandler{CreateArchiveService: archiveService}
}

func (h *CreateArchiveHandler) CreateArchive(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request method and URL
	slog.Info("Handling CreateArchive request", "method", r.Method, "url", r.URL)

	if r.Method != "POST" {
		slog.Error("Invalid method", "method", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form with a limit of 10 MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error("Failed to parse multipart form data", "error", err)
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		slog.Warn("No files provided in request")
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	// Validate and zip the provided files
	zipBuffer, err := h.CreateArchiveService.ValidateAndZipFiles(files)
	if err != nil {
		slog.Error("Failed to validate and zip files", "error", err)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// Log successful creation of ZIP archive
	slog.Info("Successfully created ZIP archive", "file_count", len(files))

	// Set headers and write the ZIP file to the response
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")

	if _, err := w.Write(zipBuffer.Bytes()); err != nil {
		slog.Error("Failed to write ZIP archive to response", "error", err)
		http.Error(w, "Failed to send archive", http.StatusInternalServerError)
	}
}
