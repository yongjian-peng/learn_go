package mypay

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
	model.TransactionID = c.CreateOrderBody.MyPayTransactionId
	model.PaymentLink = c.CreateOrderBody.Url
	model.Msg = c.CreateOrderBody.Message
}

// CreateOrderBody 分层中的返回 创建收款订单 上游Response
type CreateOrderBody struct {
	Success            string `json:"success,omitempty"`
	Url                string `json:"url,omitempty"`
	Message            string `json:"message,omitempty"`
	MyPayTransactionId string `json:"myPayTransactionId,omitempty"`
}

//type CreateOrderBody struct {
//	StatusCode      string           `json:"statuscode,omitempty"`
//	Status          string           `json:"status,omitempty"`
//	Data            interface{}      `json:"data,omitempty"`
//	CreateOrderData *CreateOrderData `json:"data,omitempty"`
//}

// CreatePayoutBody 分层中的返回 创建提现订单 上游Response
type CreatePayoutBody struct {
	StatusCode       string      `json:"statuscode,omitempty"`
	Status           string      `json:"status,omitempty"`
	Data             interface{} `json:"data,omitempty"`
	CreatePayoutData *CreatePayoutData
}

// QueryOrderBody 分层中的返回 查询收款订单 上游Response
type QueryOrderBody struct {
	StatusCode     string      `json:"statuscode,omitempty"`
	Status         string      `json:"status,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	QueryOrderData *QueryOrderData
}

// QueryPayoutBody 分层中的返回 查询提现订单 上游Response
type QueryPayoutBody struct {
	StatusCode      string      `json:"statuscode,omitempty"`
	Status          string      `json:"status,omitempty"`
	Data            interface{} `json:"data,omitempty"`
	QueryPayoutData *QueryPayoutData
	TxnType         string `json:"txnType,omitempty"`
	Timestamp       string `json:"timestamp,omitempty"`
	MyPayGuid       string `json:"MyPayGuid,omitempty"`
}

// 分层中的返回 创建收款订单 上游Response data
type CreateOrderData struct {
	Status         string `json:"status,omitempty"`         // Success
	RechargeAmount string `json:"rechargeAmount,omitempty"` //
	OperatorRef    int    `json:"operatorRef,omitempty"`    //
	OperatorCode   string `json:"operatorCode,omitempty"`   //
	UserTxnId      string `json:"userTxnId,omitempty"`      //
	MobileNo       string `json:"mobileNo,omitempty"`       //
	MyPayTxnId     string `json:"myPayTxnId,omitempty"`     //
}

// 分层中的返回 查询收款订单 上游Response data
type QueryOrderData struct {
	Status         string `json:"status,omitempty"`         // Success
	RechargeAmount string `json:"rechargeAmount,omitempty"` //
	OperatorRef    int    `json:"operatorRef,omitempty"`    //
	OperatorCode   string `json:"operatorCode,omitempty"`   //
	UserTxnId      string `json:"userTxnId,omitempty"`      //
	MobileNo       string `json:"mobileNo,omitempty"`       //
	MyPayTxnId     string `json:"myPayTxnId,omitempty"`     //
}

// 分层中的返回 创建提现订单 上游Response data
type CreatePayoutData struct {
	UserId             string `json:"userId,omitempty"`             // Success
	UserName           string `json:"userName,omitempty"`           //
	Token              string `json:"token,omitempty"`              //
	IpAdderess         string `json:"ipAdderess,omitempty"`         //
	AccountNumber      string `json:"accountNumber,omitempty"`      //
	IfscCode           string `json:"ifscCode,omitempty"`           //
	BeneficiaryName    string `json:"beneficiaryName,omitempty"`    //
	Amount             string `json:"amount,omitempty"`             //
	MobileNo           string `json:"mobileNo,omitempty"`           //
	Email              string `json:"email,omitempty"`              //
	OperatorId         string `json:"operatorId,omitempty"`         //
	TransactionId      string `json:"transactionId,omitempty"`      //
	MyPayTransactionId string `json:"myPayTransactionId,omitempty"` //
	Status             string `json:"status,omitempty"`             //
	StatusCode         string `json:"statusCode,omitempty"`         //
	StatusMessage      string `json:"statusMessage,omitempty"`      //
	Remark             string `json:"remark,omitempty"`             //
	PayoutId           string `json:"payoutId,omitempty"`           //
	BankRef            string `json:"bankRef,omitempty"`            //
	CreateDate         string `json:"createDate,omitempty"`         //
}

// 分层中的返回 查询提现订单 上游Response data
type QueryPayoutData struct {
	Status         string    `json:"transactionStatusCode,omitempty"`    // Success
	RechargeAmount string    `json:"transactionStatus,omitempty"`        //
	OperatorRef    int       `json:"transactionStatusMessage,omitempty"` //
	OperatorCode   string    `json:"transactionAmount,omitempty"`        //
	UserTxnId      string    `json:"transactionBankRef,omitempty"`       //
	TxnOrder       *TxnOrder //
}

// 分层中返回 查询状态 上游 Response data txnOrder
type TxnOrder struct {
	DealerTxnId        string `json:"dealerTxnId,omitempty"`        //
	MyPayTransactionId string `json:"MyPayTransactionId,omitempty"` //
	SpKey              string `json:"spKey,omitempty"`              //
	Account            string `json:"account,omitempty"`            //
	Optional1          string `json:"optional1,omitempty"`          //
	Optional2          string `json:"optional2,omitempty"`          //
	Optional3          string `json:"optional3,omitempty"`          //
	Optional4          string `json:"optional4,omitempty"`          //
	Optional5          string `json:"optional5,omitempty"`          //
}

// ErrorResponse 分层中的返回 错误返回 上游Response
type ErrorResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
