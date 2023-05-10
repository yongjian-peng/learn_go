package req

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"

	"github.com/spf13/cast"
)

//	{
//	    "version":"version",
//	    "payment":"",
//
// "mid":"",
// "sign":"",
// "order_id":"order_id",
// "order_currency":"order_currency",
// "order_amount":321456,
// "return_url":"",
// "notify_url ":"",
// "order_note":"",
//
//	    "customer_name":"",
//	    "customer_phone":"",
//	    "customer_email":""
//	}
//
// AspPayment 支付下单请求参数
type AspPayment struct {
	PaymentMethod string `label:"payment_method" json:"payment_method" validate:"oneof=sunny.h5 sunny.wappay" comment:"商户ID"`
	OrderID       string `label:"order_id" json:"order_id" validate:"required,max=32"`
	UserID        string `label:"user_id" json:"user_id" validate:"required"`
	OrderCurrency string `label:"order_currency" json:"order_currency" validate:"required,eq=INR"`
	OrderAmount   int    `label:"order_amount" json:"order_amount" validate:"required,numeric,gte=10000"`
	OrderName     string `label:"order_name" json:"order_name" validate:"required,lt=64"`
	ReturnURL     string `label:"return_url" json:"return_url" validate:"omitempty,url"`
	NotifyURL     string `label:"notify_url" json:"notify_url" validate:"omitempty,url"`
	CustomerName  string `label:"customer_name" json:"customer_name" validate:"required,lt=64"`
	CustomerPhone string `label:"customer_phone" json:"customer_phone" validate:"required,lt=20"`
	CustomerEmail string `label:"customer_email" json:"customer_email"`
	DeviceInfo    string `label:"device_info" json:"device_info" validate:"required"`
	OrderNote     string `label:"order_note" json:"order_note"`
	Attach        string `label:"attach" json:"attach"`
	PayParams
}

type AspPaymentHeader struct {
	Version   string `label:"Version" json:"Version" validate:"oneof=1.0"`
	Timestamp int    `label:"Timestamp" json:"Timestamp" validate:"required,numeric"`
	AppId     string `label:"AppId" json:"AppId" validate:"required" comment:"商户ID"`
	Signature string `label:"Signature" json:"Signature" validate:"required" comment:"签名"`
}

// PayParams 支付请求的中的 所有用到的请求参数
type PayParams struct {
	PaymentMethod string `json:"payment_method" comment:"商户ID"`
	AppId         string `json:"appid" comment:"商户ID"`
	OrderID       string `json:"order_id"`
	UserID        string `json:"user_id"`
	OrderCurrency string `json:"order_currency"`
	OrderAmount   int    `json:"order_amount"`
	Timestamp     int    `json:"Timestamp"`
	OrderName     string `json:"order_name"`
	ReturnURL     string `json:"return_url" validate:"omitemptyurl"`
	NotifyURL     string `json:"notify_url" validate:"omitemptyurl"`
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	CustomerEmail string `json:"customer_email"`
	DeviceInfo    string `json:"device_info"`
	OrderNote     string `json:"order_note"`
	Attach        string `json:"attach"`
	Sign          string `json:"sign"`

	ClientIp  string `json:"client_ip"`
	TradeType string `json:"trade_type"`
	Provider  string `json:"provider"`
	Adapter   string `json:"adapter"`
	Sn        string `json:"sn"`
}

// Generate 结构体数据转化 从 SysConfigControl 至 system.SysConfig 对应的模型
func (s *AspPayment) Generate(model *model.AspOrder, aspMerchantProject *model.AspMerchantProject, aspMerchantProjectConfig *model.AspMerchantProjectConfig, merchantProjectCurrencyInfo *model.AspMerchantProjectCurrency, deptTradeTypeInfo *DeptTradeTypeInfo) {
	model.Sn = s.PayParams.Sn
	model.DepartId = deptTradeTypeInfo.DepartID
	model.ChannelId = deptTradeTypeInfo.ChannelID
	model.CurrencyId = cast.ToInt(merchantProjectCurrencyInfo.CurrencyId)
	model.ChannelMchId = ""
	model.MchId = aspMerchantProject.MchId
	model.MchProjectId = aspMerchantProject.Id
	model.UserId = 0
	model.Body = s.OrderName
	model.TransactionId = ""
	model.AdapterTransId = ""
	model.OutTradeNo = s.PayParams.OrderID
	model.FeeType = constant.ORDER_FEE_TYPE_INR
	model.TotalFee = cast.ToUint(s.OrderAmount)
	model.Discount = 0
	model.CashFeeType = constant.ORDER_FEE_TYPE_INR
	model.CashFee = 0
	model.Rate = 0
	model.FeeRate = cast.ToFloat64(aspMerchantProjectConfig.OutFeeRate)                    // 代付费率
	model.ChargeFee = goutils.ChargeFee(s.OrderAmount, aspMerchantProjectConfig.InFeeRate) // 传输金额 * 手续费
	model.FixedAmount = aspMerchantProjectConfig.FixedInAmount                             // 固定手续费(分)
	model.FixedCurrency = aspMerchantProjectConfig.FixedCurrency
	model.TotalChargeFee = model.ChargeFee + model.FixedAmount // 总计手续费（手续费+固定手续费
	model.SpbillCreateIp = s.PayParams.ClientIp
	model.ClientIp = s.PayParams.ClientIp
	model.NotifyUrl = s.NotifyURL
	model.TradeType = s.PayParams.TradeType
	model.H5Type = deptTradeTypeInfo.H5Type
	model.AuthCode = ""
	model.DeviceInfo = s.DeviceInfo
	model.CustomerId = s.UserID
	model.CustomerName = s.CustomerName
	model.CustomerEmail = s.CustomerEmail
	model.CustomerPhone = s.CustomerPhone
	model.OrderToken = ""
	model.PaymentLink = ""
	model.PaymentsUrl = ""
	model.SettlementsUrl = ""
	model.Attach = s.PayParams.Attach
	model.TimeStart = cast.ToString(goutils.GetDateTimeUnix())
	model.TimeExpire = "0"
	model.TimeEnd = "0"
	model.GoodsTag = ""
	model.TradeState = constant.ORDER_TRADE_STATE_PENDING
	model.BankType = ""
	model.IsCallback = 0
	model.CreateTime = cast.ToUint64(goutils.GetDateTimeUnix())
	model.FinishTime = 0
	model.SettlementTime = 0

	model.IsBill = 0
	model.Wallet = constant.ORDER_FEE_TYPE_INR
	model.Region = constant.ORDER_REGION
	model.Provider = s.PayParams.Provider
	model.Adapter = s.PayParams.Adapter
	model.Qrcode = ""
	model.Note = s.OrderNote
	model.ReturnUrl = s.ReturnURL
	model.UnionpayAppend = "{}"

}
