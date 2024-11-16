package main

import (
	handler "doodocs-days/internal/domain"
	"fmt"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/archive/information", handler.ArchiveInfoHandlers)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error during running server")
	}

}
