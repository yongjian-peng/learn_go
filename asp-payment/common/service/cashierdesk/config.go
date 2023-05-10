package cashierdesk

import (
	"asp-payment/common/pkg/constant"
	"asp-payment/common/service/cashierdesk/impl/amarquickpay"
	"asp-payment/common/service/cashierdesk/impl/fynzonpay"
	"asp-payment/common/service/cashierdesk/impl/haodapay"
	"asp-payment/common/service/cashierdesk/impl/mypay"
	"asp-payment/common/service/cashierdesk/interfaces"
)

var mypayCashierDesk = mypay.NewCashierDeskImpl()
var fynzonpayCashierDesk = fynzonpay.NewCashierDeskImpl()
var amarquickpayCashierDesk = amarquickpay.NewCashierDeskImpl()
var haodapayCashierDesk = haodapay.NewCashierDeskImpl()

func GetCashierDeskByCode(code string) interfaces.CashierDeskInterface {
	switch code {
	case constant.TradeTypeMyPay:
		return mypayCashierDesk
	case constant.TradeTypeFynzonPay:
		return fynzonpayCashierDesk
	case constant.TradeTypeAmarquickPay:
		return amarquickpayCashierDesk
	case constant.TradeTypeHaoDaPay:
		return haodapayCashierDesk
	}
	return nil
}
