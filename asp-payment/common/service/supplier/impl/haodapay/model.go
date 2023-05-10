package haodapay

import (
	supplier "asp-payment/common/service/supplier/interfaces"
)

// PaymentRsp 分层中的返回 创建收款订单
type PaymentRsp struct {
	StatusCode          int
	StatusMsg           string
	ErrorResponse       *ErrorResponse
	CreateOrderBody     *CreateOrderBody
	CreatePayoutBody    *CreatePayoutBody
	CreatePayoutUpiBody *CreatePayoutUpiBody
	QueryOrderBody      *QueryOrderBody
	QueryPayoutBody     *QueryPayoutBody
	UpiValidateBody     *UpiValidateBody
}

// GenerateCreateOrder 分层中的赋值 创建收款订单
func (c *PaymentRsp) GenerateCreateOrder(model *supplier.ScanData) {
	model.TransactionID = c.CreateOrderBody.CreateOrderData.Reference
	model.PaymentLink = c.CreateOrderBody.CreateOrderData.PaymentLink
	model.Msg = c.StatusMsg
}

// CreateOrderBody 分层中的返回 创建收款订单 上游Response
type CreateOrderBody struct {
	StatusCode      string      `json:"status_code,omitempty"`
	Status          string      `json:"status,omitempty"`
	Message         string      `json:"message,omitempty"`
	Type            string      `json:"type,omitempty"`
	Data            interface{} `json:"data,omitempty"`
	CreateOrderData *CreateOrderData
}

// CreatePayoutBody 分层中的返回 创建提现订单 上游Response
type CreatePayoutBody struct {
	StatusCode string `json:"status_code,omitempty"`
	Status     string `json:"status,omitempty"`
	Message    string `json:"message,omitempty"`
	PayoutID   string `json:"payout_id,omitempty"`
	Code       string `json:"code,omitempty"`
	Type       string `json:"type,omitempty"`
}

// QueryOrderBody 分层中的返回 查询收款订单 上游Response
type QueryOrderBody struct {
	StatusCode     string      `json:"status_code,omitempty"`
	Status         string      `json:"status,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	QueryOrderData *QueryOrderData
}

// QueryPayoutBody 分层中的返回 查询提现订单 上游Response
type QueryPayoutBody struct {
	StatusCode      string      `json:"status_code,omitempty"`
	Status          string      `json:"status,omitempty"`
	Data            interface{} `json:"data,omitempty"`
	QueryPayoutData *QueryPayoutData
	Message         string `json:"message,omitempty"`
	Code            string `json:"code,omitempty"`
	Type            string `json:"type,omitempty"`
}

// UpiValidateBody 分层中的返回 验证 upi 是否合法 上游Response
type UpiValidateBody struct {
	StatusCode   string `json:"status_code,omitempty"`
	Status       string `json:"status,omitempty"`
	CustomerName string `json:"customerName,omitempty"`
	Channel      string `json:"channel,omitempty"`
	Message      string `json:"message,omitempty"`
	Code         string `json:"code,omitempty"`
	Type         string `json:"type,omitempty"`
}

// CreatePayoutUpiBody 分层中的返回 创建提现订单UPI支付方式 上游Response
type CreatePayoutUpiBody struct {
	StatusCode   string `json:"status_code,omitempty"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message,omitempty"`
	PayoutID     string `json:"payout_id,omitempty"`
	Reference    string `json:"reference,omitempty"`
	CustomerName string `json:"customerName,omitempty"`
	Code         string `json:"code,omitempty"`
	Type         string `json:"type,omitempty"`
}

type CreateOrderData struct {
	IntentLink  string `json:"intent_link,omitempty"`
	OrderID     string `json:"order_id,omitempty"`
	PaymentLink string `json:"payment_link,omitempty"`
	QrLink      string `json:"qr_link,omitempty"`
	Reference   string `json:"reference,omitempty"`
}

type QueryOrderData struct {
	UTR string `json:"UTR,omitempty"`
}

type QueryPayoutData struct {
	Status string `json:"status,omitempty"`
	UTR    string `json:"UTR,omitempty"`
}

// ErrorResponse 分层中的返回 错误返回 上游Response
type ErrorResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
