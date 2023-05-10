package zpayimpl

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// H5下单API
// 返回的是 支付的跳转的 payment_link
//
//	Status = 200 is success
func (c *Client) CreatePayout(ctx context.Context, bm model.BodyMap) (zpRsq *CreatePayoutRsp, err error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "ZPay ", "CreatePayout", time.Now(), float64(3))
	signature := c.signature(bm, c.PayoutSerectKey)
	bm.Set("sign", signature)
	res, bs, err := c.doProdGet(ctx, bm, orderPayout)
	if err != nil {
		return nil, err
	}
	zpRsq = &CreatePayoutRsp{Code: http.StatusOK}
	zpRsq.ErrorResponse = new(ErrorResponse)
	zpRsq.Response = new(CreatePayoutDetail)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &zpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		zpRsq.Code = res.StatusCode
		zpRsq.Msg = zpRsq.ErrorResponse.Message
		return zpRsq, nil
	} else {

		if err = goutils.JsonDecodeByte(bs, &zpRsq.Response); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		zpRsq.Msg = zpRsq.Response.Message
	}
	return zpRsq, nil
}

// 查询订单API
// ctx 全局中的 context
// bm 支付全局 body
// payoutNo 支付的订单号 上游返回的
// 返回的是 支付的跳转的 payment_link
//
//	Status = 200 is success
func (c *Client) QueryPayout(ctx context.Context, bm model.BodyMap) (*QueryPayoutRsp, *appError.Error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "ZPay ", "QueryPayout", time.Now(), float64(3))
	signature := c.signature(bm, c.PayoutSerectKey)
	bm.Set("sign", signature)
	res, bs, err := c.doProdGet(ctx, bm, queryPayout)
	if err != nil {
		return nil, appError.NewError(err.Error())
	}
	zpRsq := &QueryPayoutRsp{Code: http.StatusOK}
	zpRsq.Response = new(QueryPayoutDetail)
	zpRsq.ErrorResponse = new(ErrorResponse)
	// fmt.Print("res--------------------", fmt.Sprintf("%+v", res))
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &zpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
		zpRsq.Code = res.StatusCode
		zpRsq.Msg = zpRsq.ErrorResponse.Message
		return zpRsq, nil
	} else {
		if err = goutils.JsonDecodeByte(bs, &zpRsq.Response); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
		result, ok := zpRsq.Response.Data.(string)
		if !ok {
			// fmt.Println("zpRsq.Response.Data:", zpRsq.Response.Data)
			// 转字符串
			jsonbody, err := json.Marshal(zpRsq.Response.Data)
			if err != nil {
				logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
				return nil, appError.NewError(err.Error())
			}
			queryPayoutData := new(QueryPayoutData)
			// 通过字符串转结构体
			if err := goutils.JsonDecodeByte(jsonbody, &queryPayoutData); err != nil {
				logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
				return nil, appError.NewError(err.Error())
			}
			zpRsq.Msg = zpRsq.Response.Message

			zpRsq.Response.QueryPayoutData = queryPayoutData

		} else {
			zpRsq.Msg = result
		}
		zpRsq.Code = res.StatusCode
	}
	return zpRsq, nil
}
