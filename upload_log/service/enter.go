package service

import (
	"context"
	"upload_log/service/system"
)

func UploadLogOss(ctx context.Context) {
	set := system.ReadFileService{}
	set.ReadFile("./")
}
