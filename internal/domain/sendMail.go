package domain

import (
	"doodocs-days/internal/service"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type SendMailHandler struct {
	SendMailService service.SendMailService
}

func NewSendMailHandler(svc service.SendMailService) *SendMailHandler {
	return &SendMailHandler{SendMailService: svc}
}

func (h *SendMailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	isValid, err := h.SendMailService.ValidateFile(handler)
	if !isValid || err != nil {
		http.Error(w, "Incorrect type of file", http.StatusUnsupportedMediaType)
		return
	}

	emails, ok := r.MultipartForm.Value["emails"]
	if !ok || len(emails) == 0 {
		http.Error(w, "No email addresses provided", http.StatusBadRequest)
		return
	}

	if err := sendEmail(emails, handler.Filename, file); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func sendEmail(to []string, fileName string, file multipart.File) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("FROM_EMAIL"),
		os.Getenv("FROM_EMAIL_PASSWORD"),
		os.Getenv("FROM_EMAIL_SMTP"),
	)

	boundary := "my-boundary-123456"

	header := make(map[string]string)
	header["From"] = os.Getenv("FROM_EMAIL")
	header["To"] = strings.Join(to, ",")
	header["Subject"] = "Attached File"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf(`multipart/mixed; boundary="%s"`, boundary)

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n"

	message += fmt.Sprintf("--%s\r\n", boundary)
	message += "Content-Type: text/plain; charset=UTF-8\r\n\r\n"
	message += "Please find the attached file.\r\n\r\n"

	message += fmt.Sprintf("--%s\r\n", boundary)
	message += "Content-Type: application/octet-stream\r\n"
	message += fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName)
	message += "Content-Transfer-Encoding: base64\r\n\r\n"

	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		message += base64.StdEncoding.EncodeToString(buffer[:n])
		message += "\r\n"
	}

	message += fmt.Sprintf("--%s--\r\n", boundary)

	return smtp.SendMail(

		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		to,
		[]byte(message),
	)
}
