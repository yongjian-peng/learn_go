package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/spf13/cast"
)

/*
* 获取小写的MD5
 */
func GetMD5LOWER(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

/*
* 对map的key值进行排序
 */
func SortMap(m map[string]interface{}) []string {
	var arr []string
	for k := range m {
		arr = append(arr, cast.ToString(k))
	}
	sort.Strings(arr)
	return arr
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
	// fmt.Println("str------------", str)
	sign := GetMD5LOWER(str)
	return sign
}

func GetSignature(params map[string]interface{}, key string) (sign string) {
	delete(params, "sign")
	delete(params, "key")
	keys := SortMap(params)

	//fmt.Println("keys: ", keys)
	sign = GetMD5Sign(params, keys, key)
	//fmt.Println("sign: ", sign)
	return
}

func main() {

	/*
		{
		    "accountname":"ericluzhonghua",
		    "bankname":"ICIC",
		    "cardnumber":"006663300006572",
		    "email":"ericluzhonghua@gmail.com",
		    "extdata":"",
		    "ifsc":"ICIC0006631",
		    "mchid":"202301104",
		    "mobile":"9036830689",
		    "money":100,
		    "notifyurl":"https://hopepatti.fun/callback/smartpaytime/payout",
		    "out_trade_no":"122023020816242300000001",
		    "pay_md5sign":"d66f10e3995dceaebd86eab355b2d5df",
		    "paymentmode":"IMPS"
		}
	*/
	bm := make(map[string]interface{})
	bm["accountname"] = "ericluzhonghua"
	bm["bankname"] = "ICIC"
	bm["cardnumber"] = "006663300006572"
	bm["email"] = "ericluzhonghua@gmail.com"
	bm["extdata"] = ""
	bm["ifsc"] = "ICIC0006631"
	bm["mchid"] = "202301104"
	bm["mobile"] = "9036830689"
	bm["money"] = 100
	bm["notifyurl"] = "https://hopepatti.fun/callback/smartpaytime/payout"
	bm["out_trade_no"] = "122023020816242300000001"
	bm["paymentmode"] = "IMPS"

	SecretKey := "soixBAwgZamPvJYEQeLMDCHkSpcVOUbh"
	sign := GetSignature(bm, SecretKey)

	fmt.Println(sign)

	sign2 := GetMD5LOWER("Shanghai")
	fmt.Println(sign2)

	sign3 := GetMD5LOWER("accountname=ericluzhonghua&bankname=ICIC&cardnumber=006663300006572&email=ericluzhonghua@gmail.com&ifsc=ICIC0006631&mchid=202301104&mobile=9036830689&money=100&notifyurl=https://hopepatti.fun/callback/smartpaytime/payout&out_trade_no=122023020816404500000001&paymentmode=IMPS&key=soixBAwgZamPvJYEQeLMDCHkSpcVOUbh")
	fmt.Println(sign3)

}
