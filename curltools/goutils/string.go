package goutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func SearchString(slice []string, s string) int {
	for i, v := range slice {
		if s == v {
			return i
		}
	}
	return -1
}

func SafeHtml(s string) string {
	r := strings.NewReplacer("<input", "&lt;input", "<a ", "&lt; a")
	return r.Replace(s)
}

func substr(s string, pos int, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func JoinString(strings ...string) string {
	var buf bytes.Buffer
	for _, str := range strings {
		buf.WriteString(str)
	}
	return buf.String()
}

func Truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

/*
IsBlank checks if a string is whitespace or empty (""). Observe the following behavior:

	goutils.IsBlank("")        = true
	goutils.IsBlank(" ")       = true
	goutils.IsBlank("bob")     = false
	goutils.IsBlank("  bob  ") = false

Parameter:

	str - the string to check

Returns:

	true - if the string is whitespace or empty ("")
*/
func IsBlank(str string) bool {
	strLen := len(str)
	if str == "" || strLen == 0 {
		return true
	}
	for i := 0; i < strLen; i++ {
		if unicode.IsSpace(rune(str[i])) == false {
			return false
		}
	}
	return true
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

// IsEmpty checks if a string is empty (""). Returns true if empty, and false otherwise.
func IsEmpty(str string) bool {
	return len(str) == 0
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// Substr 截取字符串
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

func Equals(a, b string) bool {
	return a == b
}

func EqualsIgnoreCase(a, b string) bool {
	return a == b || strings.ToUpper(a) == strings.ToUpper(b)
}

// RuneLen 字符成长度
func RuneLen(s string) int {
	bt := []rune(s)
	return len(bt)
}

// GetSummary 获取summary
func GetSummary(s string, length int) string {
	s = strings.TrimSpace(s)
	summary := Substr(s, 0, length)
	if RuneLen(s) > length {
		summary += "..."
	}
	return summary
}

/**
 * 字符串首字母转化为大写 ios -> Ios
 */
func strFirstToUpper(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			vv[i] -= 32
			upperStr += string(vv[i]) // + string(vv[i+1])
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// FirstCharUpper ios_bbbbbbbb -> IosBbbbbbbbb
func FirstCharUpper(str string) string {
	showName := ""
	if strings.Contains(str, "_") {
		names := strings.Split(str, "_")
		for _, nameItem := range names {
			showName += strFirstToUpper(nameItem)
		}
	} else {
		showName = strFirstToUpper(str)
	}
	return showName
}

// SecondCharUpper ios_bbbbbbbb -> iosBbbbbbbbb
func SecondCharUpper(str string) string {
	showName := ""
	if strings.Contains(str, "_") {
		names := strings.Split(str, "_")
		for index, nameItem := range names {
			if index == 0 {
				showName += nameItem
			} else {
				showName += strFirstToUpper(nameItem)
			}
		}
	} else {
		showName = str
	}
	return showName
}

/*
*
检查是否是空空字符
*/
func CheckIsEmptyString(str string) (string, bool) {
	str = strings.NewReplacer(" ", "", "　", "", "\t", "", "\n", "", "\r", "").Replace(str)
	if str == "" {
		return "", true
	}
	return str, false
}

// 反转字符串
func ReverseSliceString(l []string) {
	for i := 0; i < len(l)/2; i++ {
		li := len(l) - i - 1
		l[i], l[li] = l[li], l[i]
	}
}

func TrimSpace(str string) string {
	nick := strings.NewReplacer(" ", "", "　", "", "\t", "", "\n", "", "\r", "").Replace(str)
	return nick
}

func CreateToken(identify string) string {
	return Md5(JoinString(identify, strconv.FormatInt(time.Now().UnixNano(), 10)))
}

func ConvertToString(v interface{}) (str string) {
	if v == nil {
		return NULL
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = json.Marshal(v); err != nil {
		return NULL
	}
	str = string(bs)
	return
}

/**
 * 工具类：
 * 转成字符串
 */
func ToString(val interface{}) (str string, err error) {
	// fmt.Printf("%T", val)
	var s string
	if vv, ok := val.(float64); ok {
		s = strconv.Itoa(int(vv))
	} else if vv, ok := val.(float32); ok {
		s = strconv.Itoa(int(vv))
	} else if vv, ok := val.(int); ok {
		s = strconv.Itoa(int(vv))
	} else if vv, ok := val.(int64); ok {
		s = strconv.Itoa(int(vv))
	} else if vv, ok := val.(string); ok {
		s = string(vv)
	} else {
		return s, errors.New("不支持的参数类型")
	}
	return s, nil
}

/**
** MapInterface 转 MapString
 */
func MapInterfaceToMapString(params map[string]interface{}) (map[string]string, error) {
	tmp := make(map[string]string)
	// 转字符串
	for k, v := range params {
		//val, err := ToString(v)
		//if err != nil {
		//	return tmp, err
		//}
		val := cast.ToString(v)
		tmp[k] = val
	}

	return tmp, nil
}

/**
** MapInterface 转 MapString
 */
func MapInterfaceToMapStringAndFistLower(params map[string]interface{}) (map[string]string, error) {
	tmp := make(map[string]string)
	// 转字符串
	for k, v := range params {
		val, err := ToString(v)
		if err != nil {
			return tmp, err
		}
		k = MakeFirstLowerCase(k)
		tmp[k] = val
	}

	return tmp, nil
}

func MakeFirstLowerCase(s string) string {

	if len(s) < 2 {
		return strings.ToLower(s)
	}

	bts := []byte(s)

	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}
