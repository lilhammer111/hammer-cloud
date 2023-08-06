package db

import (
	"database/sql"
	"fmt"
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

// GetFileMetaList : 从mysql批量获取文件元信息
func GetFileMetaList(limit int) ([]TableFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file " +
			"where status=1 limit ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	cloumns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cloumns))
	var tfiles []TableFile
	for i := 0; i < len(values) && rows.Next(); i++ {
		tfile := TableFile{}
		err = rows.Scan(&tfile.FileHash, &tfile.FileAddr,
			&tfile.FileName, &tfile.FileSize)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	fmt.Println(len(tfiles))
	return tfiles, nil
}

// UpdateFileLocation : 更新文件的存储地址(如文件被转移了)
func UpdateFileLocation(fileHash string, fileAddr string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"update tbl_file set`file_addr`=? where  `file_sha1`=? limit 1")
	if err != nil {
		fmt.Println("预编译sql失败, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileAddr, fileHash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("更新文件location失败, filehash:%s", fileHash)
		}
		return true
	}
	return false
}
