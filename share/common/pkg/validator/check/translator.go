package check

import (
	enLocal "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
)

// 初始化翻译器
func TranslateInit(validate *validator.Validate) ut.Translator {
	enCh := enLocal.New()
	uni := ut.New(enCh)                 // 万能翻译器，保存所有的语言环境和翻译数据
	trans, _ := uni.GetTranslator("en") // 翻译器
	_ = en.RegisterDefaultTranslations(validate, trans)

	// 添加额外翻译
	_ = validate.RegisterTranslation("ulen", trans, func(ut ut.Translator) error {
		return ut.Add("ulen", "{0} Length equal to{1}!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ulen", fe.Field(), fe.Param())
		return t
	})

	_ = validate.RegisterTranslation("omitemptyurl", trans, func(ut ut.Translator) error {
		return ut.Add("omitemptyurl", "{0} Not a legal address {1}!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("omitemptyurl", fe.Field(), fe.Param())
		return t
	})

	return trans
}
