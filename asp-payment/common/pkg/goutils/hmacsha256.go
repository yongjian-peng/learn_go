/***************************************************
 ** @Desc : 获取一个hmac_sha256的字符串
 ** @Time : 2022/9/13 16:06
 ** @Author : yongjian
 ** @File : hmac_sha256
 * @Last Modified by: mikey.zhaopengpeng
 * @Last Modified time: 2022-09-22 20:02:20:20
 ** @Software: Vscode
****************************************************/
package goutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/spf13/cast"
	"hash"
)

// 获取支付正式环境Sign值
func GetReleaseSign(params map[string]interface{}, keys []string, paySecret string) (sign string) {

	str := ""
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		if len(cast.ToString(params[k])) == 0 {
			continue
		}
		str += k + "=" + cast.ToString(params[k]) + "&"
	}
	str += "paySecret=" + paySecret

	//fmt.Println("str----------", str)

	var h hash.Hash

	h = hmac.New(sha256.New, []byte(paySecret))
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))
	// fmt.Println("sha----------", sha)
	return base64.StdEncoding.EncodeToString([]byte(sha))
	// return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
