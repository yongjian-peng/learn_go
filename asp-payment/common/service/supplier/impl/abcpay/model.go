package abcpay

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
	model.TransactionID = c.CreateOrderBody.OrderId
	model.PaymentLink = c.CreateOrderBody.PayUrl
	model.Msg = c.StatusMsg
}

// CreateOrderBody 分层中的返回 创建收款订单 上游Response
type CreateOrderBody struct {
	Code    int    `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	PayUrl  string `json:"payUrl,omitempty"`
	OrderId string `json:"orderId,omitempty"`
}

// CreatePayoutBody 分层中的返回 创建提现订单 上游Response
type CreatePayoutBody struct {
	Code    int     `json:"code,omitempty"`
	Msg     string  `json:"msg,omitempty"`
	OrderId string  `json:"order_id,omitempty"`
	Balance float32 `json:"balance,omitempty"`
}

// QueryOrderBody 分层中的返回 查询收款订单 上游Response
type QueryOrderBody struct {
	Code       int    `json:"code,omitempty"`
	Msg        string `json:"msg,omitempty"`
	OrderId    string `json:"order_id,omitempty"`
	FinishTime int    `json:"finish_time,omitempty"`
	Status     int    `json:"status,omitempty"` //订单状态
	OutBizNo   string `json:"out_biz_no,omitempty"`
	Sign       string `json:"sign,omitempty"`
}

type QueryPayoutData struct {
	Status      int    `json:"status,omitempty"`
	SuccessTime int    `json:"success_time,omitempty"`
	OrderId     string `json:"order_id,omitempty"`
}

// QueryPayoutBody 分层中的返回 查询提现订单 上游Response
type QueryPayoutBody struct {
	Msg  string          `json:"msg,omitempty"`
	Code int             `json:"code"`
	Data QueryPayoutData `json:"data,omitempty"`
}

// ErrorResponse 分层中的返回 错误返回 上游Response
type ErrorResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
