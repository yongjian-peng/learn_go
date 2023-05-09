/***************************************************
 ** @Desc : This file for ...
 ** @Time : 2019/10/26 11:17
 ** @Author : yuebin
 ** @File : sort_go
 ** @Last Modified by : yuebin
 ** @Last Modified time: 2019/10/26 11:17
 ** @Software: GoLand
****************************************************/
package goutils

import (
	"github.com/spf13/cast"
	"sort"
)

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

/**
** 按照key的ascii值从小到大给map排序
 */
func SortMapByKeys(m map[string]interface{}) map[string]string {
	keys := SortMap(m)
	tmp := make(map[string]string)
	for _, key := range keys {
		tmp[key] = cast.ToString(m[key])
	}

	return tmp
}
