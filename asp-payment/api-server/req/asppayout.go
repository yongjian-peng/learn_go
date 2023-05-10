package req

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"

	"github.com/spf13/cast"
)

//	{
//
// "version":"version",
// "appid":"",
//
// "sign":"",
// "order_id":"order_id",
// "order_currency":"order_currency",
// "order_amount":321456,
// "order_name":"",
// "notify_url ":"",
// "order_note":"",
//
// "customer_name":"",
// "customer_phone":"",
// "customer_email":""
//
// "ifsc": "testpay", string - ifsc-只有印度银行卡提现需要. 注意：非印度银行卡，该字段不要赋值
// "bank_card":"14434233123123", // string - 提现目标银行卡id/pix账号/clabe卡号，其他提现方式可以忽略参数
// "bank_code":"9871", // string - 银行编码-墨西哥clabe提现需要
// "vpa": "xxxx", // string - vpa号,印度upi提现方式需要提供
// "app_order_id":"test-123131asd", // string - 调用方生成的order_id
// "amount": 10, // integer - 金额

// "pay_type":"bank"   // string - 提现类型(枚举值) bank:提现到银行卡(印度地区支持)
//
//	clabe:墨西哥特有的方式，入前确认
//	PIX: 巴西特有方式
//	upi: 仅印度地区
//
// "address": 10, // 允许收款人地址、字母数字和空格
// "city": 10, // 收款城市
//
//	}
//
// AspPayout 提现下单请求参数
type AspPayout struct {
	OrderID       string `label:"order_id" json:"order_id" validate:"required,max=32"`
	UserID        string `label:"user_id" json:"user_id" validate:"required"`
	OrderCurrency string `label:"order_currency" json:"order_currency" validate:"required,eq=INR"`
	OrderAmount   int    `label:"order_amount" json:"order_amount" validate:"required,numeric,gte=1000"`
	OrderName     string `label:"order_name" json:"order_name" validate:"required,lt=64"`
	NotifyURL     string `label:"notify_url" json:"notify_url" validate:"omitempty,url"`
	CustomerName  string `label:"customer_name" json:"customer_name" validate:"required,lt=64"`
	CustomerPhone string `label:"customer_phone" json:"customer_phone" validate:"required,lt=20"`
	CustomerEmail string `label:"customer_email" json:"customer_email"`
	DeviceInfo    string `label:"device_info" json:"device_info" validate:"required"`
	OrderNote     string `label:"order_note" json:"order_note"`
	Attach        string `label:"attach" json:"attach"`
	Ifsc          string `label:"ifsc" json:"ifsc" validate:"required_if=PayType bank"`
	BankCard      string `label:"bank_card" json:"bank_card" validate:"required_if=PayType bank"`
	BankCode      string `label:"bank_code" json:"bank_code"`
	Vpa           string `label:"vpa" json:"vpa" validate:"required_if=PayType upi"`
	PayType       string `label:"pay_type" json:"pay_type" validate:"oneof=bank upi"`
	Address       string `label:"address" json:"address"`
	City          string `label:"city" json:"city"`
	PayoutCode    string `label:"payout_code" json:"payout_code"` // 状态 00 是正常 01 即其他 是需要审核的状态
	PayoutParams
}

type AspPayoutAudit struct {
	Action      string `label:"action" json:"action" validate:"oneof=pass return" comment:"审核类型"`
	OperationID int    `label:"operation_id" json:"operation_id" validate:"required"`
	PayoutID    int    `label:"payout_id" json:"payout_id" validate:"required"`
	PayoutCode  string `label:"payout_code" json:"payout_code"` // 状态 00 是正常 01
}

// AspBeneficiary 提现下单请求参数
type AspBeneficiary struct {
	// NotifyURL     string `label:"notify_url" json:"notify_url" validate:"omitempty,url"`
	CustomerName  string `label:"customer_name" json:"customer_name" validate:"required,lt=64"`
	CustomerPhone string `label:"customer_phone" json:"customer_phone" validate:"required,lt=20"`
	CustomerEmail string `label:"customer_email" json:"customer_email"`
	Ifsc          string `label:"ifsc" json:"ifsc" validate:"required"`
	BankCard      string `label:"bank_card" json:"bank_card" validate:"required"`
	BankCode      string `label:"bank_code" json:"bank_code"`
	Address       string `label:"address" json:"address"`
	City          string `label:"city" json:"city"`
}

// 提现请求的中的 所有用到的请求参数
type PayoutParams struct {
	AppId         string `json:"appid" comment:"商户ID"`
	OrderID       string `json:"order_id"`
	UserID        string `json:"user_id"`
	OrderCurrency string `json:"order_currency"`
	OrderAmount   int    `json:"order_amount"`
	Timestamp     int    `json:"Timestamp"`
	OrderName     string `json:"order_name"`
	NotifyURL     string `json:"notify_url"`
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	CustomerEmail string `json:"customer_email"`
	DeviceInfo    string `json:"device_info"`
	OrderNote     string `json:"order_note"`
	Ifsc          string `json:"ifsc"`
	BankCard      string `json:"bank_card"`
	BankCode      string `json:"bank_code"`
	Vpa           string `json:"vpa"`
	PayType       string `json:"pay_type"`
	Address       string `json:"address"`
	City          string `json:"city"`
	BenefiaryId   string `json:"benefiary_id"` // 受益人id
	Sign          string `json:"sign"`

	ClientIp  string `json:"client_ip"`
	TradeType string `json:"trade_type"`
	Provider  string `json:"provider"`
	Adapter   string `json:"adapter"`
	Sn        string `json:"sn"`
}

// Generate 结构体数据转化 从 SysConfigControl 至 system.SysConfig 对应的模型
func (s *AspPayout) Generate(model *model.AspPayout, aspMerchantProject *model.AspMerchantProject, aspMerchantProjectConfig *model.AspMerchantProjectConfig, merchantProjectCurrencyInfo *model.AspMerchantProjectCurrency, deptTradeTypeInfo *DeptTradeTypeInfo) {
	model.Sn = s.PayoutParams.Sn
	model.DepartId = cast.ToInt(deptTradeTypeInfo.DepartID)
	model.ChannelId = cast.ToUint(deptTradeTypeInfo.ChannelID)
	model.CurrencyId = cast.ToInt(merchantProjectCurrencyInfo.CurrencyId)
	model.MchId = cast.ToUint(aspMerchantProject.MchId)
	model.MchProjectId = aspMerchantProject.Id
	model.UserId = cast.ToUint(s.UserID)
	model.Body = s.OrderName
	model.TransactionId = ""
	model.AdapterTransId = ""
	model.OutTradeNo = s.PayoutParams.OrderID
	model.FeeType = constant.PAYOUT_FEE_TYPE_INR
	model.TotalFee = cast.ToUint(s.OrderAmount)
	model.CashFeeType = ""
	model.CashFee = 0
	model.Rate = 0
	model.FeeRate = cast.ToFloat64(aspMerchantProjectConfig.OutFeeRate)                     // 代付费率
	model.ChargeFee = goutils.ChargeFee(s.OrderAmount, aspMerchantProjectConfig.OutFeeRate) // 传输金额 * 手续费
	model.FixedAmount = aspMerchantProjectConfig.FixedOutAmount                             // 固定手续费(分)
	model.FixedCurrency = aspMerchantProjectConfig.FixedCurrency
	model.TotalChargeFee = model.ChargeFee + model.FixedAmount // 总计手续费（手续费+固定手续费
	model.SpbillCreateIp = s.PayoutParams.ClientIp
	model.ClientIp = s.PayoutParams.ClientIp
	model.NotifyUrl = s.NotifyURL
	model.TradeType = s.PayoutParams.TradeType
	model.DeviceInfo = s.DeviceInfo
	model.CustomerName = s.CustomerName
	model.CustomerEmail = s.CustomerEmail
	model.CustomerPhone = s.CustomerPhone
	model.Vpa = s.Vpa
	model.Ifsc = s.Ifsc
	model.BankCard = s.BankCard
	model.BankCode = s.BankCode
	model.Address = s.Address
	model.City = s.City
	model.PayType = s.PayType
	model.OrderToken = ""
	model.Attach = s.Attach
	model.TimeStart = goutils.GetNowTimesTamp()
	model.TimeExpire = "0"
	model.TimeEnd = "0"
	model.TradeState = constant.PAYOUT_TRADE_STATE_APPLY // 待审核
	model.BankType = ""
	model.IsCallback = 0
	model.BeneficiaryId = s.PayoutParams.BenefiaryId
	model.CreateTime = cast.ToUint64(goutils.GetDateTimeUnix())
	model.FinishTime = 0

	model.Provider = s.PayoutParams.Provider
	model.Adapter = s.PayoutParams.Adapter
	model.Note = s.OrderNote
	model.SupplierReturnCode = ""
	model.SupplierReturnMsg = ""

}

// AspPayoutQuery 提现查询 body 信息
type AspPayoutQuery struct {
	Sn             string `query:"sn" validate:"required`
	OutTradeNo     string `query:"out_trade_no, omitempty"`
	IsCallUpstream string `query:"is_call_upstream, omitempty"`
}

// PayoutAmountList 聚合 sum 查询 amount 的值
type PayoutAmountList struct {
	Amount int `json:"amount"`
}

// PayoutCountList 聚合 sum 查询 条数 的值
type PayoutCountList struct {
	Count int `json:"count"`
}

// AspPayoutUpiValidate 验证upi账号结构体
type AspPayoutUpiValidate struct {
	Vpa string `label:"vpa" json:"vpa"`
}
