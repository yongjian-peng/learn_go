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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"sort"
)

func GetMD5Sign(params map[string]string, keys []string, paySecret string) string {
	str := ""
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		if len(params[k]) == 0 {
			continue
		}
		str += k + "=" + params[k] + "&"
	}
	str += "paySecret=" + paySecret
	sign := GetMD5Upper(str)
	return sign
}

/*
* 验签
 */
func Md5Verify(params map[string]string, paySecret string) bool {
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
func GetHmacSHA256(params map[string]string, paySecret string) bool {
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
func HmacSHA256Verify(params map[string]string, paySecret string) bool {
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

func GetSign(params map[string]interface{}, key string) (sign string, err error) {
	// 删除sign字段(不参与签名)
	delete(params, "sign")
	delete(params, "key")
	// 排序
	keys := make([]string, len(params))
	i := 0
	// 注意：map的遍历是乱序的，也就是循环得到随机key
	for k, _ := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys) //排序字符串
	// var new_params = make(map[string]interface{})
	var str string
	for _, v := range keys {
		val, err := ToString(params[v])
		// fmt.Println("hhh", val, err)
		if err != nil {
			// fmt.Println("hhh", val, err)
			return "", err
		}
		if len(val) == 0 {
			continue
		}
		str = str + v + "=" + val + "&"
	}
	str = str + "paySecret=" + key

	// data := []byte(str)
	// has := md5.Sum(data)
	// md5str1 := fmt.Sprintf("%x", has) // 将[]byte转成16进制
	// sign = strings.ToUpper(md5str1)

	var h hash.Hash
	h = hmac.New(sha256.New, []byte(key))
	h.Write([]byte(str))
	sha := hex.EncodeToString(h.Sum(nil))
	// fmt.Println("sha----------", sha)
	// fmt.Println("簽名字符串:", str, params, sha)
	sign = base64.StdEncoding.EncodeToString([]byte(sha))
	// fmt.Println("簽名字符串:", str, params, sha, sign)
	return sign, nil
}
