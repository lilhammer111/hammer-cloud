package main

import (
	h "github.com/lilhammer111/hammer-cloud/handler"
	m "github.com/lilhammer111/hammer-cloud/middleware"
	"log"
	"net/http"
)

func main() {
	// user api
	http.HandleFunc("/user/signup", h.SignUpHandler)
	http.HandleFunc("/user/signin", h.LoginHandler)
	http.HandleFunc("/user/info", m.TokenAuth(h.UserInfoHandler))
	// file api
	http.HandleFunc("/file/upload", m.TokenAuth(h.UploadHandler))
	http.HandleFunc("/file/upload/suc", m.TokenAuth(h.UploadSucHandler))
	http.HandleFunc("/file/meta", m.TokenAuth(h.GetFileMetaHandler))
	http.HandleFunc("/file/query", m.TokenAuth(h.FileQueryHandler))
	http.HandleFunc("/file/download", m.TokenAuth(h.DownloadFileHandler))
	http.HandleFunc("/file/update", m.TokenAuth(h.FileMetaUpdateHandler))
	http.HandleFunc("/file/delete", m.TokenAuth(h.FileDeleteHandler))
	//fast upload api
	http.HandleFunc("/file/fastupload", m.TokenAuth(h.TryFastUploadHandler))

	// multi part upload api
	http.HandleFunc("/file/mpupload/init", m.TokenAuth(h.InitMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", m.TokenAuth(h.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", m.TokenAuth(h.CompleteUploadHandler))

	// oss api
	http.HandleFunc("/file/downloadurl", m.TokenAuth(h.DownloadURLHandler))

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Println("Failed to start server", err)
	}
}
