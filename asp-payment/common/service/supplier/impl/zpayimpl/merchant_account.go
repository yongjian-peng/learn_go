package zpayimpl

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// 查询订单API
// ctx 全局中的 context
// bm 支付全局 body
// orderNo 支付的订单号 上游返回的
// 返回的是 账户的金额
//
//	Status = 200 is success
func (c *Client) QueryMerchantAccount(ctx context.Context, bm model.BodyMap) (*QueryMerchantAccountRsp, *appError.Error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "ZPay ", "QueryMerchantAccount", time.Now(), float64(3))
	signature := c.signature(bm, c.PayoutSerectKey)

	bm.Set("sign", signature)
	res, bs, err := c.doProdGet(ctx, bm, queryMerchantAccount)
	if err != nil {
		return nil, appError.NewError(err.Error())
	}
	fpRsq := &QueryMerchantAccountRsp{Code: http.StatusOK}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.Response = new(QueryMerchantAccountDetail)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.QueryMerchantAccount ", zap.Any("err", err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
		fpRsq.Code = res.StatusCode
		fpRsq.Msg = fpRsq.ErrorResponse.Message
		return fpRsq, nil
	} else {

		if err = goutils.JsonDecodeByte(bs, &fpRsq.Response); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.QueryMerchantAccount ", zap.Any("err", err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
		fpRsq.Msg = fpRsq.Response.Message
	}
	return fpRsq, nil
}
