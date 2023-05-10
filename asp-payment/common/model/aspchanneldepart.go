package model

// AspChannelDepartConfig 渠道商户关系表channel_config 对应 admin_departs
type AspChannelDepartConfig struct {
	Id           int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	DepartId     uint   `gorm:"column:depart_id;type:int(10) unsigned;default:0;comment:平台）商户或代理商id，对应 admin_departs 主键 ID;NOT NULL" json:"depart_id"`
	ChannelId    int    `gorm:"column:channel_id;type:int(11);default:0;comment:渠道id(y2p)　对应　channel_config 主键 ID;NOT NULL" json:"channel_id"`
	ChannelMchId string `gorm:"column:channel_mch_id;type:varchar(45);comment:渠道商户号：渠道方为商户平台申请的商户号。;NOT NULL" json:"channel_mch_id"`
	Status       int    `gorm:"column:status;type:int(11);default:0;comment:渠道商户审核状态: 0 未审核 1:审核通过 2:审核未通过;NOT NULL" json:"status"`
	Remark       string `gorm:"column:remark;type:varchar(1024);comment:其他：如退回原因;NOT NULL" json:"remark"`
	Config       string `gorm:"column:config;type:text;comment:商户自定义配置;NOT NULL" json:"config"`
	PassTime     uint64 `gorm:"column:pass_time;type:bigint(20) unsigned;default:0;comment:代理开通(审核通过)时间戳;NOT NULL" json:"pass_time"`
	PayoutCount  uint   `gorm:"column:payout_count;type:int(10) unsigned;default:20;comment:商户API接口，一天可允许的次数，0代表无限;NOT NULL" json:"payout_count"`
}

func (m *AspChannelDepartConfig) TableName() string {
	return "asp_channel_depart_config"
}

type AspChannelDepartConfigInfo struct {
	Appid           string `json:"appid"`
	Signature       string `json:"signature"`
	PartnerId       string `json:"partnerId"`
	ApplicationId   string `json:"applicationId"`
	PayoutSignature string `json:"payoutSignature"`
}
