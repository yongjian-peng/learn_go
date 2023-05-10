package seveneightpayreq

// CallBackOrderReq 收款回调 body 返回参数
type CallBackOrderReq struct {
	MemberId      string `form:"memberid" json:"memberid"`             // string - 商户编号
	OrderId       string `form:"orderid" json:"orderid"`               // string - 订单号
	TransactionId string `form:"transaction_id" json:"transaction_id"` // string - 交易号
	Amount        string `form:"amount" json:"amount"`                 // string - 上游生成的代收号
	Sign          string `form:"sign" json:"sign"`                     // string - 签名值
	Datetime      string `form:"datetime" json:"datetime"`             // Long - 商户应用ID
	ReturnCode    string `form:"returncode" json:"returncode"`         // integer - 代收金额 单位分
	Attach        string `form:"attach" json:"attach"`                 // string - 附加参数
}

// CallBackPayoutReq 收款回调 body 返回参数
type CallBackPayoutReq struct {
	MchId         string `form:"mchid,omitempty" json:"mchid,omitempty"`                   // string - 商户编号
	OutTradeNo    string `form:"out_trade_no,omitempty" json:"out_trade_no,omitempty"`     // string - 订单号
	TransactionId string `form:"transaction_id,omitempty" json:"transaction_id,omitempty"` // string - 订单号
	Money         string `form:"money,omitempty" json:"money,omitempty"`                   // string - 订单金额
	Status        string `form:"status,omitempty" json:"status,omitempty"`                 // string - 转账状态
	Msg           string `form:"msg,omitempty" json:"msg,omitempty"`                       // string - 状态说明
	Sign          string `form:"sign,omitempty" json:"sign,omitempty"`                     // string - sign
}
