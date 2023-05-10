package seveneight

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/service/supplier/api"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Client struct {
	AppId        string
	SecretKey    string
	RequestId    string
	LogFileName  string
	Notify       string
	PayoutNotify string
}

func NewClient(channelDepartInfo *model.AspChannelDepartConfig, requestId, LogFileName string) (*Client, error) {
	var channelConfigInfo model.AspChannelDepartConfigInfo
	goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo)
	CallBack := config.AppConfig.Urls
	if CallBack.SevenEightH5NotifyUrl == "" || CallBack.SevenEightPayoutNotifyUrl == "" {
		logger.ApiWarn(constant.FirstPayLogFileName, requestId, "seveneight notify url err ", zap.Any("CallBack", CallBack))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		err := (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartPaymentCallBackConfigErrMsg)
		return nil, err
	}
	return &Client{
		AppId:        channelConfigInfo.Appid,
		SecretKey:    channelConfigInfo.Signature,
		RequestId:    requestId,
		Notify:       CallBack.SevenEightH5NotifyUrl,
		PayoutNotify: CallBack.SevenEightPayoutNotifyUrl,
		LogFileName:  LogFileName,
	}, nil
}

func (c *Client) GetHeader() map[string]interface{} {
	header := make(map[string]interface{})
	return header
}

func (c *Client) GetHttpServer() *api.HttpServer {
	return &api.HttpServer{
		BaseUrlProd: BaseUrlProd,
		RequestId:   c.RequestId,
		LogFileName: c.LogFileName,
	}
}

func GetMD5Sign(params map[string]interface{}, keys []string, paySecret string) string {
	str := ""
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		if len(cast.ToString(params[k])) == 0 {
			continue
		}
		str += k + "=" + cast.ToString(params[k]) + "&"
	}
	str += "key=" + paySecret
	//fmt.Println("str:", str)
	// fmt.Println("str------------", str)
	sign := goutils.GetMD5Upper(str)
	return sign
}

func GetSignature(params map[string]interface{}, key string) (sign string) {
	delete(params, "sign")
	delete(params, "key")
	keys := goutils.SortMap(params)

	//fmt.Println("keys: ", keys)
	sign = GetMD5Sign(params, keys, key)
	//fmt.Println("sign: ", sign)
	return
}

// CreateOrder
// H5下单API
// 返回的是 支付的跳转的 payment_link
//
//	StatusCode = 200 is success
func (c *Client) CreateOrder(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "seveneightPay ", "CreateOrder", time.Now(), float64(3))
	bm.Set("pay_md5sign", GetSignature(bm, c.SecretKey))
	bm.Set("pay_response", PayResponseJson)
	res, bs, err := c.GetHttpServer().PostForm(OrderCreateUrl, c.GetHeader(), bm)
	//fmt.Println("bs", string(bs))
	if err != nil {
		return nil, err
	}
	fpRsq := &PaymentRsp{StatusCode: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.CreateOrderBody = new(CreateOrderBody)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreateOrder ", zap.Error(err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		fpRsq.StatusCode = res.StatusCode
		fpRsq.StatusMsg = string(bs)
		return fpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &fpRsq.CreateOrderBody); err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "Body json err ", zap.Error(err))
		return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
	}

	return fpRsq, nil
}

func (c *Client) QueryOrder(bm model.BodyMap) (*PaymentRsp, error) {
	fpRsq := &PaymentRsp{StatusCode: Success}
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "seveneightPay ", "QueryOrder", time.Now(), float64(3))
	bm.Set("pay_md5sign", GetSignature(bm, c.SecretKey))
	res, bs, err := c.GetHttpServer().PostForm(QueryOrderUrl, c.GetHeader(), bm)
	//fmt.Println("bs", string(bs))
	if err != nil {
		return nil, err
	}
	fpRsq = &PaymentRsp{StatusCode: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.QueryOrderBody = new(QueryOrderBody)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreateOrder ", zap.Error(err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		fpRsq.StatusCode = res.StatusCode
		fpRsq.StatusMsg = string(bs)
		return fpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &fpRsq.QueryOrderBody); err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "Body json err ", zap.Error(err))
		return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
	}

	return fpRsq, nil
}

func (c *Client) QueryPayout(bm model.BodyMap) (*PaymentRsp, error) {
	fpRsq := &PaymentRsp{StatusCode: Success}
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "seveneightPay ", "QueryPayout", time.Now(), float64(3))
	bm.Set("pay_md5sign", GetSignature(bm, c.SecretKey))
	res, bs, err := c.GetHttpServer().PostForm(QueryPayoutUrl, c.GetHeader(), bm)
	if err != nil {
		return nil, err
	}
	fpRsq = &PaymentRsp{StatusCode: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.QueryPayoutBody = new(QueryPayoutBody)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreateOrder ", zap.Error(err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		fpRsq.StatusCode = res.StatusCode
		fpRsq.StatusMsg = string(bs)
		return fpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &fpRsq.QueryPayoutBody); err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "Body json err ", zap.Error(err))
		return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
	}

	return fpRsq, nil
}

// CreatePayout
// 代付
// 返回的是 支付的跳转的 payment_link
//
//	StatusCode = 200 is success
func (c *Client) CreatePayout(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "seveneightPay ", "CreatePayout", time.Now(), float64(3))
	bm.Set("pay_md5sign", GetSignature(bm, c.SecretKey))
	res, bs, err := c.GetHttpServer().PostForm(PayOutCreateUrl, c.GetHeader(), bm)
	fmt.Println("bs", string(bs))
	if err != nil {
		return nil, err
	}
	fpRsq := &PaymentRsp{StatusCode: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.CreatePayoutBody = new(CreatePayoutBody)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreateOrder ", zap.Error(err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		fpRsq.StatusCode = res.StatusCode
		fpRsq.StatusMsg = string(bs)
		return fpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &fpRsq.CreatePayoutBody); err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "Body json err ", zap.Error(err))
		return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
	}

	return fpRsq, nil
}
