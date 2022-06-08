package system

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"upload_log/global"
	"upload_log/logic"
	"upload_log/model"
)

type ReadFileService struct{}

func (ReadFileService *ReadFileService) ReadFile(file string) {
	// 依次读取目录中的文件
	uploadPath := global.Config.UploadLogData.BizSearchDir
	ossObjectPrefix := global.Config.UploadLogData.OssObjectPrefix

	//打开目录
	f, err := os.OpenFile(uploadPath, os.O_RDONLY, os.ModeDir)
	if err != nil {
		fmt.Println("openfile err:", err)
		return
	}
	defer f.Close()
	//读取目录项
	info, err := f.Readdir(-1) //-1读取目录中的所有目录项
	if err != nil {
		fmt.Println("readdir err:", err)
		return
	}
	//变量返回切片
	for _, fileInfo := range info {
		if !fileInfo.IsDir() {
			// fmt.Println("png文件有:", fileInfo.Name())
			if strings.HasSuffix(fileInfo.Name(), ".log") {
				fmt.Println("log文件有:", fileInfo.Name())
				// 组装sql数据
				filename := fileInfo.Name()

				filenameWith := logic.BuildLocalDataName(filename)
				// 写入到数据库

				// 保存数据库
				uploadLogModel := model.ProdLogUploadResultModel{}

				// SN查询
				result := global.Gorm.Where("file_name = ?", filenameWith).Last(&uploadLogModel)

				uploadLogModel.FileName = filenameWith
				uploadLogModel.OssFileName = logic.BuildOssDataName(filename)
				uploadLogModel.MTime = time.Now().Unix()
				uploadLogModel.CreateTime = time.Now().Format("2006-01-02 15:04:05")
				uploadLogModel.OriginStatus = model.PRODLOG_ORIGIN_STATUS_INIT
				uploadLogModel.Bytes = 1
				uploadLogModel.Sha1 = logic.BuildFileHash(filenameWith)
				uploadLogModel.Source = ossObjectPrefix

				// fmt.Println(uploadLogModel.Sha1)

				if result.RowsAffected > 0 {
					// 更新
					res := global.Gorm.Save(&uploadLogModel)
					if res.Error != nil {
						log.Printf("上传日志更新DB失败:[失败原因]:%+v [待插入的值]%+v", res.Error, uploadLogModel)
					}
				} else {
					// 插入
					res := global.Gorm.Create(&uploadLogModel)
					if res.Error != nil {
						log.Printf("上传日志插入DB失败:[失败原因]:%+v [待插入的值]:%+v", res.Error, uploadLogModel)
					}
				}
				// 记录日志
				log.Printf("读取文件成功，写入数据库:%+v", uploadLogModel)
			}
		}
	}

	// 扫描文件组成字符

	// 插入到数据库
}
