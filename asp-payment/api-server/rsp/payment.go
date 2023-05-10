package rsp

// 统一输出初始化
type ScanSuccessData struct {
	Sn            string `json:"sn"`
	OutTradeNo    string `json:"out_trade_no"`
	AppId         string `json:"appid"`
	FeeType       string `json:"fee_type"`
	CashFeeType   string `json:"cash_fee_type"`
	TradeType     string `json:"trade_type"`
	Provider      string `json:"provider"`
	TransactionID string `json:"transaction_id"`
	ClientIP      string `json:"client_ip"`
	OrderToken    string `json:"order_token"`
	CreateTime    uint64 `json:"create_time"`
	FinishTime    uint64 `json:"finish_time"`
	OrderAmount   uint   `json:"order_amount"`
	TotalFee      uint   `json:"total_fee"`
	PaymentLink   string `json:"payment_link"`
	TradeState    string `json:"trade_state"`
	NotifyURL     string `json:"notify_url"`
	ReturnURL     string `json:"return_url"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
}

type QuerySuccessData struct {
	Id          int    `json:"id"`
	Sn          string `json:"sn"`
	OutTradeNo  string `json:"out_trade_no"`
	AppId       string `json:"appid"`
	FeeType     string `json:"fee_type"`
	CashFeeType string `json:"cash_fee_type"`
	TradeType   string `json:"trade_type"`
	Provider    string `json:"provider"`
	// Adapter       string `json:"adapter"`
	TransactionID string `json:"transaction_id"`
	ClientIP      string `json:"client_ip"`
	OrderToken    string `json:"order_token"`
	CreateTime    uint64 `json:"create_time"`
	FinishTime    uint64 `json:"finish_time"`
	OrderAmount   uint   `json:"order_amount"`
	TotalFee      uint   `json:"total_fee"`
	PaymentLink   string `json:"payment_link"`
	TradeState    string `json:"trade_state"`
	// BankType      string `json:"bank_type"`
	ReturnURL     string `json:"return_url"`
	NotifyURL     string `json:"notify_url"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
}

type PayoutSuccessData struct {
	Sn            string `json:"sn"`
	OutTradeNo    string `json:"out_trade_no"`
	AppId         string `json:"appid"`
	FeeType       string `json:"fee_type"`
	CashFeeType   string `json:"cash_fee_type"`
	TradeType     string `json:"trade_type"`
	Provider      string `json:"provider"`
	TransactionID string `json:"transaction_id"`
	ClientIP      string `json:"client_ip"`
	CreateTime    uint64 `json:"create_time"`
	FinishTime    uint64 `json:"finish_time"`
	OrderAmount   uint   `json:"order_amount"`
	TotalFee      uint   `json:"total_fee"`
	BankType      string `json:"bank_type"`
	NotifyURL     string `json:"notify_url"`
	TradeState    string `json:"trade_state"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
}

type PayoutQuerySuccessData struct {
	//Id            int    `json:"id"`
	Sn            string `json:"sn"`
	OutTradeNo    string `json:"out_trade_no"`
	AppId         string `json:"appid"`
	FeeType       string `json:"fee_type"`
	CashFeeType   string `json:"cash_fee_type"`
	TradeType     string `json:"trade_type"`
	Provider      string `json:"provider"`
	TransactionID string `json:"transaction_id"`
	ClientIP      string `json:"client_ip"`
	CreateTime    uint64 `json:"create_time"`
	FinishTime    uint64 `json:"finish_time"`
	OrderAmount   uint   `json:"order_amount"`
	TotalFee      uint   `json:"total_fee"`
	BankType      string `json:"bank_type"`
	NotifyURL     string `json:"notify_url"`
	TradeState    string `json:"trade_state"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
}

type MerchantAccountQuerySuccessData struct {
	Balance          int64 `json:"balance"`           // 可用余额 （可以提供代付和提现的金额）
	BalanceIng       int64 `json:"balance_ing"`       // 提现中金额
	BalanceFreeze    int64 `json:"balance_freeze"`    // 被冻结金额
	BalanceUnsettled int64 `json:"balance_unsettled"` // 待结算余额 未结算的余额
	// BalanceOut int64 `json:"balance_out"` // 已提现金额
	AppId int `json:"appid"`
}

type MerchantAccountChannelQuerySuccessData struct {
	Balance          int64 `json:"balance"`           // 可用余额 （可以提供代付和提现的金额）
	BalanceIng       int64 `json:"balance_ing"`       // 提现中金额
	BalanceFreeze    int64 `json:"balance_freeze"`    // 被冻结金额
	BalanceUnsettled int64 `json:"balance_unsettled"` // 待结算余额 未结算的余额
}

type MerchantProjectChannelSuccessData struct {
	DepartID   int    `json:"depart_id"`
	ChannelID  uint   `json:"channel_id"`
	DepartName string `json:"depart_name"`
	Provider   string `json:"provider"`
	Payment    string `json:"payment"`
	TradeType  string `json:"trade_type"`
	H5Type     string `json:"h5_type"`
}

type ScanFailData struct {
	PayKey     string `json:"payKey"`
	StatusCode string `json:"statusCode"`
	Msg        string `json:"msg"`
}

// data	object	响应结构体
// └─m_id	Number	商户号
// └─balance	Number	可用余额
// └─balance_out	Number	已提现金额
// └─balance_ing	Number	提现中金额
// └─balance_freeze	Number	被冻结金额
// └─sign	string	签名(仅响应结构体参与)
// status	int	响应码
// msg	string	响应消息
