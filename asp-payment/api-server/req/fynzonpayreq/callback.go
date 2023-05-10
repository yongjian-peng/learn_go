package fynzonpayreq

// 收款回调 body 返回参数
type CallBackReq struct {
	Amount        string `json:"amount,omitempty" form:"amount,omitempty"`
	Curr          string `json:"curr,omitempty" form:"curr,omitempty"`
	Descriptor    string `json:"descriptor,omitempty" form:"descriptor,omitempty"`
	IDOrder       string `json:"id_order,omitempty" form:"id_order,omitempty"`
	Reason        string `json:"reason,omitempty" form:"reason,omitempty"`
	Status        string `json:"status,omitempty" form:"status,omitempty"`
	StatusNm      string `json:"status_nm,omitempty" form:"status_nm,omitempty"`
	Tdate         string `json:"tdate,omitempty" form:"tdate,omitempty"`
	TransactionID string `json:"transaction_id,omitempty" form:"transaction_id,omitempty"`
}

// 添加受益人回调 body 返回参数
type CallBackBeneficiaryReq struct {
	Reason string `json:"reason,omitempty" form:"reason,omitempty"` // 再次提交 0303
	Status string `json:"status,omitempty" form:"status,omitempty"`
	Notify string `json:"notify,omitempty" form:"notify,omitempty"`
	BeneId string `json:"bene_id,omitempty" form:"bene_id,omitempty"`
}

// 回调 headers 返回 参数
type CallBackHeaderReq struct {
	Signature string `json:"Signature,omitempty"` // Signature Signature=base64(hmac_sha256('回调数据', SecertKey)
}

// 提现回调 返回参数
type CallBackPayoutReq struct {
	Message           string `json:"message,omitempty" form:"message,omitempty"`
	AvailableBalance  string `json:"available_balance,omitempty" form:"available_balance,omitempty"`
	PayoutAmount      string `json:"payout_amount,omitempty" form:"payout_amount,omitempty"`
	PayoutCurrency    string `json:"payout_currency,omitempty" form:"payout_currency,omitempty"`
	Reason            string `json:"reason,omitempty" form:"reason,omitempty"`
	Status            string `json:"status,omitempty" form:"status,omitempty"`
	BankStatus        string `json:"bankStatus,omitempty" form:"bankStatus,omitempty"`
	TransactionStatus string `json:"transaction_status,omitempty" form:"transaction_status,omitempty"`
	TransactionID     string `json:"transaction_id,omitempty" form:"transaction_id,omitempty"`
	RequestId         string `json:"request_id,omitempty" form:"request_id,omitempty"`
	Notify            string `json:"notify,omitempty" form:"notify,omitempty"`
}
