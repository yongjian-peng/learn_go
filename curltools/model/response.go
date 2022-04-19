package model

type ResponseData struct {
	Code int64 `json:"code"`
	Data struct {
		Adapter        string `json:"adapter"`
		Appid          int64  `json:"appid"`
		Attach         string `json:"attach"`
		BankType       string `json:"bank_type"`
		Body           string `json:"body"`
		CashFee        string `json:"cash_fee"`
		CashFeeType    string `json:"cash_fee_type"`
		CreateTime     string `json:"create_time"`
		Detail         string `json:"detail"`
		Discount       string `json:"discount"`
		FeeType        string `json:"fee_type"`
		ID             string `json:"id"`
		IsPrint        string `json:"is_print"`
		IsSubscribe    string `json:"is_subscribe"`
		MchName        string `json:"mch_name"`
		MwebURL        string `json:"mweb_url"`
		NonceStr       string `json:"nonce_str"`
		OutTradeNo     string `json:"out_trade_no"`
		PayAmount      string `json:"pay_amount"`
		Provider       string `json:"provider"`
		Qrcode         string `json:"qrcode"`
		Refundable     int64  `json:"refundable"`
		Sign           string `json:"sign"`
		Sn             string `json:"sn"`
		TimeEnd        int64  `json:"time_end"`
		TotalFee       string `json:"total_fee"`
		TotalRefundFee int64  `json:"total_refund_fee"`
		TradeState     string `json:"trade_state"`
		TradeType      string `json:"trade_type"`
		TransactionID  string `json:"transaction_id"`
		Version        string `json:"version"`
		Wallet         string `json:"wallet"`
	} `json:"data"`
	Message string `json:"message"`
}
