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
func (c *Client) CreateOrder(ctx context.Context, bm model.BodyMap) (*CreateOrderRsp, *appError.Error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "ZPay ", "CreateOrder", time.Now(), float64(3))
	//fmt.Println("bm: ", bm)
	//fmt.Println("SerectKey: ", c.SerectKey)
	signature := c.signature(bm, c.SerectKey)
	bm.Set("sign", signature)
	//fmt.Println("bm: ", bm)
	res, bs, err := c.doProdGet(ctx, bm, orderCreate)
	if err != nil {
		return nil, appError.NewError(err.Error())
	}
	// 和 c++ 作用域 待确认 c++ 则是 函数外面取不到值的
	zpRsq := &CreateOrderRsp{Code: http.StatusOK}
	zpRsq.ErrorResponse = new(ErrorResponse)
	zpRsq.Response = new(CreateOrderDetail)
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
		zpRsq.Msg = zpRsq.Response.Message
	}
	return zpRsq, nil
}

// 查询订单API
// ctx 全局中的 context
// bm 支付全局 body
// orderNo 支付的订单号 上游返回的
// 返回的是 支付的跳转的 payment_link
//
//	Status = 200 is success
func (c *Client) QueryOrder(ctx context.Context, bm model.BodyMap) (zpRsq *QueryOrderRsp, err error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "ZPay ", "QueryOrder", time.Now(), float64(3))
	signature := c.signature(bm, c.SerectKey)
	bm.Set("sign", signature)
	res, bs, err := c.doProdGet(ctx, bm, queryOrder)
	if err != nil {
		return nil, err
	}

	zpRsq = &QueryOrderRsp{Code: http.StatusOK}
	zpRsq.ErrorResponse = new(ErrorResponse)
	zpRsq.Response = new(QueryOrderDetail)
	zpRsq.Code = res.StatusCode
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &zpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		zpRsq.Msg = zpRsq.ErrorResponse.Message
		return zpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &zpRsq.Response); err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
		return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
	}

	result, ok := zpRsq.Response.Data.(string)

	if !ok {
		// fmt.Println("fpRsq.Response.Data:", fpRsq.Response.Data)
		// 转字符串
		jsonbody, err := json.Marshal(zpRsq.Response.Data)
		if err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
			return nil, err
		}
		queryOrderData := new(QueryOrderData)
		// 通过字符串转结构体
		if err := goutils.JsonDecodeByte(jsonbody, &queryOrderData); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("err", err))
			return nil, err
		}
		zpRsq.Msg = zpRsq.Response.Message

		zpRsq.Response.QueryOrderData = queryOrderData

	} else {
		zpRsq.Msg = result
	}
	return zpRsq, nil
}
