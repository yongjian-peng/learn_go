package system

import "fmt"

type UploadFileLogService struct{}

func (UploadFileLogService *UploadFileLogService) UploadFileLog() {
	fmt.Print("uploadLog")
}
