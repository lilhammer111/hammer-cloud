package handler

const (
	ErrInvalidParam = "invalid param"
	ErrMysqlError   = "database mysql error"
	ErrRedisError   = "database rd error"
	ErrRespError    = "response write error"
	ErrStrConvError = "string convert error"
	ErrMkdirError   = "mkdir error"
	ErrCreateError  = "create file error"
	ErrReadError    = "read data failed"
	ErrWriteError   = "write data failed"
	ErrSeekError    = "failed to seek"
	ErrOSSPutError  = "oss put fail"
	ErrMarshalError = "marshal data error"

	TempPartRootDir = "/data/hammer_cloud_part"
	OSSRootDir      = "test/"
)
