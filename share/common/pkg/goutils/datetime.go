/***************************************************
 ** @Desc : This file for ...
 ** @Time : 2019/8/21 10:21
 ** @Author : yuebin
 ** @File : date_time
 ** @Last Modified by : yuebin
 ** @Last Modified time: 2019/8/21 10:21
 ** @Software: GoLand
****************************************************/
package goutils

import "time"

func GetDateTimeNot() string {
	return time.Now().Format("2006010215:04:05")
}

func GetDate() string {
	return time.Now().Format("2006-01-02")
}

func GetBasicDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetDateTimeUnix() int64 {
	return time.Now().Unix()
}

func GetDateTimeUnixNano() int64 {
	return time.Now().UnixNano()
}

func GetNowTimesTamp() string {
	return time.Now().Format("20060102150405")
}

func GetDateTimeBeforeHours(hour int) string {
	return time.Now().Add(-time.Hour * time.Duration(hour)).Format("2006-01-02 15:04:05")
}

func GetDateBeforeDays(days int) string {
	return time.Now().Add(-time.Hour * time.Duration(days) * 24).Format("2006-01-02")
}

func GetDateTimeBeforeDays(days int) string {
	return time.Now().Add(-time.Hour * time.Duration(days) * 24).Format("2006-01-02 15:04:05")
}

func GetDateAfterDays(days int) string {
	return time.Now().Add(time.Hour * time.Duration(days) * 24).Format("2006-01-02")
}

// 获取今天0点0时0分的时间戳 1664380800
func GetTodayBeginTimeStamp() int64 {
	currentTime := time.Now()

	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	return startTime.Unix()
}

// 获取今天0点0时0分的时间戳 2022-09-29 00:00:00
func GetTodayBeginTimeData() string {
	currentTime := time.Now()

	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	return startTime.Format("2006/01/02 15:04:05")
}

// 获取今天23:59:59秒的时间戳
func GetTodayEndTimeStamp() int64 {
	currentTime := time.Now()

	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())
	return endTime.Unix()
}

// 获取今天23:59:59秒的时间戳
func GetTodayEndTimeData() string {
	currentTime := time.Now()

	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())
	return endTime.Format("2006/01/02 15:04:05")
}

// 时间格式 转时间戳 例如 2022-10-21 20:36:38 转时间戳
func GetTimesTampToUnix(timeStamp string) (timeInt int64) {
	loc, _ := time.LoadLocation("Local")
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStamp, loc) //使用parseInLocation将字符串格式化返回本地时区时间
	timeInt = stamp.Unix()
	return
}
