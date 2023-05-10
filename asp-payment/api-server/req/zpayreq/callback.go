package zpayreq

// 收款回调 body 返回参数
type ZPayCallBackReq struct {
	OrderNo        string `form:"orderNo"`        // integer - 平台代收号
	PartnerOrderNo string `form:"partnerOrderNo"` // string - 商户代收号
	ChannelOrderNo string `form:"channelOrderNo"` // string - 上游生成的代收号
	Sign           string `form:"sign"`           // string - 签名值
	ApplicationId  int64  `form:"applicationId"`  // Long - 商户应用ID
	Amount         int    `form:"amount"`         // integer - 代收金额 单位分
	PayWay         int    `form:"payWay"`         // integer - 支付方式ID
	Status         int    `form:"status"`         // Integer - 状态（0：失败 1：成功）
}

// 回调 headers 返回 参数
type ZPayCallBackHeaderReq struct {
	Signature string `json:"Signature,omitempty"` // Signature Signature=base64(hmac_sha256('回调数据', SecertKey)
}

// 提现回调 返回参数
type ZPayPayoutCallBackReq struct {
	ErrorMsg          string `form:"errorMsg"`          // string - 代付上游错误信息（因渠道不同，需要获取到原始值后URLENCODE编码后验签）
	PartnerWithdrawNo string `form:"partnerWithdrawNo"` // string - 商户生成的代付号
	WithdrawNo        string `form:"withdrawNo"`        // string - 平台生成的代付号
	ChannelWithdrawNo string `form:"channelWithdrawNo"` // string - 上游生成的代收号
	Sign              string `form:"sign"`              // string - 签名值
	Amount            int    `form:"amount"`            // integer - 代收金额 单位分
	Status            int    `form:"status"`            // Integer - 状态 0：失败 1：成功
}
