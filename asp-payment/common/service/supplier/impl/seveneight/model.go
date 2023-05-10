package seveneight

import (
	supplier "asp-payment/common/service/supplier/interfaces"
)

// PaymentRsp 分层中的返回 创建收款订单
type PaymentRsp struct {
	StatusCode       int
	StatusMsg        string
	ErrorResponse    *ErrorResponse
	CreateOrderBody  *CreateOrderBody
	CreatePayoutBody *CreatePayoutBody
	QueryOrderBody   *QueryOrderBody
	QueryPayoutBody  *QueryPayoutBody
}

// GenerateCreateOrder 分层中的赋值 创建收款订单
func (c *PaymentRsp) GenerateCreateOrder(model *supplier.ScanData) {
	model.TransactionID = c.CreateOrderBody.TransactionId
	model.PaymentLink = c.CreateOrderBody.PayUrl
	model.Msg = c.StatusMsg
}

// CreateOrderBody 分层中的返回 创建收款订单 上游Response
type CreateOrderBody struct {
	Status        string `json:"status,omitempty"`
	Msg           string `json:"msg,omitempty"`
	PayUrl        string `json:"payurl,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
}

// CreatePayoutBody 分层中的返回 创建提现订单 上游Response
type CreatePayoutBody struct {
	Status        string `json:"status,omitempty"`
	Msg           string `json:"msg,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
	Balance       string `json:"balance,omitempty"`
}

// QueryOrderBody 分层中的返回 查询收款订单 上游Response
type QueryOrderBody struct {
	Status        string `json:"status,omitempty"`
	Msg           string `json:"msg,omitempty"`
	MemberId      string `json:"memberid,omitempty"`
	OrderId       string `json:"orderid,omitempty"`
	Amount        string `json:"amount,omitempty"`
	TimeEnd       string `json:"time_end,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
	ReturnCode    string `json:"returncode,omitempty"`
	TradeState    string `json:"trade_state,omitempty"`
	Sign          string `json:"sign,omitempty"`
}

// QueryPayoutBody 分层中的返回 查询提现订单 上游Response
type QueryPayoutBody struct {
	Status        string `json:"status,omitempty"`
	Msg           string `json:"msg,omitempty"`
	MchId         string `json:"mchid,omitempty"`
	OutTradeNo    string `json:"out_trade_no,omitempty"`
	Amount        string `json:"amount,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
	RefCode       string `json:"refCode,omitempty"`
	RefMsg        string `json:"refMsg,omitempty"`
	SuccessTime   string `json:"success_time,omitempty"`
}

// ErrorResponse 分层中的返回 错误返回 上游Response
type ErrorResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
