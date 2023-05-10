package zpayimpl

import (
	"asp-payment/common/model"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"

	"hash"
	"strings"
)

// 获取QQ支付正式环境Sign值
func GetReleaseSign(apiKey string, signType string, bm model.BodyMap) (sign string) {
	var h hash.Hash
	if signType == SignType_HMAC_SHA256 {
		h = hmac.New(sha256.New, []byte(apiKey))
	} else {
		h = md5.New()
	}
	h.Write([]byte(bm.EncodeWeChatSignParams(apiKey)))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
