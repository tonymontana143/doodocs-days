package domain

import (
	"doodocs-days/internal/service"
	"fmt"
	"net/http"
)

type SendMailHandler struct {
	SendMailService service.SendMailService
}

func NewSendMailHandler(svc service.SendMailService) *SendMailHandler {
	return &SendMailHandler{SendMailService: svc}
}
func (h *SendMailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("correct")
}
