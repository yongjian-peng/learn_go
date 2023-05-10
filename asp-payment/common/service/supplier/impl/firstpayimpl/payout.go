package firstpayimpl

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/pkg/xhttp"
	"fmt"

	"net/http"
	"time"

	"go.uber.org/zap"
)

// H5下单API
// 返回的是 支付的跳转的 payment_link
//
//	Status = 200 is success
func (c *Client) CreatePayout(bm model.BodyMap) (*CreatePayoutRsp, *appError.Error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "FirstPay ", "CreatePayout", time.Now(), float64(3))
	headers := c.GetHeader(bm.JsonBody())
	res, bs, err := c.GetHttpServer().PostJson(orderPayoutURL, headers, bm)
	if err != nil {
		return nil, appError.NewError(err.Error())
	}
	fpRsq := &CreatePayoutRsp{Code: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.Response = new(CreatePayoutDetail)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreatePayout.json ", zap.Error(err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
		fpRsq.Code = res.StatusCode
		fpRsq.Msg = string(bs)
		return fpRsq, nil
	} else {

		if err = goutils.JsonDecodeByte(bs, &fpRsq.Response); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreatePayout.json ", zap.Error(err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
	}
	return fpRsq, nil
}

// QueryPayout
// 查询订单API
// ctx 全局中的 context
// bm 支付全局 body
// payoutNo 支付的订单号 上游返回的
// 返回的是 支付的跳转的 payment_link
//
//	Status = 200 is success
func (c *Client) QueryPayout(TransactionID, OutTradeNo string) (*QueryPayoutRsp, *appError.Error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "FirstPay ", "QueryPayout", time.Now(), float64(3))
	// /v1/platform/inquiry_payout_status?order_id=xfafi212312ada&app_order_id=stextadaf
	uri := fmt.Sprintf(queryPayout, TransactionID) + "&app_order_id=" + OutTradeNo
	// fmt.Println("url-------------", uri)
	// fmt.Println("signature-------------", signature)
	headers := c.GetHeader("")
	res, bs, err := c.GetHttpServer().Get(uri, xhttp.TypeJSON, headers, model.BodyMap{})
	if err != nil {
		return nil, appError.NewError(err.Error())
	}
	fpRsq := &QueryPayoutRsp{Code: Success}
	fpRsq.Response = new(QueryPayoutDetail)
	fpRsq.ErrorResponse = new(ErrorResponse)
	// fmt.Print("res--------------------", res)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.QueryPayout.json ", zap.Error(err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
		fpRsq.Code = res.StatusCode
		fpRsq.Msg = string(bs)
		return fpRsq, nil
	} else {

		if err = goutils.JsonDecodeByte(bs, &fpRsq.Response); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.QueryPayout.json ", zap.Error(err))
			return nil, appError.NewError(fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs)).Error())
		}
		fpRsq.Code = res.StatusCode
		fpRsq.Msg = string(bs)
	}
	return fpRsq, nil
}
