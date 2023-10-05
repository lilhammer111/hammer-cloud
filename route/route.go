package route

import (
	"github.com/gin-gonic/gin"
	h "github.com/lilhammer111/hammer-cloud/handler"
	m "github.com/lilhammer111/hammer-cloud/middleware"
)

func Router() *gin.Engine {
	// gin framework including logger, recovery
	router := gin.Default()
	// handle static assets
	router.Static("/static/", "./static/")
	// the apis that don't need to auth first
	router.GET("/user/signup", h.SignUpHandler)
	router.POST("/user/signup", h.DoSignUpHandler)

	router.GET("/user/login", h.LoginHandler)
	router.POST("/user/login", h.DoLoginHandler)

	// middleware
	router.Use(m.TokenAuth())
	// after use middleware
	//http.HandleFunc("/user/info", m.TokenAuth(h.UserInfoHandler))
	//// file api
	//http.HandleFunc("/file/upload", m.TokenAuth(h.UploadHandler))
	//http.HandleFunc("/file/upload/suc", m.TokenAuth(h.UploadSucHandler))
	//http.HandleFunc("/file/meta", m.TokenAuth(h.GetFileMetaHandler))
	//http.HandleFunc("/file/query", m.TokenAuth(h.FileQueryHandler))
	//http.HandleFunc("/file/download", m.TokenAuth(h.DownloadFileHandler))
	//http.HandleFunc("/file/update", m.TokenAuth(h.FileMetaUpdateHandler))
	//http.HandleFunc("/file/delete", m.TokenAuth(h.FileDeleteHandler))
	////fast upload api
	//http.HandleFunc("/file/fastupload", m.TokenAuth(h.TryFastUploadHandler))
	//
	//// multi part upload api
	//http.HandleFunc("/file/mpupload/init", m.TokenAuth(h.InitMultipartUploadHandler))
	//http.HandleFunc("/file/mpupload/uppart", m.TokenAuth(h.UploadPartHandler))
	//http.HandleFunc("/file/mpupload/complete", m.TokenAuth(h.CompleteUploadHandler))
	//
	//// oss api
	//http.HandleFunc("/file/downloadurl", m.TokenAuth(h.DownloadURLHandler))

	return router
}
