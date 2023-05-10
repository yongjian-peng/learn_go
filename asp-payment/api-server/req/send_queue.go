package req

// MerchantProjectTotalFeeQueue 提现查询 headers 头部信息
type MerchantProjectTotalFeeQueue struct {
	ProductType string `json:"product_type"`
	ProductID   int    `json:"product_id"`
}

type MerchantProjectPreFlowQueue []struct {
	MchProjectID         int    `json:"mch_project_id"`
	MchProjectCurrencyID int    `json:"mch_project_currency_id"`
	PreTotalFee          int    `json:"pre_total_fee"`
	ProductType          int    `json:"product_type"`
	ProductID            int    `json:"product_id"`
	Status               int    `json:"status"`
	Currency             string `json:"currency"`
}
