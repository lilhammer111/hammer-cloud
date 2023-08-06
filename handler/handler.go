package handler

import (
	"encoding/json"
	"errors"
	"github.com/lilhammer111/hammer-cloud/common"
	"github.com/lilhammer111/hammer-cloud/config"
	"github.com/lilhammer111/hammer-cloud/db"
	"github.com/lilhammer111/hammer-cloud/meta"
	"github.com/lilhammer111/hammer-cloud/mq"
	"github.com/lilhammer111/hammer-cloud/store/oss"
	"github.com/lilhammer111/hammer-cloud/util"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
	}

	if r.Method == http.MethodPost {
		// firstly, parse the file sent form frontend by form.
		// header , aka file meta info , includes filename, size, and httpHeader etc.

		username := r.FormValue("username")
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

		_, err = nf.Seek(0, 0)
		if err != nil {
			http.Error(w, ErrSeekError, http.StatusInternalServerError)
			return
		}
		fileMeta.FileSha1 = util.FileSha1(nf)

		// at the meanwhile, to save the file into oss
		_, err = nf.Seek(0, 0)
		if err != nil {
			http.Error(w, ErrSeekError, http.StatusInternalServerError)
			return
		}
		ossPath := OSSRootDir + fileMeta.FileSha1
		//err = oss.Bucket().PutObject(ossPath, nf)
		//if err != nil {
		//	http.Error(w, ErrOSSPutError, http.StatusInternalServerError)
		//	return
		//}
		//fileMeta.Location = ossPath
		data := mq.TransferData{
			FileHash:      fileMeta.FileSha1,
			CurLocation:   fileMeta.Location,
			DestLocation:  ossPath,
			DestStoreType: common.StoreOSS,
		}
		pubData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, ErrMarshalError, http.StatusInternalServerError)
			return
		}
		suc := mq.Publish(config.TransExchangeName, config.TransOSSRoutingKey, pubData)
		if !suc {
			// todo: republish msg
		}

		//meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)

		// TODO : updates user file table
		//err = rd.ParseForm()
		//if err != nil {
		//	http.Error(w, "wrong parameter", http.StatusBadRequest)
		//	return
		//}
		//username := rd.Form.Get("username")
		suc = db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
		} else {
			http.Error(w, "upload failed", http.StatusInternalServerError)
		}
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

func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusInternalServerError)
		return
	}
	limStr := r.Form.Get("limit")
	username := r.Form.Get("username")
	lim, err := strconv.Atoi(limStr)
	if err != nil {
		http.Error(w, "failed to convert limit", http.StatusInternalServerError)
		return

	}
	userFiles, err := db.QueryUserFileMetas(username, lim)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		http.Error(w, "marshal failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "resp error", http.StatusInternalServerError)
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

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("username")
	fileHash := r.Form.Get("fileHash")
	filename := r.Form.Get("filename")
	filesizeS := r.Form.Get("filesize")
	filesize, err := strconv.ParseInt(filesizeS, 10, 64)
	if err != nil {
		http.Error(w, "str convert err", http.StatusInternalServerError)
		return
	}

	fileMeta, err := db.GetFileMeta(fileHash)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if fileMeta == nil {
		resp := util.RespBody{
			Code: -1,
			Msg:  "failed to fast upload",
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(resp.JSONBytes())
		if err != nil {
			http.Error(w, "resp error", http.StatusInternalServerError)
			return
		}
		return
	}

	suc := db.OnUserFileUploadFinished(username, fileHash, filename, filesize)

	w.Header().Set("Content-Type", "application/json")

	if suc {
		resp := util.RespBody{
			Code: 0,
			Msg:  "fast upload success",
		}

		_, err := w.Write(resp.JSONBytes())
		if err != nil {
			http.Error(w, "resp error ", http.StatusInternalServerError)
			return
		}
	} else {
		resp := util.NewRespBody(-2, "fast upload fail", nil)
		_, err := w.Write(resp.JSONBytes())
		if err != nil {
			http.Error(w, "resp error ", http.StatusInternalServerError)
			return
		}
	}
}

func DownloadURLHandler(w http.ResponseWriter, r *http.Request) {
	fileHash := r.URL.Query().Get("fileHash")
	row, _ := db.GetFileMeta(fileHash)
	signedURL := oss.DownloadURL(row.FileAddr.String)
	w.Write([]byte(signedURL))
}
