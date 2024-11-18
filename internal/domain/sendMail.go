package domain

import (
	"doodocs-days/internal/service"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
)

type SendMailHandler struct {
	SendMailService service.SendMailService
}

func NewSendMailHandler(svc service.SendMailService) *SendMailHandler {
	return &SendMailHandler{SendMailService: svc}
}

func (h *SendMailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	// Log the incoming request method and URL
	slog.Info("Handling SendMail request", "method", r.Method, "url", r.URL)

	if r.Method != http.MethodPost {
		slog.Error("Invalid method", "method", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Failed to retrieve file from form data", "error", err)
		http.Error(w, "Failed to get file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	isValid, err := h.SendMailService.ValidateFile(handler)
	if err != nil {
		slog.Error("Error validating file", "filename", handler.Filename, "error", err)
		http.Error(w, "File validation failed", http.StatusUnsupportedMediaType)
		return
	}
	if !isValid {
		slog.Warn("Invalid file type", "filename", handler.Filename)
		http.Error(w, "Incorrect type of file", http.StatusUnsupportedMediaType)
		return
	}

	// Retrieve email addresses
	emails, ok := r.MultipartForm.Value["emails"]
	if !ok || len(emails) == 0 {
		slog.Warn("No email addresses provided in request")
		http.Error(w, "No email addresses provided", http.StatusBadRequest)
		return
	}

	// Log information about email sending
	slog.Info("Sending emails with attached file", "email_count", len(emails), "filename", handler.Filename)

	// Attempt to send emails
	files := r.MultipartForm.File["file"]
	if err := h.SendMailService.SendMails(emails, files); err != nil {
		slog.Error("Failed to send emails", "error", err)
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log success
	slog.Info("File uploaded and email sent successfully", "email_count", len(emails), "filename", handler.Filename)

	w.Write([]byte("File uploaded and email sent successfully!"))
}
