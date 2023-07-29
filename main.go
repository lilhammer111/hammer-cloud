package main

import (
	"github.com/lilhammer111/hammer-cloud/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadFileHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Println("Failed to start server", err)
	}
}
