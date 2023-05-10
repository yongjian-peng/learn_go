package firstpayreq

// 收款回调 body 返回参数
type FirstPayCallBackReq struct {
	AppOrderId     string `json:"app_order_id"`               // 请求的聚合支付id sn
	OrderId        string `json:"order_id"`                   // 上游的返回的id OPIN166340390606261994dc
	Amount         int    `json:"amount"`                     // 金额（实际用户支付金额，部分**印度**渠道可能存在实际支付金额和订单金额不等情况，找商务确认）
	PaymentOrderId string `json:"payment_order_id,omitempty"` // 墨西哥电汇支付唯一id(分多次支付的时候才会存在该字段)
	ThirdDesc      string `json:"third_desc"`                 // 墨西哥电汇固定值multiTransferViaClabe(分多次支付的时候才会有该字段)
	Status         int    `json:"status"`                     // 支付状态 0支付中 1支付成功 2支付失败
}

// 回调 headers 返回 参数
type FirstPayCallBackHeaderReq struct {
	Signature string `json:"Signature,omitempty"` // Signature Signature=base64(hmac_sha256('回调数据', SecertKey)
}

// 提现回调 返回参数
type FirstPayPayoutCallBackReq struct {
	AppOrderId string `json:"app_order_id,omitempty"` // 请求的聚合支付id sn  string - 调用方order_id
	OrderId    string `json:"order_id,omitempty"`     // 上游的返回的id OPIN166340390606261994dc
	Amount     int64  `json:"amount,omitempty"`       // 金额（实际用户支付金额，部分**印度**渠道可能存在实际支付金额和订单金额不等情况，找商务确认）
	Status     int64  `json:"status,omitempty"`       // 支付状态 0支付中 1支付成功 2支付失败
}
