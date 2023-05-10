package amarquickpay

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	amarquickpaySupplier "asp-payment/common/service/supplier/impl/amarquickpay"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type CashierDeskImpl struct{}

func NewCashierDeskImpl() *CashierDeskImpl {
	return &CashierDeskImpl{}
}

type Client struct {
	AppId       string
	SecretKey   string
	RequestId   string
	LogFileName string
	ReturnUrl   string
}

func NewClient(channelDepartInfo *model.AspChannelDepartConfig, requestId, LogFileName string) (*Client, error) {
	var channelConfigInfo model.AspChannelDepartConfigInfo
	goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo)
	CallBack := config.AppConfig.Urls
	if CallBack.AmarquickpayWappayNotifyUrl == "" {
		logger.ApiWarn(LogFileName, requestId, "seveneight notify url err ", zap.Any("CallBack", CallBack))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		err := (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartPaymentCallBackConfigErrMsg)
		return nil, err
	}
	return &Client{
		AppId:       channelConfigInfo.Appid,
		SecretKey:   channelConfigInfo.Signature,
		RequestId:   requestId,
		ReturnUrl:   CallBack.AmarquickpayWappayNotifyUrl,
		LogFileName: LogFileName,
	}, nil
}

// Assembling 组装参数
func (c *CashierDeskImpl) Assembling(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (model.BodyMap, *appError.Error) {
	bm := make(model.BodyMap)

	client, err := NewClient(channelDepartInfo, requestId, constant.AmarquickPayLogFileName)
	if err != nil {
		return bm, appError.CodeSupplierInitClientCode
	}

	amount := cast.ToInt64(orderInfo.TotalFee)
	bm.Set("APP_ID", client.AppId).
		Set("ORDER_ID", orderInfo.Sn).
		Set("RETURN_URL", client.ReturnUrl). //支付成功返回的url地址
		Set("TXNTYPE", "SALE").
		Set("CUST_NAME", orderInfo.CustomerName).
		Set("CUST_PHONE", orderInfo.CustomerPhone).
		Set("CUST_EMAIL", orderInfo.CustomerEmail).
		Set("CURRENCY_CODE", "356").
		Set("CUST_ZIP", "").
		Set("CUST_STREET_ADDRESS1", "").
		Set("AMOUNT", amount) //分
	bm.Set("HASH", amarquickpaySupplier.GetSignature(bm, client.SecretKey))

	return bm, nil
}

// Rendering 渲染页面参数
// 返回参数 应该有返回的页面地址 和 页面from 表单数据
func (c *CashierDeskImpl) Rendering(bm *model.BodyMap) (string, model.BodyMap, *appError.Error) {
	res := make(map[string]interface{})

	res["APP_ID"] = bm.Get("APP_ID")
	res["ORDER_ID"] = bm.Get("ORDER_ID")
	res["RETURN_URL"] = bm.Get("RETURN_URL")
	res["TXNTYPE"] = bm.Get("TXNTYPE")
	res["CUST_NAME"] = bm.Get("CUST_NAME")
	res["CUST_PHONE"] = bm.Get("CUST_PHONE")
	res["CUST_EMAIL"] = bm.Get("CUST_EMAIL")
	res["CURRENCY_CODE"] = bm.Get("CURRENCY_CODE")
	res["AMOUNT"] = bm.Get("AMOUNT")
	res["HASH"] = bm.Get("HASH")

	return "order/checkout", res, nil
}

// GetPaymentIntentUrl 获取支付链接地址
func (c *CashierDeskImpl) GetPaymentIntentUrl(orderInfo *model.AspOrder) (string, *appError.Error) {
	return "", nil
}

// GetPaymentQrUrl 获取支付链接地址
func (c *CashierDeskImpl) GetPaymentQrUrl(orderInfo *model.AspOrder) (string, *appError.Error) {
	return "", nil
}
