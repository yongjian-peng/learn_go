package system

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"upload_log/dao"
	"upload_log/global"
	"upload_log/model"
	"upload_log/utils"
	"upload_log/utils/upload"
)

type UploadFileLogService struct{}

func (UploadFileLogService *UploadFileLogService) UploadFileLog() {
	// fmt.Print("uploadLog")
	// 读取数据库中的文件

	searchParams := model.SearchParams{}
	searchParams.PageSize = 100
	searchParams.Page = 0
	searchParams.OriginStatus = model.PRODLOG_ORIGIN_STATUS_INIT
	searchParams.Source = global.Config.UploadLogData.OssObjectPrefix

	if err, lists, total := dao.GetProdLogUploadList(searchParams); err != nil {
		fmt.Print(err)
	} else {

		// lists 是 interface{} 类型 需要使用 json 转成 model 来循环
		resByre, err := json.Marshal(lists)
		if err != nil {
			log.Fatal(err)
		}
		// os.Stdout.Write(resByre)
		var prodLogList []model.ProdLogUploadResultModel
		jsonRes := json.Unmarshal(resByre, &prodLogList)
		if jsonRes != nil {
			fmt.Printf("%v", jsonRes)
			return
		}
		// fmt.Printf("使用 json：%v", prodLogList)

		for _, val := range prodLogList {
			if val.FileName != "" {
				if utils.Exists(val.FileName) == false {
					continue
				}
				// 上传文件 & 更新数据库状态
				err = UploadFile(val.FileName)
				if err != nil {
					log.Println("上传资源失败: ", val.FileName)
				}
				// 待上传文件 依次上传到 oss 存储服务
				fmt.Println(val.FileName)
			}
		}
		fmt.Printf("list %d", total)
	}
}

// 项目参考地址 参考来实现的 https://github.com/flipped-aurora/gin-vue-admin 这里出处
//@author: [piexlmax](https://github.com/piexlmax)
//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: filename string
//@return: err error
func UploadFile(filename string) (err error) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(filename)
	if uploadErr != nil {
		return err
	}

	var uploadLogModel = model.ProdLogUploadResultModel{}
	prodLogUploadResultMap := map[string]interface{}{
		"OssFileName":  filePath,
		"MTime":        time.Now().Unix(),
		"OriginStatus": model.PRODLOG_ORIGIN_STATUS_UPLOAD_FINISH,
	}
	log.Println("上传文件成功:", filePath, key)
	db := global.Gorm.Where("file_name = ?", filename).First(&uploadLogModel)
	if uploadLogModel.OriginStatus == model.PRODLOG_ORIGIN_STATUS_INIT {
		err = db.Updates(prodLogUploadResultMap).Error
		log.Printf("更新文件成功:%+v ", uploadLogModel)
		return err
	}

	return nil
}
