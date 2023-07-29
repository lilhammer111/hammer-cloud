package main

import (
	"github.com/lilhammer111/hammer-cloud/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("file/upload/suc", handler.UploadSucHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Println("Failed to start server", err)
	}
}
