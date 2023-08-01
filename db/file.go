package db

import (
	"database/sql"
	"github.com/lilhammer111/hammer-cloud/db/mysql"
	"log"
)

// OnFileUploadFinished insert file meta into db
func OnFileUploadFinished(fileHash, filename, fileAddr string, filesize int64) bool {
	stmtIns, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_file (file_sha1, file_name, file_size, file_addr, status) values (?,?,?,?,1)")
	if err != nil {
		log.Println(err)
		return false
	}

	defer stmtIns.Close()

	execRes, err := stmtIns.Exec(fileHash, filename, filesize, fileAddr)
	if err != nil {
		log.Println(err)
		return false
	}

	if affected, err := execRes.RowsAffected(); err == nil {
		if affected <= 0 {
			log.Printf("File with hash: %s has may been uploaded before", fileHash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// GetFileMeta gets file meta struct by fileSha1
func GetFileMeta(fileHash string) (*TableFile, error) {
	stmtSel, err := mysql.DBConn().
		Prepare("select file_sha1, file_addr,file_name, file_size from tbl_file where file_sha1 = ? and status = 1 limit 1;")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmtSel.Close()

	tfile := TableFile{}
	err = stmtSel.QueryRow(fileHash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &tfile, nil
}
