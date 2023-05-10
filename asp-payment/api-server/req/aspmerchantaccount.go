package req

//		{
//		    "version":"version",
//	      	"appid":"",
//	     	"sign":"",
//			"depart_id":"depart_id",
//			"channel_id":"channel_id",
//			"email":321456,
//			"password":"",
//		}
//
// AspMerchantAccount 查询商户信息请求参数
type AspMerchantAccount struct {
	DepartId  int    `json:"depart_id" validate:"required"`
	ChannelId int    `json:"channel_id" validate:"required"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type AspMerchantAccountQuery struct {
	Currency string `json:"currency" validate:"required,eq=INR"`
}

type AspMerchantAccountHeaderReq struct {
	Version   string `json:"Version" validate:"oneof=1.0"`
	Timestamp int    `label:"timestamp" json:"timestamp" validate:"required,numeric"`
	AppId     string `json:"AppId" validate:"required" comment:"商户ID"`
	Signature string `json:"Signature" validate:"required" comment:"签名"`
}

// AspMerchantProjectChannelReq 查询cp当前渠道req
type AspMerchantProjectChannelReq struct {
	MchProId  string `json:"mch_pro_id" validate:"required" comment:"商户ID"`
	Currency  string `json:"currency" validate:"required,eq=INR"`
	TradeType string `json:"trade_type" validate:"oneof=H5 PAYOUT"`
}
