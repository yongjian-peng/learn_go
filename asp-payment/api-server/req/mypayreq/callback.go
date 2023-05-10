package mypayreq

// 代收和代付回调 body 返回参数 用参数来区分是 代收或者代付
type CallBackReq struct {
	Date               string `json:"Date,omitempty"`
	MyPayTransactionID string `json:"MyPayTransactionId,omitempty"`
	Name               string `json:"Name,omitempty"`
	OutletID           int64  `json:"OutletId,omitempty"`
	Upi                string `json:"UPI,omitempty"`
	UserID             int64  `json:"UserId,omitempty"`
	UserTxnID          string `json:"UserTxnId,omitempty"`
	AgentID            string `json:"agent_id,omitempty"`
	Amount             int64  `json:"amount,omitempty"`
	OprID              int64  `json:"opr_id,omitempty"`
	ResCode            string `json:"res_code,omitempty"`
	ResMsg             string `json:"res_msg,omitempty"`
	SpKey              string `json:"sp_key,omitempty"`
	Status             string `json:"status,omitempty"`
}
