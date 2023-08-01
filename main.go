package main

import (
	"github.com/lilhammer111/hammer-cloud/handler"
	"github.com/lilhammer111/hammer-cloud/middleware"
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

	http.HandleFunc("/user/signup", handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.LoginHandler)
	http.HandleFunc("/user/info", middleware.TokenAuthMDW(handler.UserInfoHandler))

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Println("Failed to start server", err)
	}
}
