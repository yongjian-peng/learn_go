package main

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

type UserModel struct {
	Name     string `label:"姓名" json:"name" validate:"required"`
	Age      int    `label:"年龄" json:"age" validate:"required"`
	Birthday string `label:"生日" json:"birthday" validate:"required"`
}

func main() {
	//声明翻译对象
	uni := ut.New(en.New(), zh.New())
	//设置翻译语言
	trans, _ := uni.GetTranslator("zh")
	//创建一个验证数据
	user := UserModel{
		Name:     "",
		Age:      0,
		Birthday: "",
	}
	//获取验证对象
	validate := validator.New()
	//注册默认翻译
	_ = zhtrans.RegisterDefaultTranslations(validate, trans)
	//验证数据
	err := validate.Struct(user)
	if err != nil {
		errs := err.(validator.ValidationErrors) //转换验证结果
		fmt.Println(errs.Translate(trans), "yyyyyyy")
		//翻译验证结果并转换为 slice 形式映射
		fmt.Println(fmt.Sprintf("--- : %+v", resultSlice(errs.Translate(trans), user)))
		fmt.Println(fmt.Sprintf(fmt.Sprintf("+++++++ : %+v", resultMap(errs.Translate(trans), user))))
		//翻译验证结果并转换为 map 形式映射

	}

}

func resultMap(fields map[string]string, model interface{}) map[string]string {
	m := reflect.TypeOf(model)       //反射模型
	result := map[string]string{}    //声明错误数组
	for field, err := range fields { //遍历错误
		fieldName := field[strings.Index(field, ".")+1:] //获得错误key
		structField, _ := m.FieldByName(fieldName)       //获取反射的name
		structFieldTag := structField.Tag                //获取Tag标签
		jsonName := structFieldTag.Get("json")           //获取Tag标签 json 信息
		if jsonName == "-" {                             //忽略字段
			continue
		}
		if jsonName == "" { //为空则为struct的名字
			jsonName = fieldName
		}
		//将错误信息err中的英文名替换为 中文信息 并赋值给返回参
		result[jsonName] = strings.ReplaceAll(err, fieldName, structFieldTag.Get("label"))
	}
	return result
}

func resultSlice(fields map[string]string, model interface{}) string {
	m := reflect.TypeOf(model) //反射模型
	fmt.Println("model: ", model)
	result := make([]string, 0)
	for field, err := range fields {
		fieldName := field[strings.Index(field, ".")+1:] //获得错误key
		structField, ok := m.FieldByName(fieldName)      //获取反射的name
		fmt.Println("ok: ", ok)
		fmt.Println("structField: ", structField)
		structFieldTag := structField.Tag //获取Tag标签
		result = append(result, strings.ReplaceAll(err, fieldName, structFieldTag.Get("label")))
	}

	if len(result) > 0 {
		return result[0]
	}

	return ""

}
