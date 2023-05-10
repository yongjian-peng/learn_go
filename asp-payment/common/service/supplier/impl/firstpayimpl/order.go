package firstpayimpl

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"fmt"

	"net/http"
	"time"

	"go.uber.org/zap"
)

// CreateOrder
// H5下单API
// 返回的是 支付的跳转的 payment_link
//
//	Status = 200 is success
func (c *Client) CreateOrder(bm model.BodyMap) (*CreateOrderRsp, *appError.Error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "FirstPay ", "CreateOrder", time.Now(), float64(3))
	headers := c.GetHeader(bm.JsonBody())
	res, bs, err := c.GetHttpServer().PostJson(OrderCreateUrl, headers, bm)
	if err != nil {
		return nil, appError.NewError(err.Error())
	}
	// 和 c++ 作用域 待确认 c++ 则是 函数外面取不到值的
	fpRsq := &CreateOrderRsp{Code: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.Response = new(CreateOrderDetail)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreateOrder ", zap.Error(err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
		fpRsq.Code = res.StatusCode
		fpRsq.Msg = string(bs)
		return fpRsq, nil
	} else {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.Response); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Error(err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
	}
	return fpRsq, nil
}

// 查询订单API
// ctx 全局中的 context
// bm 支付全局 body
// orderNo 支付的订单号 上游返回的
// 返回的是 支付的跳转的 payment_link
//
//	Status = 200 is success
func (c *Client) QueryOrder(bm model.BodyMap, orderNo string) (*QueryOrderRsp, *appError.Error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "FirstPay ", "QueryOrder", time.Now(), float64(3))
	uri := retryCollectCallback + orderNo
	headers := c.GetHeader(bm.JsonBody())
	res, bs, err := c.GetHttpServer().PostJson(uri, headers, bm)
	if err != nil {
		return nil, appError.NewError(err.Error())
	}
	fpRsq := &QueryOrderRsp{Code: Success}
	fpRsq.Response = new(QueryOrderDetail)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.Response); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.QueryOrder.json ", zap.Error(err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
		fpRsq.Code = res.StatusCode
		fpRsq.Msg = string(bs)
		return fpRsq, nil
	} else {

		if err = goutils.JsonDecodeByte(bs, &fpRsq.Response); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.QueryOrder.json ", zap.Error(err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
	}
	return fpRsq, nil
}
