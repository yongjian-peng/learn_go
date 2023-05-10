package fynzonpay

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/aes"
	"asp-payment/common/pkg/appError"
	"asp-payment/common/pkg/config"
	"asp-payment/common/pkg/constant"
	"asp-payment/common/pkg/goutils"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/pkg/xhttp"
	"asp-payment/common/service/supplier/api"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	AppId             string
	SecretKey         string
	PartnerId         string
	ApplicationId     string
	PayoutSignature   string
	RequestId         string
	LogFileName       string
	Notify            string
	BeneficiaryNotify string
	PayoutNotify      string
}

func NewClient(channelDepartInfo *model.AspChannelDepartConfig, requestId, LogFileName string) (*Client, error) {
	// {"appid":"1515","partnerId":"53bb71c1edeea8e1258ae678472d2f6feb1dc8283a03f0e7f6efd019f398c17e","applicationId":"anlZYlVJekdPdlFSVDJkeVdma0hlNEQ0VjhrZk9WZk5WVWk2V2VmcFQwOHBHWFdsQm9xeWJBeTRDZS96WVFUYQ","signature":"MTEyOTNfMTUxNV8yMDIzMDExMzAwMTMzOA","payoutSignature":"Qc_9Z2pinBs@#Am"}
	var channelConfigInfo model.AspChannelDepartConfigInfo
	goutils.JsonDecode(channelDepartInfo.Config, &channelConfigInfo)
	CallBack := config.AppConfig.Urls
	if CallBack.FynzonpayWappayNotifyUrl == "" || CallBack.FynzonpayWappayErrorUrl == "" || CallBack.FynzonpayPayoutNotifyUrl == "" || CallBack.FynzonpayBeneficiaryNotifyUrl == "" || CallBack.FynzonpayOrderReturnUrl == "" {
		logger.ApiWarn(constant.FynzonpayLogFileName, requestId, "fynzonpay notify url err ", zap.Any("CallBack", CallBack))
		MissNotFoundErrCode := *appError.MissNotFoundErrCode
		err := (&MissNotFoundErrCode).FormatMessage(constant.MissChannelDepartPaymentCallBackConfigErrMsg)
		return nil, err
	}
	return &Client{
		AppId:             channelConfigInfo.Appid,
		SecretKey:         channelConfigInfo.Signature,
		PartnerId:         channelConfigInfo.PartnerId,
		ApplicationId:     channelConfigInfo.ApplicationId,
		PayoutSignature:   channelConfigInfo.PayoutSignature,
		RequestId:         requestId,
		Notify:            CallBack.FynzonpayWappayNotifyUrl,
		BeneficiaryNotify: CallBack.FynzonpayBeneficiaryNotifyUrl,
		PayoutNotify:      CallBack.FynzonpayPayoutNotifyUrl,
		LogFileName:       LogFileName,
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

// signature 加密
func (c *Client) signature(payToken string, str string) string {
	//encryptMethod := "AES-256-CBC"
	//var h hash.Hash
	//h = hmac.New(sha256.New, []byte(payToken))
	//h.Write([]byte(str))
	//sha := hex.EncodeToString(h.Sum(nil))
	h := sha256.New()
	h.Write([]byte(payToken))

	sha := hex.EncodeToString(h.Sum(nil))

	iv := sha[0:16]
	//fmt.Println("sha: ", sha)
	//fmt.Println("iv: ", iv)
	//fmt.Println("payToken: ", payToken)
	secretKey := c.PartnerId[0:32]
	encryptData, err := aes.CBCEncrypt([]byte(str), []byte(secretKey), []byte(iv))
	//fmt.Println("encryptData: ", encryptData)
	//fmt.Println("encryptData-string--------: ", string(encryptData))
	//fmt.Println("err: ", err)
	if err != nil {
		return ""
	}
	//  $output = rtrim( strtr( base64_encode( openssl_encrypt( $string, $encrypt_method, $secret_key, 0, $iv ) ), '+/', '-_'), '=');

	encryptText := base64.StdEncoding.EncodeToString(encryptData)
	//fmt.Println("encryptText: ", encryptText)

	encryptText2 := base64.StdEncoding.EncodeToString([]byte(encryptText))
	//fmt.Println("encryptText2: ", encryptText2)
	//全部替换
	output := strings.Replace(encryptText2, "+/", "-_", -1)
	//fmt.Println(strings.Replace(encryptText, "+/", "-_", -1))

	return strings.TrimRight(output, "=")
}

// CreateOrder
// H5下单API
// 返回的是 支付的跳转的 payment_link
//
//	StatusCode = 200 is success
func (c *Client) CreateOrder(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "fynzonpay ", "CreateOrder", time.Now(), float64(3))
	param := bm.EncodeURLParams()
	url := OrderCreateUrl + "?" + param
	logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreateOrder ", zap.Any("param:", param), zap.Any("bm:", bm))
	res, bs, err := c.GetHttpServer().Get(url, xhttp.TypeForm, c.GetHeader(), bm)
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
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "fynzonpay ", "QueryOrder", time.Now(), float64(3))

	param := bm.EncodeURLParams()
	logger.ApiWarn(c.LogFileName, c.RequestId, "client.QueryOrder ", zap.Any("param:", param), zap.Any("bm:", bm))
	url := QueryOrderUrl + "?" + param
	res, bs, err := c.GetHttpServer().Get(url, xhttp.TypeForm, c.GetHeader(), bm)
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
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "fynzonpayPay ", "QueryPayout", time.Now(), float64(3))
	param := bm.EncodeURLParams()
	logger.ApiWarn(c.LogFileName, c.RequestId, "client.QueryPayout ", zap.Any("param:", param), zap.Any("bm:", bm))
	encode := c.signature(c.ApplicationId, param)
	logger.ApiWarn(c.LogFileName, c.RequestId, "params ", zap.Any("param", param))
	bq := make(model.BodyMap)
	encode = encode + c.ApplicationId
	bq.Set("pram_encode", encode)
	res, bs, err := c.GetHttpServer().PostForm(QueryPayoutUrl, c.GetHeader(), bq)
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

	return fpRsq, nil
}

// AddbeneFiciary
// 代付
// 添加受益人
//
//	StatusCode = 200 is success
func (c *Client) AddbeneFiciary(bm model.BodyMap) (*PaymentRsp, error) {
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "fynzonpayPay ", "AddbeneFiciary", time.Now(), float64(3))

	param := bm.EncodeURLParams()
	logger.ApiWarn(c.LogFileName, c.RequestId, "client.AddbeneFiciary ", zap.Any("param:", param), zap.Any("bm:", bm))
	//fmt.Println("param: ", param)
	encode := c.signature(c.ApplicationId, param)
	//fmt.Println("encode: ", encode)
	bq := make(model.BodyMap)
	encode = encode + c.ApplicationId
	bq.Set("pram_encode", encode)
	//return nil, fmt.Errorf("error")
	//bm.Set("pay_md5sign", GetSignature(bm, c.SecretKey))
	res, bs, err := c.GetHttpServer().PostForm(PayoutBeneficiaryCreateUrl, c.GetHeader(), bq)
	fmt.Println("bs", string(bs))
	if err != nil {
		return nil, err
	}
	fpRsq := &PaymentRsp{StatusCode: Success}
	fpRsq.ErrorResponse = new(ErrorResponse)
	fpRsq.CreateBeneficiaryBody = new(CreateBeneficiaryBody)
	if res.StatusCode != http.StatusOK {
		if err = goutils.JsonDecodeByte(bs, &fpRsq.ErrorResponse); err != nil {
			logger.ApiWarn(c.LogFileName, c.RequestId, "client.AddbeneFiciary ", zap.Error(err))
			return nil, fmt.Errorf("[%s]: %v, bytes: %s", "unmarshal error", err, string(bs))
		}
		fpRsq.StatusCode = res.StatusCode
		fpRsq.StatusMsg = res.Status
		return fpRsq, nil
	}

	if err = goutils.JsonDecodeByte(bs, &fpRsq.CreateBeneficiaryBody); err != nil {
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
	defer goutils.ExecutionTime(c.LogFileName, c.RequestId, "fynzonpayPay ", "CreatePayout", time.Now(), float64(3))
	param := bm.EncodeURLParams()
	logger.ApiWarn(c.LogFileName, c.RequestId, "client.CreatePayout ", zap.Any("param:", param), zap.Any("bm:", bm))
	encode := c.signature(c.ApplicationId, param)

	bq := make(model.BodyMap)
	encode = encode + c.ApplicationId
	bq.Set("pram_encode", encode)
	res, bs, err := c.GetHttpServer().PostForm(PayOutCreateUrl, c.GetHeader(), bq)
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

	return fpRsq, nil
}
