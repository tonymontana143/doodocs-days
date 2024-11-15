package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error during running server")
	}
}
