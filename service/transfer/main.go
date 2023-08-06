package main

import (
	"bufio"
	"encoding/json"
	"github.com/lilhammer111/hammer-cloud/config"
	"github.com/lilhammer111/hammer-cloud/db"
	"github.com/lilhammer111/hammer-cloud/mq"
	"github.com/lilhammer111/hammer-cloud/store/oss"
	"log"
	"os"
)

// ProcessTransfer is the real func to do the work of uploading files
func ProcessTransfer(msg []byte) bool {
	// 1. parse msg
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err)
		return false
	}
	// 2. according to temporary file path to create FD
	filed, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err)
		return false
	}

	// 3. upload file by FD
	err = oss.Bucket().PutObject(pubData.DestLocation, bufio.NewReader(filed))
	if err != nil {
		log.Println(err)
		return false
	}

	// 4. update file's path to tbl_file
	suc := db.UpdateFileLocation(pubData.FileHash, pubData.DestLocation)
	if !suc {
		return false
	}
	return true
}

func main() {

	log.Println("start to listen transfer task queue")
	mq.StartConsume(config.TransOSSQueueName, "transfer_oss", ProcessTransfer)
}
