package domain

import (
	"doodocs-days/internal/service"
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
	if err := godotenv.Load(); err != nil {
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	files := r.MultipartForm.File["file"]
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

	if err := h.SendMailService.SendMails(emails, files); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File uploaded and email sent successfully!"))
}

// func sendEmail(to []string, fileName string, file multipart.File) error {
// 	fmt.Println("Sending email to:", to)

// 	auth := smtp.PlainAuth(
// 		"",
// 		os.Getenv("FROM_EMAIL"),
// 		os.Getenv("FROM_EMAIL_PASSWORD"),
// 		os.Getenv("FROM_EMAIL_SMTP"),
// 	)

// 	boundary := "my-boundary-123456"

// 	header := make(map[string]string)
// 	header["From"] = os.Getenv("FROM_EMAIL")
// 	header["To"] = strings.Join(to, ",")
// 	header["Subject"] = "Attached File"
// 	header["MIME-Version"] = "1.0"
// 	header["Content-Type"] = fmt.Sprintf(`multipart/mixed; boundary="%s"`, boundary)

// 	var message strings.Builder
// 	for k, v := range header {
// 		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
// 	}
// 	message.WriteString("\r\n")

// 	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
// 	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
// 	message.WriteString("Please find the attached file.\r\n\r\n")

// 	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
// 	message.WriteString("Content-Type: application/octet-stream\r\n")
// 	message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))
// 	message.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")

// 	buffer := make([]byte, 4096)
// 	for {
// 		n, err := file.Read(buffer)
// 		if err != nil && err != io.EOF {
// 			// Log error during file read
// 			fmt.Println("Error reading file:", err)
// 			return err
// 		}
// 		if n == 0 {
// 			break
// 		}
// 		encoded := base64.StdEncoding.EncodeToString(buffer[:n])
// 		message.WriteString(encoded + "\r\n")
// 	}

// 	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

// 	err := smtp.SendMail(
// 		os.Getenv("SMTP_ADDR"),
// 		auth,
// 		os.Getenv("FROM_EMAIL"),
// 		to,
// 		[]byte(message.String()),
// 	)
// 	if err != nil {
// 		fmt.Println("Failed to send email:", err)
// 	}
// 	return err
// }
