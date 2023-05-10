package interfaces

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
)

type CashierDeskInterface interface {
	// 组装页面参数
	Assembling(string, *model.AspChannelDepartConfig, *model.AspOrder) (model.BodyMap, *appError.Error)
	// 渲染页面参数
	Rendering(*model.BodyMap) (string, model.BodyMap, *appError.Error)

	// 获取支付地址
	GetPaymentIntentUrl(*model.AspOrder) (string, *appError.Error)
	// 获取支付地址
	GetPaymentQrUrl(*model.AspOrder) (string, *appError.Error)
}
