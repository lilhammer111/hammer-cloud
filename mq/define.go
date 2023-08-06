package mq

import "github.com/lilhammer111/hammer-cloud/common"

// TransferData means message structure
type TransferData struct {
	FileHash      string
	CurLocation   string
	DestLocation  string
	DestStoreType common.StoreType
}
