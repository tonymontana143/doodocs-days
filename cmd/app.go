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

	//archive info
	archiveInfoService := service.NewArchiveService()
	archiveInfoHandler := handler.NewFileHandler(archiveInfoService)
	//create archive
	createArchiveService := service.NewCreateArchiveService()
	createArchiveHandler := handler.NewCreateArchiveHandler(createArchiveService)

	mux.HandleFunc("/api/archive/files", createArchiveHandler.CreateArchive)
	mux.HandleFunc("/api/archive/information", archiveInfoHandler.ArchiveInfoHandle)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error during running server")
	}

}
