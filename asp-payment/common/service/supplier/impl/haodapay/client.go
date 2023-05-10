package haodapay

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/service/supplier/api"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"hash"
	"net/http"
	"time"
)

type Client struct {
	AppId                string
	SecretKey            string
	PayoutClientSecret   string
	PayinChecksumSecret  string
	PayoutChecksumSecret string
	RequestId            string
	LogFileName          string
	Notify               string
	PayoutNotify         string
}

// {"x_client_id":"6HKGkNQWZy1265","x_client_secret":"Xy613U20tv230315035509"}
type ChannelConfigInfo struct {
	XClientId            string `json:"x_client_id"`
	XClientSecret        string `json:"x_client_secret"`
	PayinChecksumSecret  string `json:"payin_checksum_secret"`
	PayoutChecksumSecret string `json:"payout_checksum_secret"`
	PayoutClientSecret   string `json:"payout_client_secret"`
	//Extion               map[string]interface{} `json:"extion"`
}

func NewClient(channelDepartInfo *model.AspChannelDepartConfig, requestId, LogFileName string) (*Client, error) {
	var channelConfigInfo ChannelConfigInfo
	goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo)
	//CallBack := config.AppConfig.Urls
	//if CallBack.SevenEightH5NotifyUrl == "" || CallBack.SevenEightPayoutNotifyUrl == "" {
	//	logger.ApiWarn(constant.FirstPayLogFileName, requestId, "haodapayreq notify url err ", zap.Any("CallBack", CallBack))
	//	MissNotFoundErrCode := *appError.MissNotFoundErrCode
	//	err := (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartPaymentCallBackConfigErrMsg)
	//	return nil, err
	//}
	return &Client{
		AppId:                channelConfigInfo.XClientId,
		SecretKey:            channelConfigInfo.XClientSecret,
		PayoutClientSecret:   channelConfigInfo.PayoutClientSecret,
		PayinChecksumSecret:  channelConfigInfo.PayinChecksumSecret,
		PayoutChecksumSecret: channelConfigInfo.PayoutChecksumSecret,
		RequestId:            requestId,
		LogFileName:          LogFileName,
	}, nil
}

func (c *Client) GetHeader() map[string]interface{} {
	header := make(map[string]interface{})
	header["x-client-id"] = c.AppId
	header["x-client-secret"] = c.SecretKey
	return header
}

func (c *Client) GetHttpServer() *api.HttpServer {
	return &api.HttpServer{
		BaseUrlProd: BaseUrlProd,
		RequestId:   c.RequestId,
		LogFileName: c.LogFileName,
	}
}

func (c *Client) SetHttpServerProxyUrl(httpServer *api.HttpServer, proxyUrl string) *api.HttpServer {
	httpServer.ProxyUrl = proxyUrl
	return httpServer
}

func (c *Client) SetHttpServerBaseUrl(httpServer *api.HttpServer, baseUrl string) *api.HttpServer {
	httpServer.BaseUrlProd = baseUrl
	return httpServer
}

func GetSignature(params map[string]interface{}, key string) string {
	signatureStr := cast.ToString(params["order_id"]) + cast.ToString(params["reference"]) + cast.ToString(params["payer_UPIID"]) + key + cast.ToString(params["amount"]) + cast.ToString(params["UTR"])
	var h hash.Hash
	h = hmac.New(sha256.New, []byte(key))
	h.Write([]byte(signatureStr))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func GetPayoutSignature(params map[string]interface{}, key string) string {
	signatureStr := cast.ToString(params["payout_id"]) + cast.ToString(params["reference"]) + cast.ToString(params["beneficiary_account_num"]) + key + cast.ToString(params["beneficiary_account_ifsc"]) + cast.ToString(params["amount"]) + cast.ToString(params["UTR"])
	var h hash.Hash
	h = hmac.New(sha256.New, []byte(key))
	h.Write([]byte(signatureStr))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func GetPayoutUpiSignature(params map[string]interface{}, key string) string {
	signatureStr := cast.ToString(params["payout_id"]) + cast.ToString(params["reference"]) + cast.ToString(params["beneficiary_upi_handle"]) + key + cast.ToString(params["amount"]) + cast.ToString(params["UTR"])
	var h hash.Hash
	h = hmac.New(sha256.New, []byte(key))
	h.Write([]byte(signatureStr))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// CreateOrder
// H5下单API
// 返回的是 支付的跳转的 payment_link
//
//	StatusCode = 200 is success
func (c *Client) CreateOrder(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "haodaPay ", "CreateOrder", time.Now(), float64(3))

	httpServer := c.GetHttpServer()
	httpServer = c.SetHttpServerProxyUrl(httpServer, BasePayoutProxyUrl)

	res, bs, err := httpServer.PostJson(OrderCreateUrl, c.GetHeader(), bm)
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

	_, ok := fpRsq.CreateOrderBody.Data.(string)
	if !ok {
		jsonBody, errJson := json.Marshal(fpRsq.CreateOrderBody.Data)
		if errJson != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errJson", errJson))
			return nil, errJson
		}
		createOrderData := new(CreateOrderData)
		if errOrderData := goutils.JsonDecodeByte(jsonBody, &createOrderData); errOrderData != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errOrderData", errOrderData))
			return nil, errOrderData
		}
		fpRsq.StatusMsg = fpRsq.CreateOrderBody.Status
		fpRsq.CreateOrderBody.CreateOrderData = createOrderData
	}
	fpRsq.StatusMsg = fpRsq.CreateOrderBody.Message

	return fpRsq, nil
}

func (c *Client) QueryOrder(bm model.BodyMap) (*PaymentRsp, error) {
	fpRsq := &PaymentRsp{StatusCode: Success}
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "haodaPay ", "QueryOrder", time.Now(), float64(3))
	httpServer := c.GetHttpServer()
	//httpServer = c.SetHttpServerBaseUrl(httpServer, BasePayoutUrlProd)
	httpServer = c.SetHttpServerProxyUrl(httpServer, BasePayoutProxyUrl)
	//bm.Set("pay_md5sign", GetSignature(bm, c.SecretKey))
	res, bs, err := httpServer.PostJson(QueryOrderUrl, c.GetHeader(), bm)
	//fmt.Println("bs", string(bs))
	if err != nil {
		return nil, err
	}
	fpRsq = &PaymentRsp{StatusCode: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.QueryOrderBody = new(QueryOrderBody)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.QueryOrder ", zap.Error(err))
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

	_, ok := fpRsq.QueryOrderBody.Data.(string)
	if !ok {
		jsonBody, errJson := json.Marshal(fpRsq.QueryOrderBody.Data)
		if errJson != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errJson", errJson))
			return nil, errJson
		}
		queryOrderData := new(QueryOrderData)
		if errOrderData := goutils.JsonDecodeByte(jsonBody, &queryOrderData); errOrderData != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errOrderData", errOrderData))
			return nil, errOrderData
		}
		fpRsq.QueryOrderBody.QueryOrderData = queryOrderData
	}
	fpRsq.StatusMsg = fpRsq.QueryOrderBody.Status

	return fpRsq, nil
}

func (c *Client) QueryPayout(bm model.BodyMap) (*PaymentRsp, error) {
	fpRsq := &PaymentRsp{StatusCode: Success}
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "haodaPay ", "QueryPayout", time.Now(), float64(3))
	httpServer := c.GetHttpServer()
	httpServer = c.SetHttpServerBaseUrl(httpServer, BasePayoutUrlProd)
	httpServer = c.SetHttpServerProxyUrl(httpServer, BasePayoutProxyUrl)
	header := c.GetHeader()
	header["x-client-secret"] = c.PayoutClientSecret
	res, bs, err := httpServer.PostJson(QueryPayoutUrl, header, bm)
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

	_, ok := fpRsq.QueryPayoutBody.Data.(string)
	if !ok {
		jsonBody, errJson := json.Marshal(fpRsq.QueryPayoutBody.Data)
		if errJson != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errJson", errJson))
			return nil, errJson
		}
		queryPayoutData := new(QueryPayoutData)
		// 通过字符串转结构体
		if errPayoutData := goutils.JsonDecodeByte(jsonBody, &queryPayoutData); errPayoutData != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errPayoutData", errPayoutData))
			return nil, errPayoutData
		}
		fpRsq.QueryPayoutBody.QueryPayoutData = queryPayoutData
	}
	fpRsq.StatusMsg = fpRsq.QueryPayoutBody.Status

	return fpRsq, nil
}

// CreatePayout
// 代付
// 返回的是 支付的跳转的 payment_link
//
//	StatusCode = 200 is success
func (c *Client) CreatePayout(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "haodaPay ", "CreatePayout", time.Now(), float64(3))
	//bm.Set("pay_md5sign", GetSignature(bm, c.SecretKey))
	httpServer := c.GetHttpServer()
	httpServer = c.SetHttpServerBaseUrl(httpServer, BasePayoutUrlProd)
	httpServer = c.SetHttpServerProxyUrl(httpServer, BasePayoutProxyUrl)
	header := c.GetHeader()
	header["x-client-secret"] = c.PayoutClientSecret
	res, bs, err := httpServer.PostJson(PayOutCreateUrl, header, bm)
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
	fpRsq.StatusMsg = fpRsq.CreatePayoutBody.Message

	return fpRsq, nil
}

// CreatePayoutUpi
// 代付
// 返回的是 支付的跳转的 payment_link
//
//	StatusCode = 200 is success
func (c *Client) CreatePayoutUpi(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "haodaPay ", "CreatePayoutUpi", time.Now(), float64(3))
	//bm.Set("pay_md5sign", GetSignature(bm, c.SecretKey))
	httpServer := c.GetHttpServer()
	httpServer = c.SetHttpServerBaseUrl(httpServer, BasePayoutUrlProd)
	httpServer = c.SetHttpServerProxyUrl(httpServer, BasePayoutProxyUrl)
	header := c.GetHeader()
	header["x-client-secret"] = c.PayoutClientSecret
	res, bs, err := httpServer.PostJson(PayOutCreateUpiUrl, header, bm)
	fmt.Println("bs", string(bs))
	if err != nil {
		return nil, err
	}
	fpRsq := &PaymentRsp{StatusCode: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.CreatePayoutUpiBody = new(CreatePayoutUpiBody)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreatePayoutUpiBody ", zap.Error(err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		fpRsq.StatusCode = res.StatusCode
		fpRsq.StatusMsg = res.Status
		return fpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &fpRsq.CreatePayoutUpiBody); err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "Body json err ", zap.Error(err))
		return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
	}
	fpRsq.StatusMsg = fpRsq.CreatePayoutUpiBody.Message

	return fpRsq, nil
}

// UpiValidate 验证upi账户是否合法
//
//	StatusCode = 200 is success
func (c *Client) UpiValidate(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "haodaPay ", "UpiValidate", time.Now(), float64(3))

	httpServer := c.GetHttpServer()
	httpServer = c.SetHttpServerBaseUrl(httpServer, BasePayoutUrlProd)
	httpServer = c.SetHttpServerProxyUrl(httpServer, BasePayoutProxyUrl)

	res, bs, err := httpServer.PostJson(UpiValidate, c.GetHeader(), bm)
	//fmt.Println("bs", string(bs))
	if err != nil {
		return nil, err
	}
	fpRsq := &PaymentRsp{StatusCode: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.UpiValidateBody = new(UpiValidateBody)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.UpiValidate ", zap.Error(err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		fpRsq.StatusCode = res.StatusCode
		fpRsq.StatusMsg = res.Status
		return fpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &fpRsq.UpiValidateBody); err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "Body json err ", zap.Error(err))
		return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
	}
	fpRsq.StatusMsg = fpRsq.UpiValidateBody.Message

	return fpRsq, nil
}
