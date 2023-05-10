package model

// 渠道配置
type AspChannelConfig struct {
	Id            int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	ChannelType   int    `gorm:"column:channel_type;type:tinyint(1);default:1;comment:渠道类型（1收款，2提现，3都可以）;NOT NULL" json:"channel_type"`
	Provider      string `gorm:"column:provider;type:varchar(32);comment:支付的渠道供应商 (小写字符) 筛选项;NOT NULL" json:"provider"`
	Name          string `gorm:"column:name;type:varchar(45);comment:渠道名称 筛选项;NOT NULL" json:"name"`
	Appid         string `gorm:"column:appid;type:varchar(45);comment:保留;NOT NULL" json:"appid"`
	Secret        string `gorm:"column:secret;type:varchar(256);comment:保留;NOT NULL" json:"secret"`
	Key           string `gorm:"column:key;type:varchar(256);comment:保留;NOT NULL" json:"key"`
	Cert          string `gorm:"column:cert;type:varchar(256);comment:保留;NOT NULL" json:"cert"`
	CertKey       string `gorm:"column:cert_key;type:varchar(256);comment:保留;NOT NULL" json:"cert_key"`
	Currency      string `gorm:"column:currency;type:varchar(30);comment:币种;NOT NULL" json:"currency"`
	CurrencyId    int    `gorm:"column:currency_id;type:int(11);default:0;comment:币种id;NOT NULL" json:"currency_id"`
	ChannelConfig string `gorm:"column:channel_config;type:text;comment:渠道配置JSON格式 保留;NOT NULL" json:"channel_config"`
	SerialNo      string `gorm:"column:serial_no;type:varchar(255);comment:机构证书序列号 保留;NOT NULL" json:"serial_no"`
	Status        int    `gorm:"column:status;type:tinyint(1);default:1;comment:渠道状态(0关闭1开启) ;NOT NULL" json:"status"`
	Sort          int    `gorm:"column:sort;type:int(11);default:0;comment:权重(值越大优先级越高，值相同ID大的优先) ;NOT NULL" json:"sort"`
	H5Type        string `gorm:"column:h5_type;type:varchar(32);comment:h5类型（H5 WAPPAY）;NOT NULL" json:"h5_type"`
	CreateTime    int64  `gorm:"column:create_time;type:bigint(20);default:0;NOT NULL" json:"create_time"`
	UpdateTime    int64  `gorm:"column:update_time;type:bigint(20);default:0;NOT NULL" json:"update_time"`
}

func (m *AspChannelConfig) TableName() string {
	return "asp_channel_config"
}
