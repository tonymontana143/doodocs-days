package main

import (
	handler "doodocs-days/internal/domain"
	"doodocs-days/internal/service"
	"fmt"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	archiveService := service.NewArchiveService()
	archiveHandler := handler.NewFileHandler(archiveService)
	mux.HandleFunc("/api/archive/information", archiveHandler.ArchiveInfoHandle)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error during running server")
	}

}
