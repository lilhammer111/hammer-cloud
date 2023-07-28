package handler

import (
	"io"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		html, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internal Server")
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(html)
		if err != nil {
			return
		}
	} else if r.Method == "POST" {

	}
}
