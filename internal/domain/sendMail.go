package domain

import (
	"doodocs-days/internal/service"
	"fmt"
	"net/http"
	"net/smtp"
	"os"

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
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	// Check if method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle file upload
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file
	isValid, err := h.SendMailService.ValidateFile(handler)
	if !isValid || err != nil {
		http.Error(w, "Incorrect type of file", http.StatusUnsupportedMediaType)
		return
	}

	// Handle emails
	emails, ok := r.MultipartForm.Value["emails"]
	if !ok || len(emails) == 0 {
		http.Error(w, "No email addresses provided", http.StatusBadRequest)
		return
	}

	fmt.Println(handler.Filename, emails)

	// Send email
	if err := sendEmail(emails, "hello", "test"); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func sendEmail(to []string, subject string, body string) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("FROM_EMAIL"),
		os.Getenv("FROM_EMAIL_PASSWORD"),
		os.Getenv("FROM_EMAIL_SMTP"),
	)

	// Format the email message
	message := "Subject: " + subject + "\n" + "Content-Type: text/plain; charset=UTF-8\n\n" + body

	return smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		to,
		[]byte(message),
	)
}
