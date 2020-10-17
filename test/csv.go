package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("CSV_PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, http.FileServer(http.Dir("./"))))
}
