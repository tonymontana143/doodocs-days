package domain

import (
	"doodocs-days/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
)

// CreateArchiveHandler handles requests for creating archives
type CreateArchiveHandler struct {
	CreateArchiveService service.CreateArchiveService
}

// NewCreateArchiveHandler creates a new instance of CreateArchiveHandler
func NewCreateArchiveHandler(archiveService service.CreateArchiveService) *CreateArchiveHandler {
	return &CreateArchiveHandler{CreateArchiveService: archiveService}
}

// CreateArchive handles the creation of a ZIP archive
func (h *CreateArchiveHandler) CreateArchive(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request method and URL
	slog.Info("Handling CreateArchive request", "method", r.Method, "url", r.URL)

	// Method check
	if r.Method != "POST" {
		slog.Error("Invalid method", "method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "InvalidMethod",
			Message: "Method Not Allowed",
		})
		return
	}

	// Parse the multipart form with a limit of 10 MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error("Failed to parse multipart form data", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "FormParsingError",
			Message: "Unable to parse form data",
		})
		return
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		slog.Warn("No files provided in request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "NoFilesError",
			Message: "No files provided",
		})
		return
	}

	// Validate and zip the provided files
	zipBuffer, err := h.CreateArchiveService.ValidateAndZipFiles(files)
	if err != nil {
		slog.Error("Failed to validate and zip files", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "ZipCreationError",
			Message: err.Error(),
		})
		return
	}

	// Log successful creation of ZIP archive
	slog.Info("Successfully created ZIP archive", "file_count", len(files))

	// Set headers and write the ZIP file to the response
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")

	if _, err := w.Write(zipBuffer.Bytes()); err != nil {
		slog.Error("Failed to write ZIP archive to response", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "FileWriteError",
			Message: "Failed to send archive",
		})
	}
}
