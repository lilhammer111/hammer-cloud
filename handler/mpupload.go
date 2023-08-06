package handler

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/lilhammer111/hammer-cloud/cache/rd"
	"github.com/lilhammer111/hammer-cloud/db"
	"github.com/lilhammer111/hammer-cloud/util"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	chunkSize = 5 * 1024 * 1024
)

type MultipartUploadInfo struct {
	ChunkSize  int
	ChunkCount int
	Filesize   int
	FileHash   string
	UploadID   string
}

func InitMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse params
	//err := rd.ParseForm()
	//if err != nil {
	//	return
	//}
	//username := rd.Form.Get("username")
	//fileHash := rd.Form.Get("fileHash")
	//filesizeS := rd.Form.Get("filesize")
	params, err := ParseRequestParams(r, "username", "fileHash", "filesize")
	if err != nil {
		http.Error(w, ErrInvalidParam, http.StatusBadRequest)
		return
	}
	filesize, err := strconv.Atoi(params["filesize"])
	if err != nil {
		http.Error(w, ErrStrConvError, http.StatusInternalServerError)
		return
	}

	// 2. get a rd conn
	rConn := rd.Pool().Get()
	defer rConn.Close()

	// 3. init multi part info
	upInfo := MultipartUploadInfo{
		Filesize:   filesize,
		ChunkSize:  chunkSize,
		ChunkCount: int(math.Ceil(float64(filesize) / chunkSize)),
		FileHash:   params["fileHash"],
		UploadID:   params["username"] + fmt.Sprintf("%x", time.Now().UnixNano()),
	}

	// 4. save info into rd
	_, err = rConn.Do("hset", "MP_"+upInfo.UploadID,
		"chunkCount", upInfo.ChunkCount,
		"fileHash", upInfo.FileHash,
		"filesize", upInfo.Filesize)
	if err != nil {
		http.Error(w, ErrMysqlError, http.StatusInternalServerError)
		return
	}
	// 5. return info to client
	_, err = w.Write(util.NewRespBody(0, "OK", upInfo).JSONBytes())
	if err != nil {
		http.Error(w, ErrRespError, http.StatusInternalServerError)
		return
	}
}

func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse params
	//params, err := ParseRequestParams(r, "uploadID", "index")
	//if err != nil {
	//	http.Error(w, ErrInvalidParam, http.StatusBadRequest)
	//	return
	//}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, ErrInvalidParam, http.StatusBadRequest)
		return
	}
	uploadID := r.FormValue("uploadID")
	index := r.FormValue("index")
	file, _, err := r.FormFile("file") // "file" 是文件字段的名称
	if err != nil {
		http.Error(w, ErrInvalidParam, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. get a rd conn
	rConn := rd.Pool().Get()
	defer rConn.Close()

	// 3. get file description to save chunks
	fpath := TempPartRootDir + uploadID + "/" + index
	err = os.MkdirAll(path.Dir(fpath), 0744)
	if err != nil {
		http.Error(w, ErrMkdirError, http.StatusInternalServerError)
		return
	}
	fd, err := os.Create(fpath)
	if err != nil {
		http.Error(w, ErrCreateError, http.StatusInternalServerError)
		return
	}
	defer fd.Close()

	//buf := make([]byte, 1024*1024)
	//for {
	//	n, err := r.Body.Read(buf)
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		http.Error(w, ErrReadError, http.StatusInternalServerError)
	//		return
	//	}
	//	_, err = fd.Write(buf[:n])
	//
	//	if err != nil {
	//		http.Error(w, ErrWriteError, http.StatusInternalServerError)
	//		return
	//	}
	//}
	_, err = io.Copy(fd, file)
	if err != nil {
		http.Error(w, ErrReadError, http.StatusInternalServerError)
		return
	}

	// 4. update redis cache state
	_, err = rConn.Do("hset", "MP_"+uploadID, "chkidx_"+index, 1)
	if err != nil {
		http.Error(w, ErrRedisError, http.StatusInternalServerError)
		return
	}

	// 5. return result
	_, err = w.Write(util.NewRespBody(0, "OK", nil).JSONBytes())
	if err != nil {
		http.Error(w, ErrWriteError, http.StatusInternalServerError)
		return
	}
}

func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse params
	params, err := ParseRequestParams(r, "uploadID", "username", "fileHash", "filesize", "filename")
	if err != nil {
		http.Error(w, ErrInvalidParam, http.StatusBadRequest)
		return
	}

	// 2. get a rd conn
	rConn := rd.Pool().Get()
	defer rConn.Close()

	// 3. judge if all chunks have been uploaded by uploadID
	resArr, err := redis.Values(rConn.Do("hgetall", "MP_"+params["uploadID"]))
	if err != nil {
		http.Error(w, ErrRedisError, http.StatusInternalServerError)
		return
	}
	totalCount := 0
	var chunkCount int
	for i := 0; i < len(resArr); i += 2 {
		k := string(resArr[i].([]byte))
		v := string(resArr[i+1].([]byte))
		if k == "chunkCount" {
			chunkCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			totalCount++
		}
	}

	if totalCount != chunkCount {
		_, err := w.Write(util.NewRespBody(-2, "invalid request", nil).JSONBytes())
		if err != nil {
			http.Error(w, ErrRespError, http.StatusInternalServerError)
		}
		return
	}

	// 4. TODO: MERGE CHUNK

	// 5. update tbl_user_file and tbl_file
	filesize, err := strconv.ParseInt(params["filesize"], 10, 64)
	if err != nil {
		http.Error(w, ErrStrConvError, http.StatusInternalServerError)
	}

	db.OnFileUploadFinished(params["fileHash"], params["filename"], "", filesize)
	db.OnUserFileUploadFinished(params["username"], params["fileHash"], params["filename"], filesize)

	// 6. response normally
	_, err = w.Write(util.NewRespBody(0, "OK", nil).JSONBytes())
	if err != nil {
		http.Error(w, ErrWriteError, http.StatusInternalServerError)
		return
	}
}

func CancelUploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 删除已存在的分块文件
	// 删除redis缓存状态
	// 更新mysql 用户上传记录表status
}

func MultipartUploadStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 检查分块上传状态是否有效
	// 获取分块初始化信息
	// 获取已上传的分块信息
}
