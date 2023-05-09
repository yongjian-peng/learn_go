package check

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strconv"
	"unicode/utf8"
)

// 注册自定义tag校验
func RegisterAll(valid *validator.Validate) {
	_ = valid.RegisterValidation("ulen", ulen)
	_ = valid.RegisterValidation("ugt", ugt)
	_ = valid.RegisterValidation("ult", ult)
	_ = valid.RegisterValidation("omitemptyurl", omitemptyurl)
}

// 汉字长度校验
func ulen(fl validator.FieldLevel) bool {
	length := utf8.RuneCountInString(fl.Field().String())
	param, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}
	if length == param {
		return true
	}
	return false
}

// 汉字长度大于
func ugt(fl validator.FieldLevel) bool {
	length := utf8.RuneCountInString(fl.Field().String())
	param, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}
	if length > param {
		return true
	}
	return false
}

// 汉字长度小于
func ult(fl validator.FieldLevel) bool {
	length := utf8.RuneCountInString(fl.Field().String())
	param, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}
	if length < param {
		return true
	}
	return false
}

// omitemptyurl 如果参数有值 则验证url的合法性
func omitemptyurl(fl validator.FieldLevel) bool {
	//fmt.Println("fl.Field().String(): ", fl.Field().String())
	urlString := fl.Field().String()

	if urlString != "" {
		re := regexp.MustCompile(`(api|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`)
		result := re.FindAllStringSubmatch(urlString, -1)
		if result == nil {
			return false
		}
	}

	return true
}
