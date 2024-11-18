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

// ErrorResponse is the structure for all error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (h *ArchiveInfoHandler) ArchiveInfoHandle(w http.ResponseWriter, r *http.Request) {
	// Log request method and URL
	slog.Info("Handling ArchiveInfo request", "method", r.Method, "url", r.URL)

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

	// File retrieval
	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Failed to retrieve file from form data", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "InvalidFile",
			Message: "Failed to get file from form data",
		})
		return
	}
	defer file.Close()

	// Check if the file is a valid zip
	isZip, err := h.ArchiveService.IsZipFile(file)
	if err != nil {
		slog.Error("Error checking if file is a zip", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "FileCheckError",
			Message: "Failed to check file type",
		})
		return
	}
	if !isZip {
		slog.Warn("File is not a valid zip file", "filename", handler.Filename)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "InvalidZipFile",
			Message: "File is not a valid zip file",
		})
		return
	}

	// Reset the file pointer before reading again
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		slog.Error("Failed to reset file pointer", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "FilePointerError",
			Message: "Failed to process ZIP file",
		})
		return
	}

	// Get archive info
	archiveInfo, err := h.ArchiveService.GetZipFileInfo(file, handler)
	if err != nil {
		slog.Error("Failed to process ZIP file", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "ArchiveProcessingError",
			Message: "Failed to process ZIP file",
		})
		return
	}

	// Log success in processing
	slog.Info("Successfully processed ZIP file", "filename", handler.Filename)

	// Send the archive info as response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(archiveInfo); err != nil {
		slog.Error("Failed to encode archive info as JSON", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "JSONEncodingError",
			Message: "Failed to encode response",
		})
		return
	}
}
