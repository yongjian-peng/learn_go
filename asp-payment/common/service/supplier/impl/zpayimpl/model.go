package zpayimpl

import (
	supplier "asp-payment/common/service/supplier/interfaces"
)

// 分层中的返回 创建代收订单
type CreateOrderRsp struct {
	Code          int                `json:"code"`
	Msg           string             `json:"msg"`
	ErrorResponse *ErrorResponse     `json:"-"`
	Response      *CreateOrderDetail `json:"response,omitempty"`
}

// 分层中的返回 创建提现订单
type CreatePayoutRsp struct {
	Code          int                 `json:"code"`
	Msg           string              `json:"msg"`
	ErrorResponse *ErrorResponse      `json:"-"`
	Response      *CreatePayoutDetail `json:"response,omitempty"`
}

// 分层中的返回 查询代收订单
type QueryOrderRsp struct {
	Code          int               `json:"code"`
	Msg           string            `json:"msg"`
	ErrorResponse *ErrorResponse    `json:"-"`
	Response      *QueryOrderDetail `json:"response,omitempty"`
}

// 分层中的返回 查询提现订单
type QueryPayoutRsp struct {
	Code          int                `json:"code"`
	Msg           string             `json:"msg"`
	ErrorResponse *ErrorResponse     `json:"-"`
	Response      *QueryPayoutDetail `json:"response,omitempty"`
}

// 分层中的返回 查询账户详情
type QueryMerchantAccountRsp struct {
	Code          int                         `json:"code"`
	Msg           string                      `json:"msg"`
	ErrorResponse *ErrorResponse              `json:"-"`
	Response      *QueryMerchantAccountDetail `json:"response,omitempty"`
}

// 分层中的赋值 创建代收订单
func (c *CreateOrderRsp) Generate(model *supplier.ScanData) {
	model.PaymentLink = c.Response.Data
	model.Msg = c.Msg
}

// 分层中的赋值 查询代收订单
func (c *QueryOrderRsp) Generate(model *supplier.ThirdQueryData) {
	model.TransactionID = c.Response.QueryOrderData.OrderNo
	model.Msg = c.Msg
}

// 分层中的赋值 创建提现订单
func (c *CreatePayoutRsp) Generate(model *supplier.ThirdPayoutCreateData) {
	// model.TransactionID = c.Response.CreatePayoutData.OrderId
	// model.CashFee = c.Response.CreatePayoutData.Amount
	model.Msg = c.Msg
}

// 分层中的赋值 查询提现订单
func (c *QueryPayoutRsp) Generate(model *supplier.ThirdPayoutQueryData) {
	// model.TransactionID = c.Response.QueryPayoutData.WithdrawNo
	// model.CashFee = c.Response.QueryPayoutData.Amount
	model.Msg = c.Msg
}

// 分层中的赋值 查询账户详情
func (c *QueryMerchantAccountRsp) Generate(model *supplier.ThirdMerchantAccountQueryData) {
	model.Balance = int64(c.Response.QueryPayoutData.AvailableAmount * 100)
	model.UnsettledBalance = 0
	model.Msg = c.Msg
}

// 分层中的返回 创建代收订单 上游Response
// 分层中的返回 创建代收订单 上游Response data
type CreateOrderDetail struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

// 分层中的返回 创建提现订单 上游Response
// 分层中的返回 创建提现订单 上游Response data

type CreatePayoutDetail struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

// 分层中的返回 查询代收订单 上游Response
type QueryOrderDetail struct {
	Code           string          `json:"code"`
	Message        string          `json:"message"`
	Data           interface{}     `json:"data,omitempty"`
	QueryOrderData *QueryOrderData //返回成功时才赋值
}

// 分层中的返回 查询提现订单 上游Response
type QueryPayoutDetail struct {
	Code            string           `json:"code,omitempty"`
	Message         string           `json:"message,omitempty"`
	Data            interface{}      `json:"data,omitempty"`
	QueryPayoutData *QueryPayoutData //返回成功时才赋值
}

// 分层中的返回 查询账户详情 上游Response
type QueryMerchantAccountDetail struct {
	Code            string                    `json:"code,omitempty"`
	Message         string                    `json:"message,omitempty"`
	QueryPayoutData *QueryMerchantAccountData `json:"data,omitempty"`
}

// 分层中的返回 查询代收订单 上游Response data
type QueryOrderData struct {
	PartnerId      int64  `json:"partnerId"`      // Long - 商户ID 支付中心分配的商户号
	ApplicationId  int64  `json:"applicationId"`  // Long - 商户应用ID
	OrderNo        string `json:"orderNo"`        // integer - 平台代收号
	PartnerOrderNo string `json:"partnerOrderNo"` // string - 商户代收号
	Amount         int    `json:"amount"`         // integer - 代收金额 单位分
	Status         int    `json:"status"`         // Integer - 状态（0：创建订单；1：代收中；2：代收成功；3：代收失败）
	CreateTime     string `json:"createTime"`     // Date - 创建时间
	SuccessTime    string `json:"successTime"`    // Date - 成功时间
}

// 分层中的返回 查询提现订单 上游Response data
type QueryPayoutData struct {
	PayPartnerId      int64  `json:"payPartnerId"`      // Long - 支付中心分配的商户号
	WithdrawNo        string `json:"withdrawNo"`        // string - 平台生成的代付号
	PartnerWithdrawNo string `json:"partnerWithdrawNo"` // string - 商户生成的代付号
	Amount            int    `json:"amount"`            // integer - 代收金额
	Status            int    `json:"status"`            // Integer - 状态（0：创建订单；1：代付中；2：代付成功；3：代付失败）
	CreateTime        string `json:"createTime"`        // Date - 创建时间
	SuccessTime       string `json:"successTime"`       // Date - 成功时间
}

// 分层中的返回 查询账户详情 上游Response data
type QueryMerchantAccountData struct {
	PartnerId       int     `json:"partnerId"`       // 商户号
	AvailableAmount float64 `json:"availableAmount"` // 可用金额
}

// 分层中的返回 错误返回 上游Response
type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}
