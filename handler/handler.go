package handler

import (
	"encoding/json"
	"errors"
	"github.com/lilhammer111/hammer-cloud/meta"
	"github.com/lilhammer111/hammer-cloud/util"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		html, err := os.ReadFile("./static/view/upload.html")
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
		//meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Upload finished!")
	if err != nil {
		return
	}
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			// judge if the error is caused by bad request.
			var er *url.Error
			if errors.As(err, &er) {
				http.Error(w, "Wrong request format", http.StatusBadRequest)
			} else {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
		filehash := r.Form["filehash"][0]
		fMeta, err := meta.GetFileMetaDB(filehash)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(fMeta)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(data)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	fsha1 := r.Form.Get("filehash")
	fm, err := meta.GetFileMetaDB(fsha1)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	f, err := os.Open(fm.Location)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")

	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	opType := r.Form.Get("op")

	if opType != "0" {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	curFileMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	curFileMeta.FileName = newFileName
	//meta.UpdateFileMeta(curFileMeta)
	_ = meta.UpdateFileMetaDB(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		log.Println()
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	fileSha1 := r.Form.Get("filehash")
	fMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// delete file by path
	err = os.Remove(fMeta.Location)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// delete meta info
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}
