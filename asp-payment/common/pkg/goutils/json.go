package goutils

import (
	"github.com/bytedance/sonic"
)

// JsonEncode json转码
func JsonEncode(data interface{}) (string, error) {
	res, err := sonic.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// JsonDecode json转码
func JsonDecode(str string, data interface{}) error {
	err := sonic.Unmarshal([]byte(str), data)
	if err != nil {
		return err
	}
	return nil
}

// JsonDecodeByte json转码 入参是 []byte
func JsonDecodeByte(str []byte, data interface{}) error {
	err := sonic.Unmarshal(str, data)
	if err != nil {
		return err
	}
	return nil
}
