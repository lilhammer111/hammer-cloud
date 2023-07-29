package handler

import (
	"github.com/lilhammer111/hammer-cloud/meta"
	"github.com/lilhammer111/hammer-cloud/util"
	"io"
	"net/http"
	"os"
	"time"
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
		// firstly, parse the file sent form frontend by form.
		// header , aka file meta info , includes filename, size, and httpHeader etc.
		f, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file from request", http.StatusInternalServerError)
			return
		}
		defer f.Close()
		// and then, save the header of the file
		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: "/tmp/" + header.Filename,
			UploadAt: time.Now().Format("2006/01/02 15:04:05"),
		}

		// Secondly, create a new file in order to save the file later.
		// Actually, Create() is just to invoke os.OpenFile() and returns a file struct pointer.
		nf, err := os.Create(fileMeta.Location)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer nf.Close()

		// Finally, copy file to new file.
		fileMeta.FileSize, err = io.Copy(nf, f)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		nf.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(nf)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Upload finished!")
	if err != nil {
		return
	}
}
