package convert

import (
	"reflect"
	"strconv"
	"strings"
)

// IntToBin 十进制转二进制
func IntToBin(value interface{}) string {
	switch value.(type) {
	case int32:
		return strconv.FormatInt(int64(value.(int32)), 2)
	case int64:
		return strconv.FormatInt(value.(int64), 2)
	case int:
		return strconv.FormatInt(int64(value.(int)), 2)
	case int16:
		return strconv.FormatInt(int64(value.(int16)), 2)
	case int8:
		return strconv.FormatInt(int64(value.(int8)), 2)
	case uint:
		return strconv.FormatUint(uint64(value.(uint)), 2)
	case uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 2)
	case uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 2)
	case uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 2)
	default:
		return ""
	}
}

// IntToHex 十进制转16进制
func IntToHex(value interface{}) string {
	switch value.(type) {
	case int32:
		return strconv.FormatInt(int64(value.(int32)), 16)
	case int64:
		return strconv.FormatInt(value.(int64), 16)
	case int:
		return strconv.FormatInt(int64(value.(int)), 16)
	case int16:
		return strconv.FormatInt(int64(value.(int16)), 16)
	case int8:
		return strconv.FormatInt(int64(value.(int8)), 16)
	case uint:
		return strconv.FormatUint(uint64(value.(uint)), 16)
	case uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 16)
	case uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 16)
	case uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 16)
	default:
		return ""
	}
}

// IntToOct 十进制转8进制
func IntToOct(value interface{}) string {
	switch value.(type) {
	case int32:
		return strconv.FormatInt(int64(value.(int32)), 8)
	case int64:
		return strconv.FormatInt(value.(int64), 8)
	case int:
		return strconv.FormatInt(int64(value.(int)), 8)
	case int16:
		return strconv.FormatInt(int64(value.(int16)), 8)
	case int8:
		return strconv.FormatInt(int64(value.(int8)), 8)
	case uint:
		return strconv.FormatUint(uint64(value.(uint)), 8)
	case uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 8)
	case uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 8)
	case uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 8)
	default:
		return ""
	}
}

// BinToInt 二进制转十进制
func BinToInt(s string) int {
	var i2 = 0
	if i, err := strconv.ParseInt(s, 2, 64); err != nil {
		i2 = 0
	} else {
		i2 = int(i)
	}
	return i2
}

// HexToInt 十六进制转十进制
func HexToInt(s string) int {
	var i16 = 0
	if i, err := strconv.ParseInt(s, 16, 64); err != nil {
		i16 = 0
	} else {
		i16 = int(i)
	}
	return i16
}

// OctToInt 八进制转十进制
func OctToInt(s string) int {
	var i8 = 0
	if i, err := strconv.ParseInt(s, 8, 64); err != nil {
		i8 = 0
	} else {
		i8 = int(i)
	}
	return i8
}

// ModelsToIntSlice model中类型提取其中的 idField(int 类型) 属性组成 slice 返回
func ModelsToIntSlice(models interface{}, idField string) []int {
	if models == nil {
		return []int{}
	}

	// 类型检查
	modelsValue := reflect.ValueOf(models)
	if modelsValue.Kind() != reflect.Slice {
		return []int{}
	}

	var modelValue reflect.Value

	length := modelsValue.Len()
	ids := make([]int, 0, length)

	for i := 0; i < length; i++ {
		modelValue = reflect.Indirect(modelsValue.Index(i))
		if modelValue.Kind() != reflect.Struct {
			continue
		}

		val := modelValue.FieldByName(idField)
		if val.Kind() != reflect.Int {
			continue
		}

		ids = append(ids, int(val.Int()))
	}

	return ids
}

// UnderscoreName 驼峰式写法转为下划线写法
func UnderscoreName(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
