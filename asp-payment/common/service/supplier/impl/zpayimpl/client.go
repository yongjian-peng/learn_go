package zpayimpl

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/pkg/xhttp"
	"context"
	"net/http"
	"net/http/httputil"
	"sync"

	"go.uber.org/zap"
)

type Client struct {
	PartnerId        string
	ApplicationId    string
	SerectKey        string
	PayoutSerectKey  string
	ZpayPayoutNotify string
	ZpayH5Notify     string
	mu               sync.RWMutex
	RequestId        string
	LogFileName      string
}

// 初始化First Pay
//
//	appId：应用ID
//	mchId：商户ID
//	ApiKey：API秘钥值
//	IsProd：是否是正式环境
func NewClient(request_id, PartnerId, ApplicationId, SerectKey, PayoutSerectKey string) (client *Client, err error) {

	// config.ExtendConfig = &ext.ExtConfig
	// 注入配置扩展项
	CallBack := config.AppConfig.Urls
	if CallBack.ZpayPayoutNotifyUrl == "" || CallBack.ZpayH5NotifyUrl == "" {
		logger.ApiWarn(constant.FirstPayLogFileName, request_id, "Zpay notify url err ", zap.Any("CallBack", CallBack))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		err = (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartPaymentCallBackConfigErrMsg)
		return client, err
	}
	return &Client{
		PartnerId:        PartnerId,
		ApplicationId:    ApplicationId,
		SerectKey:        SerectKey,
		PayoutSerectKey:  PayoutSerectKey,
		ZpayH5Notify:     CallBack.ZpayH5NotifyUrl,
		ZpayPayoutNotify: CallBack.ZpayPayoutNotifyUrl,
	}, nil
}

func (c *Client) doProdPost(ctx context.Context, bm model.BodyMap, uri string) (res *http.Response, bs []byte, err error) {
	var url = baseUrlProd + uri
	httpClient := xhttp.NewClient()
	logger.ApiWarn(c.LogFileName, c.RequestId, "client.doProdPost Request ",

		zap.String("request", bm.JsonBody()),
		zap.String("PartnerId", c.PartnerId),
		zap.String("ApplicationId", c.ApplicationId),
		zap.String("url", url))

	httpClient.Header.Add("Content-Type", "application/json; charset=utf-8")
	res, bs, err = httpClient.Type(xhttp.TypeJSON).Post(url).SendBodyMap(bm).EndBytes(ctx)
	if err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "httpClient.client.Post err ", zap.String("response", string(bs)), zap.Error(err))
		return nil, nil, err
	}

	// 读取 response 的信息 并且记录到日志
	dump, err := httputil.DumpResponse(res, false) // better way
	if err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "client.doProdPost Response httputil.DumpResponse err ", zap.String("response", string(bs)), zap.Error(err))
		return res, bs, nil
	}
	logger.ApiWarn(c.LogFileName, c.RequestId, "client.doProdPost Response ", zap.String("response", string(bs)), zap.String("res", string(dump)))
	// xlog.Debugf("Wechat_Response: %d > %s", res.StatusCode, string(bs))
	// xlog.Debugf("Wechat_Headers: %#v", res.Header)
	// xlog.Debugf("Wechat_SignInfo: %#v", si)
	return res, bs, nil
}

func (c *Client) doProdGet(ctx context.Context, bm model.BodyMap, uri string) (res *http.Response, bs []byte, err error) {
	var url = baseUrlProd + uri

	param := bm.EncodeURLParams()
	url = url + "?" + param

	httpClient := xhttp.NewClient()
	logger.ApiWarn(c.LogFileName, c.RequestId, "client.doProdGet Request ",
		zap.String("request", bm.JsonBody()),
		zap.String("PartnerId", c.PartnerId),
		zap.String("ApplicationId", c.ApplicationId),
		zap.String("url", url))

	httpClient.Header.Add("Content-Type", "application/json; charset=utf-8")
	res, bs, err = httpClient.Get(url).SendBodyMap(bm).EndBytes(ctx)
	if err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "httpClient.client.Get err ", zap.String("response", string(bs)), zap.Error(err))
		return nil, nil, err
	}

	// 读取 response 的信息 并且记录到日志
	dump, err := httputil.DumpResponse(res, false) // better way
	if err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "client.doProdGet Response httputil.DumpResponse err ", zap.String("response", string(bs)), zap.Error(err))
		return res, bs, nil
	}
	logger.ApiWarn(c.LogFileName, c.RequestId, "client.doProdGet Response ", zap.String("response", string(bs)), zap.String("res", string(dump)))
	// xlog.Debugf("Wechat_Response: %d > %s", res.StatusCode, string(bs))
	// xlog.Debugf("Wechat_Headers: %#v", res.Header)
	// xlog.Debugf("Wechat_SignInfo: %#v", si)
	return res, bs, nil
}
