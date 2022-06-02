package model

const (
	PRODLOG_ORIGIN_STATUS_INIT          = 0  // 原始文件状态: 未处理
	PRODLOG_ORIGIN_STATUS_UPLOAD_FINISH = 1  // 原始文件状态: 已上传到OSS
	PRODLOG_ORIGIN_STATUS_DELETE_FINISH = 2  // 原始文件状态: 已经删除原始文件
	PRODLOG_ORIGIN_STATUS_IS_EMPTY      = 99 // 原始文件状态: 原始文件为空文件，字节数为零
)

type ProdLogUploadResultModel struct {
	Bytes        int64  `json:"bytes"`
	CreateTime   string `json:"create_time"`
	FileName     string `json:"file_name"`
	ID           int64  `json:"id"`
	MTime        int64  `json:"m_time"`
	OriginStatus int64  `json:"origin_status"`
	OssFileName  string `json:"oss_file_name"`
	Sha1         string `json:"sha1"`
	Source       string `json:"source"`
}

// 自定义表名
func (ProdLogUploadResultModel) TableName() string {
	return "prod_log_upload_result"
}
