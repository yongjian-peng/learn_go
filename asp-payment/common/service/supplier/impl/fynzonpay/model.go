package fynzonpay

import (
	supplier "asp-payment/common/service/supplier/interfaces"
)

// PaymentRsp 分层中的返回 创建收款订单
type PaymentRsp struct {
	StatusCode            int
	StatusMsg             string
	ErrorResponse         *ErrorResponse
	CreateOrderBody       *CreateOrderBody
	CreateBeneficiaryBody *CreateBeneficiaryBody
	CreatePayoutBody      *CreatePayoutBody
	QueryOrderBody        *QueryOrderBody
	QueryPayoutBody       *QueryPayoutBody
}

// GenerateCreateOrder 分层中的赋值 创建收款订单
func (c *PaymentRsp) GenerateCreateOrder(model *supplier.ScanData) {
	model.TransactionID = c.CreateOrderBody.TransactionId
	model.PaymentLink = c.CreateOrderBody.Authurl
	model.Msg = c.StatusMsg
}

// CreateOrderBody 分层中的返回 创建收款订单 上游Response
type CreateOrderBody struct {
	Error         string `json:"Error,omitempty"`
	Status        string `json:"status,omitempty"`
	Message       string `json:"Message,omitempty"`
	Authurl       string `json:"authurl,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
}

// CreatePayoutBody 分层中的返回 创建提现订单 上游Response
type CreatePayoutBody struct {
	AvailableBalance string `json:"available_balance,omitempty"`
	PayoutAmount     string `json:"payout_amount,omitempty"`
	PayoutCurrency   string `json:"payout_currency,omitempty"`
	Reason           string `json:"reason,omitempty"`
	Status           string `json:"status,omitempty"`
	StatusNm         int64  `json:"status_nm,omitempty"`
	TransactionID    string `json:"transaction_id,omitempty"`
}

// CreateBeneficiaryBody
type CreateBeneficiaryBody struct {
	BeneId string `json:"bene_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Remark string `json:"remark,omitempty"`
	Status string `json:"status,omitempty"`
	Notify string `json:"notify,omitempty"`
}

// QueryOrderBody 分层中的返回 查询收款订单 上游Response
type QueryOrderBody struct {
	Error         string `json:"error,omitempty"`
	Amount        string `json:"amount,omitempty"`
	Authurl       string `json:"authurl,omitempty"`
	Cardtype      string `json:"cardtype,omitempty"`
	Curr          string `json:"curr,omitempty"`
	Descriptor    string `json:"descriptor,omitempty"`
	IDOrder       string `json:"id_order,omitempty"`
	Reason        string `json:"reason,omitempty"`
	Status        string `json:"status,omitempty"`
	StatusNm      string `json:"status_nm,omitempty"`
	Tdate         string `json:"tdate,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
}

// QueryPayoutBody 分层中的返回 查询提现订单 上游Response
type QueryPayoutBody struct {
	Status              string `json:"status,omitempty"`
	TransactionType     string `json:"transaction_type,omitempty"`
	TransactionFor      string `json:"transaction_for,omitempty"`
	TransactionDate     string `json:"transaction_date,omitempty"`
	OrderAmount         string `json:"order_amount,omitempty"`
	OrderCurrency       string `json:"order_currency,omitempty"`
	TransactionAmount   string `json:"transaction_amount,omitempty"`
	TransactionCurrency string `json:"transaction_currency,omitempty"`
	MdrAmt              string `json:"mdr_amt,omitempty"`
	MdrPercentage       string `json:"mdr_percentage,omitempty"`
	PayoutAmount        string `json:"payout_amount,omitempty"`
	AvailableBalance    string `json:"available_balance,omitempty"`
	SenderName          string `json:"sender_name,omitempty"`
	BeneficiaryID       string `json:"beneficiary_id,omitempty"`
	Remarks             string `json:"remarks,omitempty"`
	Narration           string `json:"narration,omitempty"`
	TransactionStatus   string `json:"transaction_status,omitempty"`
	NotifyURL           string `json:"notify_url,omitempty"`
	HostName            string `json:"host_name,omitempty"`
}

// ErrorResponse 分层中的返回 错误返回 上游Response
type ErrorResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
