package firstpayimpl

import (
	"asp-payment/common/pkg/goutils"
	supplier "asp-payment/common/service/supplier/interfaces"
	"github.com/spf13/cast"
)

// 分层中的返回 创建收款订单
type CreateOrderRsp struct {
	Code          int                `json:"-"`
	Msg           string             `json:"-"`
	ErrorResponse *ErrorResponse     `json:"-"`
	Response      *CreateOrderDetail `json:"response,omitempty"`
}

// 分层中的返回 创建提现订单
type CreatePayoutRsp struct {
	Code          int                 `json:"-"`
	Msg           string              `json:"-"`
	ErrorResponse *ErrorResponse      `json:"-"`
	Response      *CreatePayoutDetail `json:"response,omitempty"`
}

// 分层中的返回 查询收款订单
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

// 分层中的赋值 创建收款订单
func (c *CreateOrderRsp) Generate(model *supplier.ScanData) {
	model.TransactionID = c.Response.CreateOrderData.OrderId
	model.PaymentLink = c.Response.CreateOrderData.PaymentLink
	amount := goutils.Yuan2Fen(cast.ToFloat64(c.Response.CreateOrderData.Amount))
	model.CashFee = cast.ToInt(amount)
	model.Msg = c.Msg
}

// 分层中的赋值 创建提现订单
func (c *CreatePayoutRsp) Generate(model *supplier.ThirdPayoutCreateData) {
	model.TransactionID = c.Response.CreatePayoutData.OrderId
	amount := goutils.Yuan2Fen(cast.ToFloat64(c.Response.CreatePayoutData.Amount))
	model.CashFee = cast.ToInt(amount)
	model.Msg = c.Msg
}

// 分层中的赋值 查询提现订单
func (c *QueryPayoutRsp) Generate(model *supplier.ThirdPayoutQueryData) {
	model.TransactionID = c.Response.QueryPayoutData.AppOrderId
	amount := goutils.Yuan2Fen(cast.ToFloat64(c.Response.QueryPayoutData.Amount))
	model.CashFee = cast.ToInt(amount)
	model.Msg = c.Msg
}

// 分层中的赋值 查询账户详情
func (c *QueryMerchantAccountRsp) Generate(model *supplier.ThirdMerchantAccountQueryData) {
	balance := goutils.Yuan2Fen(cast.ToFloat64(c.Response.QueryPayoutData.Balance))
	unsettledBalance := goutils.Yuan2Fen(cast.ToFloat64(c.Response.QueryPayoutData.UnsettledBalance))
	model.Balance = balance
	model.UnsettledBalance = unsettledBalance
	model.Msg = c.Msg
}

// 分层中的返回 创建收款订单 上游Response
type CreateOrderDetail struct {
	Code            string           `json:"code,omitempty"`
	Msg             string           `json:"msg,omitempty"`
	CreateOrderData *CreateOrderData `json:"data,omitempty"`
}

// 分层中的返回 创建提现订单 上游Response
type CreatePayoutDetail struct {
	Code             string            `json:"code,omitempty"`
	Msg              string            `json:"msg,omitempty"`
	CreatePayoutData *CreatePayoutData `json:"data,omitempty"`
}

// 分层中的返回 查询收款订单 上游Response
type QueryOrderDetail struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// 分层中的返回 查询提现订单 上游Response
type QueryPayoutDetail struct {
	Code            string           `json:"code,omitempty"`
	Msg             string           `json:"msg,omitempty"`
	QueryPayoutData *QueryPayoutData `json:"data,omitempty"`
}

// 分层中的返回 查询账户详情 上游Response
type QueryMerchantAccountDetail struct {
	Code            string                    `json:"code,omitempty"`
	Msg             string                    `json:"msg,omitempty"`
	QueryPayoutData *QueryMerchantAccountData `json:"data,omitempty"`
}

// 分层中的返回 创建收款订单 上游Response data
type CreateOrderData struct {
	AppOrderId  string `json:"app_order_id,omitempty"` // 请求的id
	OrderId     string `json:"order_id,omitempty"`     // 上游的返回的id
	Amount      int    `json:"amount,omitempty"`       // 支付的金额
	PaymentLink string `json:"payment_link,omitempty"` // 支付的链接 Url
}

// 分层中的返回 创建提现订单 上游Response data
type CreatePayoutData struct {
	AppOrderId string `json:"app_order_id"` // string - 调用方流水id，唯一标识
	OrderId    string `json:"order_id"`     // string - 平台流水id
	Amount     int    `json:"amount"`       // integer - 金额
	Status     int    `json:"status"`       // integer - 状态：0-支付中，1-成功，2-失败
}

// 分层中的返回 查询提现订单 上游Response data
type QueryPayoutData struct {
	AppOrderId string `json:"app_order_id"` // string - 调用方流水id，唯一标识
	OrderId    string `json:"order_id"`     // string - 平台流水id
	Amount     int    `json:"amount"`       // integer - 金额
	Status     int    `json:"status"`       // integer - 状态：0-支付中，1-成功，2-失败
}

// 分层中的返回 查询账户详情 上游Response data
type QueryMerchantAccountData struct {
	Balance          int64 `json:"balance"`          // 账户余额 账户总额 = 余额 + 未结算的余额
	UnsettledBalance int64 `json:"unsettledBalance"` // 未结算余额
}

// 分层中的返回 错误返回 上游Response
type ErrorResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
