package meta

import (
	"github.com/lilhammer111/hammer-cloud/db"
)

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
//func UpdateFileMeta(fmeta FileMeta) {
//	fileMetas[fmeta.FileSha1] = fmeta
//}

// UpdateFileMetaDB updates or adds the FileMeta fields in DB.
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.Location, fmeta.FileSize)
}

// GetFileMeta gets a FileMeta struct by fileSha1.
//func GetFileMeta(fileSha1 string) FileMeta {
//	return fileMetas[fileSha1]
//}

// GetFileMetaDB get file meta info from db.
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, nil
	}
	var fmeta FileMeta
	if tfile.FileName.Valid && tfile.FileSize.Valid && tfile.FileAddr.Valid {
		fmeta = FileMeta{
			FileSize: tfile.FileSize.Int64,
			FileSha1: tfile.FileHash,
			FileName: tfile.FileName.String,
			Location: tfile.FileAddr.String,
		}
	}
	return fmeta, nil
}

//func GetLastFileMetas(count int) []FileMeta {
//	fMetaArray := make([]FileMeta, len(fileMetas))
//	for _, v := range fileMetas {
//		fMetaArray = append(fMetaArray, v)
//	}
//	sort.Sort(ByUploadTime(fMetaArray))
//	return fMetaArray[0:count]
//
//}

// RemoveFileMeta removes a file meta info from fileMetas.
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
