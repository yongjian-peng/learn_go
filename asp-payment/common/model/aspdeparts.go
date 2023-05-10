package model

// AspDeparts 内部代理商户/关系表
type AspDeparts struct {
	Id               int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	ParentId         int    `gorm:"column:parent_id;type:int(11);default:0;comment:代理或者商户的上级　id;NOT NULL" json:"parent_id"`
	DepartType       int    `gorm:"column:depart_type;type:tinyint(1);default:1;comment:类型 1:代理 2：商户;NOT NULL" json:"depart_type"`
	Title            string `gorm:"column:title;type:varchar(45);comment:内部商户名称;NOT NULL" json:"title"`
	CurrencyId       uint   `gorm:"column:currency_id;type:int(10) unsigned;default:0;comment:币种id ;NOT NULL" json:"currency_id"`
	Sort             int    `gorm:"column:sort;type:int(11);default:0;comment:排序;NOT NULL" json:"sort"`
	ContactsName     string `gorm:"column:contacts_name;type:varchar(45);comment:联系人姓名;NOT NULL" json:"contacts_name"`
	ContactsPhone    string `gorm:"column:contacts_phone;type:varchar(45);comment:联系人电话;NOT NULL" json:"contacts_phone"`
	ContactsEmail    string `gorm:"column:contacts_email;type:varchar(256);comment:联系人邮箱;NOT NULL" json:"contacts_email"`
	BankName         string `gorm:"column:bank_name;type:varchar(256);comment:收款银行名称;NOT NULL" json:"bank_name"`
	BankOfDeposit    string `gorm:"column:bank_of_deposit;type:varchar(256);comment:开户行名称;NOT NULL" json:"bank_of_deposit"`
	BankAccount      string `gorm:"column:bank_account;type:varchar(64);comment:银行账户账号;NOT NULL" json:"bank_account"`
	BankAccountName  string `gorm:"column:bank_account_name;type:varchar(256);comment:银行账户名称;NOT NULL" json:"bank_account_name"`
	IndustryId       uint   `gorm:"column:industry_id;type:int(10) unsigned;default:0;comment:行业类别 保留;NOT NULL" json:"industry_id"`
	CompanyName      string `gorm:"column:company_name;type:varchar(128);comment:公司全称 保留;NOT NULL" json:"company_name"`
	CompanyShort     string `gorm:"column:company_short;type:varchar(512);comment:公司简介 保留;NOT NULL" json:"company_short"`
	CompanyEmail     string `gorm:"column:company_email;type:varchar(255);comment:公司邮箱 保留;NOT NULL" json:"company_email"`
	Remark           string `gorm:"column:remark;type:varchar(128);comment:备注;NOT NULL" json:"remark"`
	Status           int    `gorm:"column:status;type:int(11);default:0;comment:状态 ：0 默认（填写资料未提交审核）； 1：已审核（渠道已申请） 2已审核（申请渠道中）3:已关闭 4 退回 5未审核/提交申请中 6渠道费率未设置;NOT NULL" json:"status"`
	CreateTime       uint64 `gorm:"column:create_time;type:bigint(20) unsigned;default:0;NOT NULL" json:"create_time"`
	UpdateTime       uint64 `gorm:"column:update_time;type:bigint(20) unsigned;default:0;NOT NULL" json:"update_time"`
	DefaultChannelId int    `gorm:"column:default_channel_id;type:int(11);default:0;comment:默认启用的渠道;NOT NULL" json:"default_channel_id"`
	SettlementType   int    `gorm:"column:settlement_type;type:tinyint(1);default:1;comment:结算类型(1.自动结算 2.手动结算);NOT NULL" json:"settlement_type"`
	PayoutCount      uint   `gorm:"column:payout_count;type:int(10) unsigned;default:20;comment:商户API接口，一天可允许的次数，0代表无限;NOT NULL" json:"payout_count"`
}

func (AspDeparts) TableName() string {
	return "asp_departs"
}
