package system

type UploadFileLogService struct{}

func (UploadFileLogService *UploadFileLogService) UploadFileLog() {
	// fmt.Print("uploadLog")
	// 读取数据库中的文件

	// 待上传文件 依次上传到 oss 存储服务

	// 上传成功后 依次更新数据库中的状态
}
