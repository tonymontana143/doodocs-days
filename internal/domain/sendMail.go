package domain

import (
	"doodocs-days/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
)

// SendMailHandler handles requests for sending emails with attachments
type SendMailHandler struct {
	SendMailService service.SendMailService
}

// NewSendMailHandler creates a new instance of SendMailHandler
func NewSendMailHandler(svc service.SendMailService) *SendMailHandler {
	return &SendMailHandler{SendMailService: svc}
}

// SendMail handles sending emails with an attachment
func (h *SendMailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "EnvLoadingError",
			Message: "Error loading .env file",
		})
		return
	}

	// Log the incoming request method and URL
	slog.Info("Handling SendMail request", "method", r.Method, "url", r.URL)

	if r.Method != http.MethodPost {
		slog.Error("Invalid method", "method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "InvalidMethod",
			Message: "Method Not Allowed",
		})
		return
	}

	// Retrieve file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Failed to retrieve file from form data", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "FileRetrievalError",
			Message: "Failed to get file from form data",
		})
		return
	}
	defer file.Close()

	// Validate file type
	isValid, err := h.SendMailService.ValidateFile(handler)
	if err != nil {
		slog.Error("Error validating file", "filename", handler.Filename, "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "FileValidationError",
			Message: "File validation failed",
		})
		return
	}
	if !isValid {
		slog.Warn("Invalid file type", "filename", handler.Filename)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "InvalidFileType",
			Message: "Incorrect type of file",
		})
		return
	}

	// Retrieve email addresses
	emails, ok := r.MultipartForm.Value["emails"]
	if !ok || len(emails) == 0 {
		slog.Warn("No email addresses provided in request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "NoEmailsError",
			Message: "No email addresses provided",
		})
		return
	}

	// Log information about email sending
	slog.Info("Sending emails with attached file", "email_count", len(emails), "filename", handler.Filename)

	// Attempt to send emails
	files := r.MultipartForm.File["file"]
	if err := h.SendMailService.SendMails(emails, files); err != nil {
		slog.Error("Failed to send emails", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "EmailSendingError",
			Message: "Failed to send email: " + err.Error(),
		})
		return
	}

	// Log success
	slog.Info("File uploaded and email sent successfully", "email_count", len(emails), "filename", handler.Filename)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File uploaded and email sent successfully!",
	})
}
