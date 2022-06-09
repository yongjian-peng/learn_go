package system

import (
	"encoding/json"
	"fmt"
	"log"
	"upload_log/dao"
	"upload_log/global"
	"upload_log/model"
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
				fmt.Println(val.FileName)
			}
		}

		fmt.Printf("list %d", total)

	}
	// 待上传文件 依次上传到 oss 存储服务

	// 上传成功后 依次更新数据库中的状态

}
