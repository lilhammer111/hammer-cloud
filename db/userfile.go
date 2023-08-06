package db

import (
	"github.com/lilhammer111/hammer-cloud/db/mysql"
	"time"
)

// UserFile is mapped by user file table
type UserFile struct {
	Username    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

func OnUserFileUploadFinished(username, fileHash, fileName string, fileSize int64) bool {
	stmt, err := mysql.DBConn().
		Prepare("insert ignore into tbl_user_file (user_name, file_sha1,file_size,file_name, upload_at) VALUES (?,?,?,?,?);")
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, fileHash, fileSize, fileName, time.Now())
	if err != nil {
		return false
	}
	return true
}

// QueryUserFileMetas get user file info in patches
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mysql.DBConn().
		Prepare("select file_sha1, file_name, file_size, upload_at, last_update from tbl_user_file where user_name = ? limit ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err := rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			return nil, err
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}
