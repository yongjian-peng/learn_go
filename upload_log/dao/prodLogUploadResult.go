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

func GetProdLogUploadList(searchParms model.SearchParams) (err error, list interface{}, total int64) {
	limit := searchParms.PageSize

	offset := searchParms.PageSize * (searchParms.Page - 1)

	prodLog := model.ProdLogUploadResultModel{}

	db := global.Gorm.Model(&prodLog)
	var prodLogList []model.ProdLogUploadResultModel

	if searchParms.OriginStatus == 0 {
		db = db.Where("origin_status = ? ", searchParms.OriginStatus)
	}
	if searchParms.Source != "" {
		db = db.Where("source = ? ", searchParms.Source)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&prodLogList).Error
	return err, prodLogList, total

}
