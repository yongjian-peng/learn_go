package fynzonpay

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	fynzonpaySupplier "asp-payment/common/service/supplier/impl/fynzonpay"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type CashierDeskImpl struct{}

func NewCashierDeskImpl() *CashierDeskImpl {
	return &CashierDeskImpl{}
}

type Client struct {
	AppId             string
	SecretKey         string
	PartnerId         string
	ApplicationId     string
	PayoutSignature   string
	RequestId         string
	LogFileName       string
	Notify            string
	BeneficiaryNotify string
	PayoutNotify      string
}

func NewClient(channelDepartInfo *model.AspChannelDepartConfig, requestId, LogFileName string) (*Client, error) {
	// {"appid":"1515","partnerId":"53bb71c1edeea8e1258ae678472d2f6feb1dc8283a03f0e7f6efd019f398c17e","applicationId":"anlZYlVJekdPdlFSVDJkeVdma0hlNEQ0VjhrZk9WZk5WVWk2V2VmcFQwOHBHWFdsQm9xeWJBeTRDZS96WVFUYQ","signature":"MTEyOTNfMTUxNV8yMDIzMDExMzAwMTMzOA","payoutSignature":"Qc_9Z2pinBs@#Am"}
	var channelConfigInfo model.AspChannelDepartConfigInfo
	goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo)
	CallBack := config.AppConfig.Urls
	if CallBack.FynzonpayWappayNotifyUrl == "" || CallBack.FynzonpayWappayErrorUrl == "" || CallBack.FynzonpayPayoutNotifyUrl == "" || CallBack.FynzonpayBeneficiaryNotifyUrl == "" || CallBack.FynzonpayOrderReturnUrl == "" {
		logger.ApiWarn(constant.FynzonpayLogFileName, requestId, "fynzonpay notify url err ", zap.Any("CallBack", CallBack))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		err := (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartPaymentCallBackConfigErrMsg)
		return nil, err
	}
	return &Client{
		AppId:             channelConfigInfo.Appid,
		SecretKey:         channelConfigInfo.Signature,
		PartnerId:         channelConfigInfo.PartnerId,
		ApplicationId:     channelConfigInfo.ApplicationId,
		PayoutSignature:   channelConfigInfo.PayoutSignature,
		RequestId:         requestId,
		Notify:            CallBack.FynzonpayWappayNotifyUrl,
		BeneficiaryNotify: CallBack.FynzonpayBeneficiaryNotifyUrl,
		PayoutNotify:      CallBack.FynzonpayPayoutNotifyUrl,
		LogFileName:       LogFileName,
	}, nil
}

// Assembling 组装参数
func (c *CashierDeskImpl) Assembling(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (model.BodyMap, *appError.Error) {
	bm := make(model.BodyMap)
	client, err := NewClient(channelDepartInfo, requestId, constant.FynzonpayLogFileName)
	if err != nil {
		return bm, appError.CodeSupplierInitClientCode
	}

	productName := goutils.IfString(orderInfo.Body == "", orderInfo.Body, fynzonpaySupplier.ProductName)

	//return scanData
	amount := goutils.Fen2Yuan(cast.ToInt64(orderInfo.TotalFee))
	bm.Set("api_token", client.SecretKey).
		Set("store_id", client.AppId).
		Set("cardsend", fynzonpaySupplier.CardSend2).
		Set("client_ip", orderInfo.ClientIp).
		Set("action", fynzonpaySupplier.Action).
		Set("source", "Host-Redirect-Card-Payment (Core PHP)").
		Set("source_url", "https://apiprod.sunnypay.net").
		Set("price", amount).
		Set("curr", fynzonpaySupplier.Curr).
		Set("product_name", productName).
		Set("fullname", orderInfo.CustomerName).
		Set("email", orderInfo.CustomerEmail).
		Set("bill_street_1", "A97B North Block").
		Set("bill_street_2", "West Vinod Nagar").
		Set("bill_city", "New Delhi").
		Set("bill_state", "DL").
		Set("bill_country", "IND").
		Set("bill_zip", "110092").
		Set("bill_phone", orderInfo.CustomerPhone).
		Set("id_order", orderInfo.Sn).
		Set("notify_url", config.AppConfig.Urls.FynzonpayWappayNotifyUrl).
		Set("success_url", config.AppConfig.Urls.FynzonpayOrderReturnUrl).
		Set("error_url", config.AppConfig.Urls.FynzonpayWappayErrorUrl)

	// api_token=MTEyOTNfMTUxM18yMDIzMDExMjE5NTkzMQ&cardsend=CHECKOUT&client_ip=192.168.147.1&action=product&source=Encode-Checkout&source_url=https%3A%2F%2Fbrightsgaming.com&price=30.00&curr=INR&product_name=Testing+Product&fullname=DEV+PRAKASH+YADAV&email=ericluzhonghua%40gmail.com&bill_street_1=A97B+North+Block&bill_street_2=West+Vinod+Nagar&bill_city=New+Delhi&bill_state=DL&bill_country=IN&bill_zip=110092&bill_phone=9036830689&id_order=202202211811¬ify_url=https%3A%2F%2Fyourdomain.com%2Fnotify.php&success_url=https%3A%2F%2Fyourdomain.com%2Fsuccess.php&error_url=https%3A%2F%2Fyourdomain.com%2Ffailed.php&checkout_url=https%3A%2F%2Fyourdomain.com%2Fcheckout_url.php&ccno=5555555555554444&ccvv=123&month=01&year=30¬es=Remark+for+transaction"

	// ZE1XTEpCeW5lWWlHU09lUFB0ZnpUN2dDL1FhRnRKb2lIVHZia3Q0cDIxdG9wQnQ2SnFkMEl3eVpqMU05amx4dGdNcUhTaUYxNkZQckQzOTdYNFZxb0lzcTN1SVNEdkkxaHlMNmlybjArUUVMU0o5VjhOL0tRejFMa3lEcVFsTU1BWVoreXY0NERLUzBLM1MzR21oWTFmSXE0RW1QTW90ZDJLemtLTHRKSkR4blpUb1V5UXhpeFBjeXFZQ1JFcTlXR1FtMXRYSTFqQituK1daTXo2VjhuWUZEeVNDU01RR2p4VFl6aDlUcC9vNHdHYzdDS29zTVZTOVI4QW9rZFRPWWFETmxvVTNyb1hYNEZmZm9uN3ljTWJOdjR0NG5KWm84UXRTWVFWeDZ2RmZkdEtpK0cyZWw1N29oOWpFMWtWNGxGeW9samdjZ3ZkYkQxc1lUZWIvZ2p1NTVuQWNDZGV3bUIrdjM0R2JoQUhMWEZIYzVOZVR2SC9oM0Frblo5M2E1UkVTTkZWSHRaS0ZqZmtBOXJEMmpWV2g5eVIwVlVLT0xaazNFVDJia2JLaTZHR3NnRDE5eUNTTVJ3NVhlcm5IV1BJUXhyWW9FUGxXblJxRFhmQmppZ3BvRFV5ZDRhMGxoZlhnMFRURFZUaFBGZ1M5Wmk4ak1YSHJsMjQzWjhHK3M0ZXpNeWJwd2lVbUtRRXF6UEQ1Y3JZNWZobCtkQnJRbVkzdytLakxGY3BQZFFnWFJXd0VNZmY1ZlAyMjhUNkZIU0NldEJ4WWMvcHNnbHh4TU4vNm1nZG8rSklkT2xzZFpxNzZwQjQ5N0cyc1FDcFd0c0xKWHJHdW8yaUQ5cU14ZWhXYmhCSE1ReWdRdnRJNUhKYzZNLzBiUXM5N2RDejF6SVU1T0ZqbzJzWW94S3Q3SnJydjRSZE1EM1ZHQ00zR1ZGeHVDZEVSTWphQnJkaFpGNmp0TmVBVFRQU0Vncll1RmZUTGJNZzFVYTB3UDJock51MDIwM1o3R2Y1Q1Y5ZXlEajRNclRSaUpYU2xUZUd1aGF2UE9BZ3ZKNzVIMmtFcGdEcnlTWU5kTEl1VzREVmh0ZHhkZHZxTkNFL3p3YWNqQlB0amNtR0RHUy9FME9Ma1dCa0pKWmtzc3c2ZC9teFUzQU9vNEY4N3lVYVVkL3ZOb08xZXN1RXg3ZmVTTTBTTVFjQWNIYkdWQmpKYUVvdHJ6YklnNHBCUXBFTVBWREw2TVlvS0dHME9IZklQdW92WlcyZW1tUVA4UzZwM0tUM1hsbkl6ejRXckEvRi83bTV5V3Z2bVZJSDdjdGpkdjUxY29HTTV0UjFIZG5VUT0"
	param := bm.EncodeURLParams()
	url := fynzonpaySupplier.OrderCreateUrl + "?" + param
	bm.Set("url", url)
	return bm, nil
}

// Rendering 渲染页面参数
// 返回参数 应该有返回的页面地址 和 页面from 表单数据
func (c *CashierDeskImpl) Rendering(bm *model.BodyMap) (string, model.BodyMap, *appError.Error) {
	res := make(map[string]interface{})

	res["api_token"] = bm.Get("api_token")
	res["store_id"] = bm.Get("store_id")
	res["action"] = bm.Get("action")
	res["client_ip"] = bm.Get("client_ip")
	res["source_url"] = bm.Get("source_url")
	res["price"] = bm.Get("price")
	res["curr"] = bm.Get("curr")
	res["product_name"] = bm.Get("product_name")
	res["fullname"] = bm.Get("fullname")
	res["email"] = bm.Get("email")
	res["bill_street_1"] = bm.Get("bill_street_1")
	res["bill_street_2"] = bm.Get("bill_street_2")
	res["bill_city"] = bm.Get("bill_city")
	res["bill_state"] = bm.Get("bill_state")
	res["bill_country"] = bm.Get("bill_country")
	res["bill_zip"] = bm.Get("bill_zip")
	res["bill_phone"] = bm.Get("bill_phone")
	res["id_order"] = bm.Get("id_order")
	res["notify_url"] = bm.Get("notify_url")
	res["success_url"] = bm.Get("success_url")
	res["cardsend"] = bm.Get("cardsend")
	res["source"] = bm.Get("source")

	return "order/fynzonpay/checkout", res, nil
}

// GetPaymentIntentUrl 获取支付链接地址
func (c *CashierDeskImpl) GetPaymentIntentUrl(orderInfo *model.AspOrder) (string, *appError.Error) {

	return "", nil
}

// GetPaymentQrUrl 获取支付链接地址
func (c *CashierDeskImpl) GetPaymentQrUrl(orderInfo *model.AspOrder) (string, *appError.Error) {
	return "", nil
}
