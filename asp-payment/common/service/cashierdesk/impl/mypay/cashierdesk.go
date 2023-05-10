package mypay

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
)

type CashierDeskImpl struct{}

func NewCashierDeskImpl() *CashierDeskImpl {
	return &CashierDeskImpl{}
}

// Assembling 组装参数
func (c *CashierDeskImpl) Assembling(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (model.BodyMap, *appError.Error) {
	bm := make(model.BodyMap)
	redirectUrl := config.AppConfig.Urls.WapPayCheckoutPreUrl + "/" + orderInfo.Sn + "/1"
	bm.Set("redirectUrl", redirectUrl)
	return bm, nil
}

// Rendering 渲染页面参数
// 返回参数 应该有返回的页面地址 和 页面from 表单数据
func (c *CashierDeskImpl) Rendering(bm *model.BodyMap) (string, model.BodyMap, *appError.Error) {
	res := make(map[string]interface{})
	res["redirectUrl"] = bm.Get("redirectUrl")
	return "order/mypay/checkout", res, nil
}

// GetPaymentIntentUrl 获取支付链接地址
func (c *CashierDeskImpl) GetPaymentIntentUrl(orderInfo *model.AspOrder) (string, *appError.Error) {
	if orderInfo.PaymentsUrl == "" {
		return "", nil
	}
	return orderInfo.PaymentsUrl, nil
}

// GetPaymentQrUrl 获取支付链接地址
func (c *CashierDeskImpl) GetPaymentQrUrl(orderInfo *model.AspOrder) (string, *appError.Error) {
	if orderInfo.PaymentsUrl == "" {
		return "", nil
	}
	return orderInfo.PaymentsUrl, nil
}
