package meta

type FileMeta struct {
	FileSize int64
	FileSha1 string
	FileName string
	Location string
	UploadAt string
}

// fileMetas is a map container of FileMeta struct.
var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta updates or adds the FileMeta struct to the map container of fileMetas.
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// GetFileMeta gets a FileMeta struct by fileSha1.
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}
