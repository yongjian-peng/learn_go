package req

import "asp-payment/common/model"

type DepartIdList struct {
	DepartId int `json:"depart_id"`
}

// fdsf
type AspDepartList struct {
	Id               int    `json:"id"`
	ParentId         int    `json:"parent_id"`
	DepartType       int    `json:"depart_type"`
	Title            string `json:"title"`
	CurrencyId       int    `json:"currency_id"`
	Sort             int    `json:"sort"`
	ContactsName     string `json:"contacts_name"`
	ContactsPhone    string `json:"contacts_phone"`
	ContactsEmail    string `json:"contacts_email"`
	BankName         string `json:"bank_name"`
	BankOfDeposit    string `json:"bank_of_deposit"`
	BankAccount      string `json:"bank_account"`
	BankAccountName  string `json:"bank_account_name"`
	CompanyName      string `json:"company_name"`
	CompanyShort     string `json:"company_short"`
	CompanyEmail     string `json:"company_email"`
	Status           int    `json:"status"`
	CreateTime       uint64 `json:"create_time"`
	UpdateTime       uint64 `json:"update_time"`
	Remark           string `json:"remark"`
	DefaultChannelId int    `json:"default_channel_id"`
	SettlementType   int    `json:"settlement_type"`
}

func (s *AspDepartList) Generate(model *model.AspDeparts) {
	model.Id = s.Id
	model.ParentId = s.ParentId
	model.DepartType = s.DepartType
	model.Title = s.Title
	model.Sort = s.Sort
	model.ContactsName = s.ContactsName
	model.ContactsPhone = s.ContactsPhone
	model.ContactsEmail = s.ContactsEmail
	model.BankName = s.BankName
	model.BankOfDeposit = s.BankOfDeposit
	model.BankAccount = s.BankAccount
	model.BankAccountName = s.BankAccountName
	model.CompanyName = s.CompanyName
	model.CompanyShort = s.CompanyShort
	model.CompanyEmail = s.CompanyEmail
	model.Status = s.Status
	model.CreateTime = s.CreateTime
	model.UpdateTime = s.UpdateTime
	model.Remark = s.Remark
	model.DefaultChannelId = s.DefaultChannelId
	model.SettlementType = s.SettlementType
}

// AspMerchantProjectDepartList 添加收益人 关联商户排序结构体 使用的是 关联表中的 排序值，给到的是渠道商户的信息
type AspMerchantProjectDepartList struct {
	Id               int    `json:"id"`
	ParentId         int    `json:"parent_id"`
	DepartType       int    `json:"depart_type"`
	Title            string `json:"title"`
	CurrencyId       int    `json:"currency_id"`
	Sort             int    `json:"sort"`
	MchProjectSort   int    `json:"mch_project_sort"`
	ContactsName     string `json:"contacts_name"`
	ContactsPhone    string `json:"contacts_phone"`
	ContactsEmail    string `json:"contacts_email"`
	BankName         string `json:"bank_name"`
	BankOfDeposit    string `json:"bank_of_deposit"`
	BankAccount      string `json:"bank_account"`
	BankAccountName  string `json:"bank_account_name"`
	CompanyName      string `json:"company_name"`
	CompanyShort     string `json:"company_short"`
	CompanyEmail     string `json:"company_email"`
	Status           int    `json:"status"`
	CreateTime       uint64 `json:"create_time"`
	UpdateTime       uint64 `json:"update_time"`
	Remark           string `json:"remark"`
	DefaultChannelId int    `json:"default_channel_id"`
	SettlementType   int    `json:"settlement_type"`
}
