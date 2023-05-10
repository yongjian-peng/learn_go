package interfaces

import (
	"asp-payment/api-server/req"
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
)

// ScanData 定义扫码支付的返回值
type ScanData struct {
	CashFeeType    string `json:"cash_fee_type"` // 现金支付货币类型
	CashFee        int    `json:"cash_fee"`      // 实际支付金额
	FinishTime     int64  `json:"finish_time"`   // 完成的时间
	TradeState     string `json:"trade_state"`   // 支付的状态
	Rate           int    `json:"rate"`          // 费率
	Qrcode         string `json:"qrcode"`
	SettlementTime uint64 `json:"settlement_time"`
	PaymentLink    string `json:"payment_link"`
	PaymentsURL    string `json:"payments_url"`
	ReturnURL      string `json:"return_url"`
	SettlementsURL string `json:"settlements_url"`
	BankType       string `json:"bank_type"`
	TransactionID  string `json:"transaction_id"`
	Code           string `json:"code"` //状态码 '00' 成功
	Msg            string `json:"msg"`  //附加的信息
}

// ThirdQueryData 定义扫码支付的返回值
type ThirdQueryData struct {
	CashFeeType   string `json:"cash_fee_type"` // 现金支付货币类型
	CashFee       int    `json:"cash_fee"`      // 实际支付金额
	TransactionID string `json:"transaction_id"`
	FinishTime    int64  `json:"finish_time"` // 完成的时间
	TradeState    string `json:"trade_state"` // 支付的状态
	Utr           string `json:"urt"`         // 支付的utr编号
	Code          string `json:"code"`        //状态码 '00' 成功
	Msg           string `json:"msg"`         //附加的信息
}

type ThirdPayoutCreateData struct {
	CashFeeType   string `json:"cash_fee_type"` // 现金支付货币类型
	CashFee       int    `json:"cash_fee"`      // 实际支付金额
	FinishTime    int64  `json:"finish_time"`   // 完成的时间
	TradeState    string `json:"trade_state"`   // 支付的状态
	BankType      string `json:"bank_type"`
	BankUtr       string `json:"bank_utr"`
	TransactionID string `json:"transaction_id"`
	Code          string `json:"code"` //状态码 '00' 成功
	Msg           string `json:"msg"`  //附加的信息
}

type ThirdPayoutQueryData struct {
	CashFeeType   string `json:"cash_fee_type"` // 现金支付货币类型
	CashFee       int    `json:"cash_fee"`      // 实际支付金额
	FinishTime    int64  `json:"finish_time"`   // 完成的时间
	TradeState    string `json:"trade_state"`   // 支付的状态
	BankType      string `json:"bank_type"`
	BankUtr       string `json:"bank_utr"`
	TransactionID string `json:"transaction_id"`
	Code          string `json:"code"` //状态码 '00' 成功
	Msg           string `json:"msg"`  //附加的信息
}

type ThirdMerchantAccountQueryData struct {
	Balance          int64  `json:"balance"`          // 账户余额 可用余额
	UnsettledBalance int64  `json:"unsettledBalance"` // 未结算余额
	Code             string `json:"code"`             //状态码 '00' 成功
	Msg              string `json:"msg"`              //附加的信息
}

type ThirdAddBeneficiary struct {
	BenefiaryId string `json:"benefiary_id"` // 返回的受益人标识 id
	Code        string `json:"code"`         //状态码 '00' 成功
	Msg         string `json:"msg"`          //附加的信息
}

type ThirdUpiValidate struct {
	Code string `json:"code"` // 状态码 '00' 成功
	Msg  string `json:"msg"`  // 返回的消息
}

type PayInterface interface {
	H5(string, *model.AspChannelDepartConfig, *model.AspOrder) (*ScanData, *appError.Error)
	WAPPAY(string, *model.AspChannelDepartConfig, *model.AspOrder) (model.BodyMap, *appError.Error)
	Payout(string, *model.AspChannelDepartConfig, *model.AspPayout) (*ThirdPayoutCreateData, *appError.Error)
	PayoutUpi(string, *model.AspChannelDepartConfig, *model.AspPayout) (*ThirdPayoutCreateData, *appError.Error) // 通过upi 方式代付
	Web(string, *model.AspChannelDepartConfig, *model.AspOrder, *req.DeptTradeTypeInfo, *model.AspMerchantProject) bool
	PayNotify()
	PayQuery(string, *model.AspChannelDepartConfig, *model.AspOrder) (*ThirdQueryData, *appError.Error)
	PayoutQuery(string, *model.AspChannelDepartConfig, *model.AspPayout) (*ThirdPayoutQueryData, *appError.Error)
	GetDepartAccountInfo(string, *model.AspChannelDepartConfig) (*ThirdMerchantAccountQueryData, *appError.Error)              // 获取账户信息
	AddBeneficiary(string, string, *model.AspChannelDepartConfig, *req.AspBeneficiary) (*ThirdAddBeneficiary, *appError.Error) // 添加代付受益人
	UpiValidate(string, *model.AspChannelDepartConfig, *req.AspPayoutUpiValidate) (*ThirdUpiValidate, *appError.Error)
}
