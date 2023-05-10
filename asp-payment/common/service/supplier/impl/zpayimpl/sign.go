package zpayimpl

import (
	"asp-payment/common/pkg/goutils"
	"crypto/md5"
	"fmt"
	"github.com/spf13/cast"
	"sort"
	"strings"
)

func (c *Client) signature(params map[string]interface{}, key string) (sign string) {

	delete(params, "sign")
	delete(params, "key")

	keys := goutils.SortMap(params)
	//fmt.Println("keys: ", keys)
	sign = GetMD5Sign(params, keys, key)
	//fmt.Println("sign: ", sign)
	return
}

func (c *Client) signature2(params map[string]interface{}, key string) (sign string, err error) {
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
	for k, v := range keys {
		// fmt.Println(k, v)
		// new_params[v] = params[v]
		if k > 0 {
			str = str + "&"
		}
		// fmt.Println(params[v])
		val, err := goutils.ToString(params[v])
		// fmt.Println("hhh", val, err)
		if err != nil {
			// fmt.Println("hhh", val, err)
			return "", err
		}
		str = str + v + "=" + val
	}
	str = str + "&key=" + key
	// fmt.Println("簽名字符串:", str, params)
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) // 将[]byte转成16进制
	sign = strings.ToUpper(md5str1)
	// sign = md5str1
	return
}

func GetMD5Sign(params map[string]interface{}, keys []string, paySecret string) string {
	str := ""
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		if len(cast.ToString(params[k])) == 0 {
			continue
		}
		str += k + "=" + cast.ToString(params[k]) + "&"
	}
	str += "key=" + paySecret
	//fmt.Println("str------------", str)
	sign := goutils.GetMD5Upper(str)
	//fmt.Println("sign: ", sign)
	return sign
}

func Md5Verify(params map[string]interface{}, paySecret string) bool {
	sign := params["sign"]
	if sign == "" {
		return false
	}

	delete(params, "sign")
	keys := goutils.SortMap(params)
	tmpSign := GetMD5Sign(params, keys, paySecret)
	if tmpSign != sign {
		return false
	} else {
		return true
	}
}

// VerifySignature 校验签名 signature
// secretKey
// bm
func VerifySignature(params map[string]interface{}, key string) (sign string) {

	delete(params, "Sign")
	delete(params, "sign")
	delete(params, "key")

	keys := goutils.SortMap(params)
	sign = GetMD5Sign(params, keys, key)

	return
}
