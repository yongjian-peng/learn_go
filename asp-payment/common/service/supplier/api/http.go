package api

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/logger"
	"asp-payment/common/pkg/xhttp"
	"context"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HttpServer struct {
	BaseUrlProd string
	RequestId   string
	LogFileName string
	ProxyUrl    string
	DebugSwitch bool
}

var (
	ctx = context.Background()
)

func (s *HttpServer) Post(uri string, typeStr xhttp.RequestType, headers map[string]interface{}, bm model.BodyMap) (res *http.Response, bs []byte, err error) {
	var prodUrl = s.BaseUrlProd + uri
	httpClient := xhttp.NewClient()
	for key, headerValue := range headers {
		httpClient.Header.Add(key, cast.ToString(headerValue))
	}

	logger.ApiWarn(s.LogFileName, s.RequestId, "client.doProdPost Request ",
		zap.Any("headers", headers),
		zap.String("request", bm.JsonBody()),
		zap.String("prodUrl", prodUrl))
	// 设置代理转发url
	if s.ProxyUrl != "" {
		proxyAddress, _ := url.Parse(s.ProxyUrl)
		httpClient = httpClient.SetProxyUrl(proxyAddress)
	}
	res, bs, err = httpClient.Type(typeStr).Post(prodUrl).SendBodyMap(bm).EndBytes(ctx)
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "httpClient.client.Post err ", zap.String("response", string(bs)), zap.Error(err))
		return nil, nil, err
	}
	// 读取 response 的信息 并且记录到日志
	dump, dErr := httputil.DumpResponse(res, false) // better way
	if dErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "client.doProdPost Body httputil.DumpResponse err ", zap.String("response", string(bs)), zap.Error(dErr))
		return res, bs, nil
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "client.doProdPost Body ", zap.String("response", string(bs)), zap.String("res", string(dump)))
	return
}

func (s *HttpServer) PostJson(uri string, headers map[string]interface{}, bm model.BodyMap) (res *http.Response, bs []byte, err error) {
	return s.Post(uri, xhttp.TypeJSON, headers, bm)
}

func (s *HttpServer) PostForm(uri string, headers map[string]interface{}, bm model.BodyMap) (res *http.Response, bs []byte, err error) {
	return s.Post(uri, xhttp.TypeForm, headers, bm)
}

func (s *HttpServer) Get(uri string, typeStr xhttp.RequestType, headers map[string]interface{}, bm model.BodyMap) (res *http.Response, bs []byte, err error) {
	var url = s.BaseUrlProd + uri
	//fmt.Println("url:",url)
	httpClient := xhttp.NewClient()
	for key, headerValue := range headers {
		httpClient.Header.Add(key, cast.ToString(headerValue))
	}
	logger.ApiWarn(s.LogFileName, s.RequestId, "client.doProdGet Request ",
		zap.Any("headers", headers),
		zap.String("request", bm.JsonBody()),
		zap.String("url", url))
	res, bs, err = httpClient.Type(typeStr).Get(url).SendBodyMap(bm).EndBytes(ctx)
	if err != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "client.doProdGet err ", zap.String("response", string(bs)), zap.Error(err))
		return nil, nil, err
	}

	// 读取 response 的信息 并且记录到日志
	dump, dErr := httputil.DumpResponse(res, false) // better way
	if dErr != nil {
		logger.ApiWarn(s.LogFileName, s.RequestId, "client.doProdGet Response", zap.String("response", string(bs)), zap.Error(dErr))
		return res, bs, nil
	}

	logger.ApiWarn(s.LogFileName, s.RequestId, "client.doProdGet Response ", zap.String("response", string(bs)), zap.String("res", string(dump)))

	return res, bs, nil
}
