package req

type DeptTradeTypeInfo struct {
	ID                            int    `json:"id"`
	DepartID                      int    `json:"depart_id"`
	ChannelID                     uint   `json:"channel_id"`
	Provider                      string `json:"provider"`
	Payment                       string `json:"payment"`
	TradeType                     string `json:"trade_type"`
	H5Type                        string `json:"h5_type"`
	InFeeRate                     string `json:"in_fee_rate"`
	OutFeeRate                    string `json:"out_fee_rate"`
	DayUpperLimit                 int    `json:"day_upper_limit"`
	UpperLimit                    int    `json:"upper_limit"`
	LowerLimit                    int    `json:"lower_limit"`
	FixedAmount                   int    `json:"fixed_amount"`
	FixedCurrency                 string `json:"fixed_currency"`
	InFeeRateUpdating             string `json:"in_fee_rate_updating"`
	OutFeeRateUpdating            string `json:"out_fee_rate_updating"`
	Sort                          int    `json:"sort"`
	DepartSort                    int    `json:"depart_sort"`
	DepartMerchantProjectLinkSort int    `json:"depart_merchant_project_link_sort"`
}
