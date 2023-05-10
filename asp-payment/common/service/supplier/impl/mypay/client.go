package mypay

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/pkg/xhttp"
	"asp-payment/common/service/supplier/api"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"strings"
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
	return &Client{
		AppId:       channelConfigInfo.Appid,
		SecretKey:   channelConfigInfo.Signature,
		RequestId:   requestId,
		LogFileName: LogFileName,
	}, nil
}

func (c *Client) GetHeader() map[string]interface{} {
	header := make(map[string]interface{})
	header["Content-Type"] = "application/json"
	header["Accept"] = "*/*"
	return header
}

func (c *Client) GetHttpServer() *api.HttpServer {
	return &api.HttpServer{
		BaseUrlProd: BaseUrlProd,
		RequestId:   c.RequestId,
		LogFileName: c.LogFileName,
	}
}

func GetHash256Sign(params map[string]interface{}, keys []string, paySecret string) string {
	str := ""
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		str += k + "=" + cast.ToString(params[k]) + "~"
	}
	//删除最后一个~
	str = str[:len(str)-1]
	str += paySecret
	fmt.Println(str)
	m := sha256.New()
	m.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

func GetSignature(params map[string]interface{}, key string) (sign string) {
	delete(params, "HASH")
	keys := goutils.SortMap(params)

	//fmt.Println("keys: ", keys)
	sign = GetHash256Sign(params, keys, key)
	//fmt.Println("sign: ", sign)
	return
}

// CreateOrder
// 代付
// 返回的是 支付的跳转的 payment_link
//
//	StatusCode = 200 is success
func (c *Client) CreateOrder(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "myPay ", "CreateOrder", time.Now(), float64(3))
	//bm.Set("pay_md5sign", GetSignature(bm, c.SecretKey))

	res, bs, err := c.GetHttpServer().PostJson(OrderCreateUrl, c.GetHeader(), bm)
	fmt.Println("bs", string(bs))
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
		fpRsq.StatusMsg = res.Status
		return fpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &fpRsq.CreateOrderBody); err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "Body json err ", zap.Error(err))
		return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
	}

	//_, ok := fpRsq.CreateOrderBody.Data.(string)
	//if !ok {
	//	// 转字符串
	//	jsonBody, errJson := json.Marshal(fpRsq.CreateOrderBody.Data)
	//	if errJson != nil {
	//		logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errJson", errJson))
	//		return nil, errJson
	//	}
	//	createOrderData := new(CreateOrderData)
	//	if errOrderData := goutils.JsonDecodeByte(jsonBody, &createOrderData); errOrderData != nil {
	//		logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errOrderData", errOrderData))
	//		return nil, errOrderData
	//	}
	//	fpRsq.StatusMsg = fpRsq.CreateOrderBody.Status
	//	fpRsq.CreateOrderBody.CreateOrderData = createOrderData
	//} else {
	//	fpRsq.StatusMsg = fpRsq.CreateOrderBody.Status
	//}
	fpRsq.StatusMsg = fpRsq.CreateOrderBody.Message
	return fpRsq, nil
}

func (c *Client) QueryOrder(bm model.BodyMap) (*PaymentRsp, error) {
	fpRsq := &PaymentRsp{StatusCode: Success}
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "mypayreq ", "QueryOrder", time.Now(), float64(3))
	param := bm.EncodeURLParams()
	queryOrderUrl := QueryOrderUrl + "?" + param
	res, bs, err := c.GetHttpServer().Get(queryOrderUrl, xhttp.TypeJSON, c.GetHeader(), bm)
	fmt.Println("bs", string(bs))
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
		fpRsq.StatusMsg = res.Status
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
		fpRsq.StatusMsg = fpRsq.QueryOrderBody.Status
		fpRsq.QueryOrderBody.QueryOrderData = queryOrderData
	} else {
		fpRsq.StatusMsg = fpRsq.QueryPayoutBody.Status
	}

	return fpRsq, nil
}

func (c *Client) QueryPayout(bm model.BodyMap) (*PaymentRsp, error) {
	fpRsq := &PaymentRsp{StatusCode: Success}
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "myPay ", "CreateOrder", time.Now(), float64(3))
	param := bm.EncodeURLParams()
	queryPayoutUrl := QueryPayoutUrl + "?" + param
	res, bs, err := c.GetHttpServer().Get(queryPayoutUrl, xhttp.TypeJSON, c.GetHeader(), bm)
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
		fpRsq.StatusMsg = res.Status
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
		fpRsq.StatusMsg = fpRsq.QueryPayoutBody.StatusCode
		fpRsq.QueryPayoutBody.QueryPayoutData = queryPayoutData
	}
	fpRsq.StatusMsg = fpRsq.QueryPayoutBody.StatusCode
	return fpRsq, nil
}

// CreatePayout
// 代付
// 返回的是 支付的跳转的 payment_link
//
//	StatusCode = 200 is success
func (c *Client) CreatePayout(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "myPay ", "CreatePayout", time.Now(), float64(3))
	res, bs, err := c.GetHttpServer().PostJson(PayOutCreateUrl, c.GetHeader(), bm)
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
		fpRsq.StatusMsg = res.Status
		return fpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &fpRsq.CreatePayoutBody); err != nil {
		logger.ApiWarn(c.LogFileName, c.RequestId, "Body json err ", zap.Error(err))
		return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
	}

	_, ok := fpRsq.CreatePayoutBody.Data.(string)
	if !ok {
		jsonBody, errJson := json.Marshal(fpRsq.CreatePayoutBody.Data)
		if errJson != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errJson", errJson))
			return nil, errJson
		}
		createPayoutData := new(CreatePayoutData)
		if errPayoutData := goutils.JsonDecodeByte(jsonBody, &createPayoutData); errPayoutData != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "Response json err ", zap.Any("errPayoutData", errPayoutData))
			return nil, errPayoutData
		}
		fpRsq.CreatePayoutBody.CreatePayoutData = createPayoutData
	}
	fpRsq.StatusMsg = fpRsq.CreatePayoutBody.Status

	return fpRsq, nil
}
