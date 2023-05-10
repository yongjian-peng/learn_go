/***************************************************
 ** @Desc : This file for ...
 ** @Time : 2019/10/26 11:08
 ** @Author : yuebin
 ** @File : sign_verify
 * @Last Modified by: mikey.zhaopeng
 * @Last Modified time: 2022-09-22 20:02:33
 ** @Software: GoLand
****************************************************/
package goutils

import (
	"github.com/spf13/cast"
)

func GetMD5Sign(params map[string]interface{}, keys []string, paySecret string) string {
	str := ""
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		if len(cast.ToString(params[k])) == 0 {
			continue
		}
		str += k + "=" + cast.ToString(params[k]) + "&"
	}
	str += "paySecret=" + paySecret
	sign := GetMD5Upper(str)
	return sign
}

/*
* 验签
 */
func Md5Verify(params map[string]interface{}, paySecret string) bool {
	sign := params["sign"]
	if sign == "" {
		return false
	}

	delete(params, "sign")
	keys := SortMap(params)
	tmpSign := GetMD5Sign(params, keys, paySecret)
	if tmpSign != sign {
		return false
	} else {
		return true
	}
}

/*
* HMAC_SHA256 加密
 */
func GetHmacSHA256(params map[string]interface{}, paySecret string) bool {
	sign := params["sign"]
	if sign == "" {
		return false
	}

	delete(params, "sign")
	keys := SortMap(params)
	tmpSign := GetReleaseSign(params, keys, paySecret)
	// fmt.Println("tmpSign--------------", tmpSign)
	if tmpSign != sign {
		return false
	} else {
		return true
	}
}

/*
* 验签 HMAC_SHA256
 */
func HmacSHA256Verify(params map[string]interface{}, paySecret string) bool {
	sign := params["sign"]
	if sign == "" {
		return false
	}

	delete(params, "sign")
	keys := SortMap(params)
	// fmt.Println("keys--------------", keys)
	// fmt.Println("paySecret--------------", paySecret)
	tmpSign := GetReleaseSign(params, keys, paySecret)
	// fmt.Println("tmpSign--------------", tmpSign)
	if tmpSign != sign {
		return false
	} else {
		return true
	}
}
