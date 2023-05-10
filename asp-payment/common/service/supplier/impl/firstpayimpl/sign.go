package firstpayimpl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash"
)

// 生成签名 signature
// func (c *Client) signature(serectKey string, bm commonpay.BodyMap) (string, error) {
// 	var h hash.Hash

// 	str := bm.JsonBody()

// 	h = hmac.New(sha256.New, []byte(serectKey))
// 	h.Write([]byte(str))
// 	sha := hex.EncodeToString(h.Sum(nil))
// 	return base64.StdEncoding.EncodeToString([]byte(sha)), nil
// }

func (c *Client) signature(serectKey string, str string) string {
	var h hash.Hash
	h = hmac.New(sha256.New, []byte(serectKey))
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

// VerifySignature 校验签名 signature
// secretKey
// bm
func VerifySignature(secretKey string, str string) string {
	var h hash.Hash

	// str := bm.JsonBody()

	h = hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}
