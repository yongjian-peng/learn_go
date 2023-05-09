package model

import (
	"time"
)

// AppConfig app配置
type AppConfig struct {
	Id                 int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:id" json:"id"`
	AppId              int       `gorm:"column:app_id;type:int(11);default:0;comment:应用id;NOT NULL" json:"app_id"`
	UserDayPayoutLimit int       `gorm:"column:user_day_payout_limit;type:int(11);default:0;comment:个人每日提现最高值;NOT NULL" json:"user_day_payout_limit"`
	PayLimit           int64     `gorm:"column:pay_limit;type:bigint(20);default:0;comment:应用提现最高值;NOT NULL" json:"pay_limit"`
	CreateTime         time.Time `gorm:"column:create_time;type:datetime;comment:创建时间" json:"create_time"`
	UpdateTime         time.Time `gorm:"column:update_time;type:datetime;comment:更新时间" json:"update_time"`
}

func (m *AppConfig) TableName() string {
	return "app_config"
}
