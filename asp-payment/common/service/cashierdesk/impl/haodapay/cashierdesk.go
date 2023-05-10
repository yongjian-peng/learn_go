package haodapay

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/goutils"
)

type CashierDeskImpl struct{}

func NewCashierDeskImpl() *CashierDeskImpl {
	return &CashierDeskImpl{}
}

type CreateOrderData struct {
	IntentLink  string `json:"intent_link,omitempty"`
	OrderID     string `json:"order_id,omitempty"`
	PaymentLink string `json:"payment_link,omitempty"`
	QrLink      string `json:"qr_link,omitempty"`
	Reference   string `json:"reference,omitempty"`
}

// Assembling 组装参数
func (c *CashierDeskImpl) Assembling(requestId string, channelDepartInfo *model.AspChannelDepartConfig, orderInfo *model.AspOrder) (model.BodyMap, *appError.Error) {
	bm := make(model.BodyMap)
	wapPayCheckoutPreUrl := config.AppConfig.Urls.WapPayCheckoutPreUrl
	wapPayCheckoutQrCodeUrl := config.AppConfig.Urls.WapPayCheckoutQrCodeUrl
	if wapPayCheckoutPreUrl == "" || wapPayCheckoutQrCodeUrl == "" {
		return bm, appError.CodeUnknown
	}
	redirectUrl := wapPayCheckoutPreUrl + "/" + orderInfo.Sn + "/1"
	redirectQrUrl := wapPayCheckoutQrCodeUrl + "/" + orderInfo.Sn
	bm.Set("redirectUrl", redirectUrl).
		Set("redirectQrUrl", redirectQrUrl)

	return bm, nil
}

// Rendering 渲染页面参数
// 返回参数 应该有返回的页面路径地址 和 页面 from 表单数据
func (c *CashierDeskImpl) Rendering(bm *model.BodyMap) (string, model.BodyMap, *appError.Error) {
	res := make(map[string]interface{})
	res["redirectUrl"] = bm.Get("redirectUrl")
	res["redirectQrUrl"] = bm.Get("redirectQrUrl")

	return "order/haodapay/checkout", res, nil
}

// GetPaymentIntentUrl 获取支付链接地址
func (c *CashierDeskImpl) GetPaymentIntentUrl(orderInfo *model.AspOrder) (string, *appError.Error) {
	if orderInfo.PaymentsUrl == "" {
		return "", nil
	}
	var createOrderData CreateOrderData
	err := goutils.JsonDecode(orderInfo.PaymentsUrl, &createOrderData)
	if err != nil {
		return "", appError.CodeUnknown
	}
	return createOrderData.IntentLink, nil
}

// GetPaymentQrUrl 获取支付链接地址 QrUrl
func (c *CashierDeskImpl) GetPaymentQrUrl(orderInfo *model.AspOrder) (string, *appError.Error) {
	if orderInfo.PaymentsUrl == "" {
		return "", nil
	}
	var createOrderData CreateOrderData
	err := goutils.JsonDecode(orderInfo.PaymentsUrl, &createOrderData)
	if err != nil {
		return "", appError.CodeUnknown
	}
	return createOrderData.QrLink, nil
}
