package system

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

		b, err := json.Marshal(lists)
		if err != nil {
			log.Fatal(err)
		}

		for _, val := range b {
			fmt.Println(val)
		}

		os.Stdout.Write(b)

		// var out bytes.Buffer
		// json.Indent(&out, b, "=", "\t")
		// out.WriteTo(os.Stdout)

		// fmt.Println(json.Marshal(lists))
		fmt.Printf("list %d", total)

	}
	// 待上传文件 依次上传到 oss 存储服务

	// 上传成功后 依次更新数据库中的状态

}
