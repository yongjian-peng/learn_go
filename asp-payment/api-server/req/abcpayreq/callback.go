package abcpayreq

// CallBackOrderReq 收款回调 body 返回参数
type CallBackOrderReq struct {
	Appid    string `form:"appid" json:"appid"`           // string - 商户号，参于签名
	Status   string `form:"status" json:"status"`         // int - 订单状态，1：支付成功；其它失败，参与签名 7 退款中 8 退款成功 9 退款失败，支付中的订单才可能退款
	Money    string `form:"money" json:"money"`           // float32 - 支付金额，单位：元，参与签名
	OrderId  string `form:"order_id" json:"order_id"`     // string - abcpay单号，参与签名
	OutBizNo string `form:"out_biz_no" json:"out_biz_no"` // string - 商户订单号，参与签名
	Sign     string `form:"sign" json:"sign"`             // string - 签名
}

// CallBackPayoutReq 收款回调 body 返回参数
type CallBackPayoutReq struct {
	AppId       string `form:"appid" json:"appid"`             // string - 商户号，参于签名
	Orderstatus string `form:"orderstatus" json:"orderstatus"` // int - 订单状态，1：转账成功；2：转账失败，参与签名
	Endtime     string `form:"endtime" json:"endtime"`         // int - 转账时间戳，参与签名
	Amount      string `form:"amount" json:"amount"`           // float32 - 转账金额，参与签名
	OutBizNo    string `form:"out_biz_no" json:"out_biz_no"`   // string - 商户订单号
	PayType     string `form:"pay_type" json:"pay_type"`       // string - 商户订单号
	Sign        string `form:"sign" json:"sign"`               // string - 签名
}
