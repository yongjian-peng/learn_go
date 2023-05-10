package req

// AspPaymentReq 支付下单请求参数
type AspOrderQuery struct {
	Sn             string `label:"sn" query:"sn" validate:"required"`
	OutTradeNo     string `label:"out_trade_no" query:"out_trade_no, omitempty"`
	IsCallUpstream string `label:"is_call_upstream" query:"is_call_upstream, omitempty"`
}

// OrderAmountList 聚合sum 查询 amount 的值
type OrderAmountList struct {
	Amount int `json:"amount"`
}
