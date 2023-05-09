package check

import (
	"errors"
	"fmt"
	"reflect"
	"share/common/pkg/appError"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Check struct {
	Validate  *validator.Validate
	Translate ut.Translator
}

func NewCheck() *Check {
	va := validator.New()
	// 自定义tag注册
	RegisterAll(va)
	// 注册万能翻译
	trans := TranslateInit(va)
	return &Check{
		Validate:  va,
		Translate: trans,
	}
}

func (c *Check) Var(field interface{}, tag string) error {
	err := c.Validate.Var(field, tag)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.New(fmt.Sprintf("%v", validationErrors.Translate(c.Translate)))
		}
	}
	return nil
}

func (c *Check) Struct(s interface{}) *appError.Error {
	err := c.Validate.Struct(s)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return appError.NewError(resultSlice(validationErrors.Translate(c.Translate), s))
		} else {
			return appError.NewError(err.Error())
		}
	}
	return nil
}

func (c *Check) ValidateMap(data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
	err := c.Validate.ValidateMap(data, rules)
	return err
}

func (c *Check) RegisterAlias(alias, tags string) {
	c.Validate.RegisterAlias(alias, tags)
	return
}

func resultSlice(fields map[string]string, model interface{}) string {
	m1 := reflect.TypeOf(model) //反射模型
	// 指针Type转为非指针Type 来之资料 https://www.cnblogs.com/timelesszhuang/p/go-reflect.html
	m := m1.Elem()
	result := make([]string, 0)
	for field, err := range fields {
		fieldName := field[strings.Index(field, ".")+1:] //获得错误key

		structField, ok := m.FieldByName(fieldName) //获取反射的name
		if ok {
			structFieldTag := structField.Tag //获取Tag标签
			result = append(result, strings.ReplaceAll(err, fieldName, structFieldTag.Get("label")))
		}
	}
	if len(result) > 0 {
		return result[0]
	}
	return ""

}
