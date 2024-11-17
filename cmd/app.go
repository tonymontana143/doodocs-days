package main

import (
	handler "doodocs-days/internal/domain"
	"doodocs-days/internal/service"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Create a new HTTP ServeMux for routing
	mux := http.NewServeMux()

	// Initialize ArchiveInfoService and its handler
	archiveInfoService := service.NewArchiveService()
	archiveInfoHandler := handler.NewFileHandler(archiveInfoService)

	// Initialize CreateArchiveService and its handler
	createArchiveService := service.NewCreateArchiveService()
	createArchiveHandler := handler.NewCreateArchiveHandler(createArchiveService)

	// Initialize SendMailService and its handler
	sendMailService := service.NewSendMailService()
	sendMailHandler := handler.NewSendMailHandler(sendMailService)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/api/mail/file", sendMailHandler.SendMail)
	mux.HandleFunc("/api/archive/files", createArchiveHandler.CreateArchive)
	mux.HandleFunc("/api/archive/information", archiveInfoHandler.ArchiveInfoHandle)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	fmt.Println("Server started on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}
