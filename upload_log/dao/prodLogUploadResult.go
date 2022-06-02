package dao

import (
	"upload_log/global"
	"upload_log/model"
)

func GetProdLogUploadResult(file_name string) (model.ProdLogUploadResultModel, error) {
	prodLog := model.ProdLogUploadResultModel{}

	result := global.Gorm.Where("file_name = ?", file_name).First(&prodLog)

	if result.Error != nil {
		return prodLog, result.Error
	}
	return prodLog, nil
}
