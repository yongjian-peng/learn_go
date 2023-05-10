package firstpayimpl

import (
	"asp-payment/common/service/supplier/api"
	"sync"
)

type Client struct {
	AppId       string
	SerectKey   string
	mu          sync.RWMutex
	RequestId   string
	LogFileName string
}

// NewClient 初始化First Pay
// appId：应用ID
// mchId：商户ID
// ApiKey：API秘钥值
// IsProd：是否是正式环境
func NewClient(appId, SerectKey string) *Client {
	return &Client{
		AppId:     appId,
		SerectKey: SerectKey,
	}
}

func (c *Client) GetHeader(str string) map[string]interface{} {
	header := make(map[string]interface{})
	header[HeaderSignature] = c.signature(c.SerectKey, str)
	header[HeaderAppId] = c.AppId
	header["Content-Type"] = "application/json; charset=utf-8"
	return header
}

func (c *Client) GetHttpServer() *api.HttpServer {
	return &api.HttpServer{
		BaseUrlProd: baseUrlProd,
		RequestId:   c.RequestId,
		LogFileName: c.LogFileName,
	}
}
