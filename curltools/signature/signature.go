package signature

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

func Signature(params map[string]string, signal string) (err error, signature string) {

	keys := make([]string, len(params))
	i := 0

	for k, _ := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	str := ""

	for _, k := range keys {
		str += k + "=" + params[k] + "&"
	}

	str += "key=" + signal
	strdata := []byte(str)
	has := md5.Sum(strdata)
	signature = fmt.Sprintf("%x", has)

	signature = strings.ToUpper(signature)

	return nil, signature
}
