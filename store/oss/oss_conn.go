package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	cfg "github.com/lilhammer111/hammer-cloud/config"
	"log"
)

var ossCli *oss.Client

// Client creates a ossCli object
func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}

	ossCli, err := oss.New(cfg.OSSEndpoint, cfg.OSSAccessKeyID, cfg.OSSAccessKeySecret)
	if err != nil {
		log.Println(err)
		return nil
	}
	return ossCli
}

// Bucket gets a bucket
func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(cfg.OSSBucket)
		if err != nil {
			log.Println(err)
			return nil
		}
		return bucket
	}
	return nil
}

func DownloadURL(objName string) string {
	signedURL, err := Bucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		log.Println(err)
		return ""
	}
	return signedURL
}
