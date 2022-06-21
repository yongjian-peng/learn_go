package system

import (
	"encoding/json"
	"log"
	"time"
	"upload_log/dao"
	"upload_log/global"
	"upload_log/model"
	"upload_log/utils"
	"upload_log/utils/upload"
)

type ClearFileService struct{}

func (ClearFileService *ClearFileService) ClearFile() {
	// 读取数据库文件 状态为上传成功的
	searchParams := model.SearchParams{}
	searchParams.PageSize = 100
	searchParams.Page = 0
	searchParams.OriginStatus = model.PRODLOG_ORIGIN_STATUS_UPLOAD_FINISH
	searchParams.Source = global.Config.UploadLogData.OssObjectPrefix

	if err, lists, _ := dao.GetProdLogUploadList(searchParams); err != nil {
		log.Println("上传日志，第三部：", err)
	} else {
		// lists 是 interface{} 类型 需要使用 json 转成model 来循环
		resByre, err := json.Marshal(lists)
		if err != nil {
			log.Fatal(err)
		}
		var prodLogList []model.ProdLogUploadResultModel
		jsonRes := json.Unmarshal(resByre, &prodLogList)
		if jsonRes != nil {
			log.Printf("%v", jsonRes)
			return
		}

		for _, val := range prodLogList {
			if val.FileName != "" {
				if utils.Exists(val.FileName) == false {
					continue
				}
				// 删除文件 并且 更新状态
				err = ClearFileAndUpdateState(val.FileName)
				if err != nil {
					log.Printf("%v", err)
				}
			}
		}
	}
	// 删除对应的文件

	// 依次更新数据库中文件 状态为已完成状态
}

func ClearFileAndUpdateState(filename string) (err error) {
	oss := upload.NewOss()
	err = oss.DeleteFile(filename)
	if err != nil {
		log.Printf("%v", err)
	}

	// 更新数据库中的状态
	var uploadLogModel = model.ProdLogUploadResultModel{}
	prodLogUploadResultMap := map[string]interface{}{
		"MTime":        time.Now().Unix(),
		"OriginStatus": model.PRODLOG_ORIGIN_STATUS_DELETE_FINISH,
	}

	db := global.Gorm.Where("file_name = ?", filename).First(&uploadLogModel)
	if uploadLogModel.OriginStatus == model.PRODLOG_ORIGIN_STATUS_UPLOAD_FINISH {
		err = db.Updates(prodLogUploadResultMap).Error
		log.Printf("删除文件，更新文件成功：%v", uploadLogModel)
		return nil
	}

	return nil
}
