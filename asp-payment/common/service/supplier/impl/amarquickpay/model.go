package amarquickpay

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
	CustName          string `json:"CUST_NAME,omitempty"`
	Txntype           string `json:"TXNTYPE,omitempty"`
	Amount            string `json:"AMOUNT,omitempty"`
	CurrencyCode      string `json:"CURRENCY_CODE,omitempty"`
	OrderId           string `json:"ORDER_ID,omitempty"`
	AppId             string `json:"APP_ID,omitempty"`
	TrxnId            string `json:"TRXN_ID,omitempty"`
	PaymentType       string `json:"PAYMENT_TYPE,omitempty"`
	MopType           string `json:"MOP_TYPE,omitempty"`
	CardMask          string `json:"CARD_MASK,omitempty"`
	PgRefNum          string `json:"PG_REF_NUM,omitempty"`
	ResponseCode      string `json:"RESPONSE_CODE,omitempty"`
	ResponseMessage   string `json:"RESPONSE_MESSAGE,omitempty"`
	Hash              string `json:"HASH,omitempty"`
	Eci               string `json:"ECI,omitempty"`
	AuthCode          string `json:"AUTH_CODE,omitempty"`
	Rrn               string `json:"RRN,omitempty"`
	Avr               string `json:"AVR,omitempty"`
	AcqId             string `json:"ACQ_ID,omitempty"`
	Status            string `json:"STATUS,omitempty"`
	CustEmail         string `json:"CUST_EMAIL,omitempty"`
	CustId            string `json:"CUST_ID,omitempty"`
	CustPhone         string `json:"CUST_PHONE,omitempty"`
	PgTxnMessage      string `json:"PG_TXN_MESSAGE,omitempty"`
	ReturnUrl         string `json:"RETURN_URL,omitempty"`
	ResponseDate      string `json:"RESPONSE_DATE,omitempty"`
	ResponseTime      string `json:"RESPONSE_TIME,omitempty"`
	ProductDesc       string `json:"PRODUCT_DESC,omitempty"`
	CardIssuerBank    string `json:"CARD_ISSUER_BANK,omitempty"`
	CardIssuerCountry string `json:"CARD_ISSUER_COUNTRY,omitempty"`
	TotalAmount       string `json:"TOTAL_AMOUNT,omitempty"`
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
