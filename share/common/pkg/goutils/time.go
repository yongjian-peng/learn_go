package goutils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// GetDateTimeFromStrDate 获取年月日根据日期
func GetDateTimeFromStrDate(date string) (year, month, day int) {
	const shortForm = "2006-01-02"
	d, err := time.Parse(shortForm, date)
	if err != nil {
		//fmt.Println("出生日期解析错误！")
		return 0, 0, 0
	}
	year = d.Year()
	month = int(d.Month())
	day = d.Day()
	return
}

// GetCurDate 获取当前日期
func GetCurDate() string {
	return time.Now().Format("2006-01-02")
}

// GetCurDateShort 获取当前日期
func GetCurDateShort() string {
	return time.Now().Format("20060102")
}

// GetCurDateHourShort 获取当前日期小时
func GetCurDateHourShort() string {
	return time.Now().Format("2006010215")
}

// GetCurDateMinuteShort 获取当前分钟
func GetCurDateMinuteShort() string {
	return time.Now().Format("200601021504")
}

// GetPreDateMinuteShort 获取前一分钟的日期
func GetPreDateMinuteShort() string {
	return time.Now().Add(time.Minute * -1).Format("200601021504")
}

func StrTimeToShortTimeStr(strTime string) string {
	tt, err := time.Parse("2006-01-02 15:04:05", strTime)
	if err != nil {
		return time.Now().Format("2006-01-02 15:04")
	}
	return tt.Format("2006-01-02 15:04")
}

// GetCurDateYearWeek 获取当前年月
func GetCurDateYearWeek() int {
	y, w := time.Now().ISOWeek()
	return y*100 + w
}

// GetPreDateShort 获取昨天的日期
func GetPreDateShort() string {
	return time.Now().AddDate(0, 0, -1).Format("20060102")
}

// GetPreDateTimeStr 获取昨天的时间
func GetPreDateTimeStr() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04")
}

func GetCurDateYearMonth() string {
	return time.Now().Format("2006-01")
}

func GetCurDateYearMonthShort() string {
	return time.Now().Format("200601")
}

func GetPreDateYearMonthShort() string {
	return time.Now().AddDate(0, -1, 0).Format("200601")
}

// GetCurTimeStr 获取当前时间
func GetCurTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetCurTimeUnixSecond 获取当前秒时间搓 10位
func GetCurTimeUnixSecond() int64 {
	return time.Now().Unix()
}

// GetCurTimeMillisecond 获取当前毫秒数 13位
func GetCurTimeMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetCurTimeUnixNano 获取当前纳秒时间搓 16位
func GetCurTimeUnixNano() int64 {
	return time.Now().UnixNano() / 1e3
}

// TimeFormat Format 跟 PHP 中 date 类似的使用方式，如果 ts 没传递，则使用当前时间
func TimeFormat(format string, ts ...time.Time) string {
	patterns := []string{
		// 年
		"Y", "2006", // 4 位数字完整表示的年份
		"y", "06", // 2 位数字表示的年份

		// 月
		"m", "01", // 数字表示的月份，有前导零
		"n", "1", // 数字表示的月份，没有前导零
		"M", "Jan", // 三个字母缩写表示的月份
		"F", "January", // 月份，完整的文本格式，例如 January 或者 March

		// 日
		"d", "02", // 月份中的第几天，有前导零的 2 位数字
		"j", "2", // 月份中的第几天，没有前导零

		"D", "Mon", // 星期几，文本表示，3 个字母
		"l", "Monday", // 星期几，完整的文本格式;L的小写字母

		// 时间
		"g", "3", // 小时，12 小时格式，没有前导零
		"G", "15", // 小时，24 小时格式，没有前导零
		"h", "03", // 小时，12 小时格式，有前导零
		"H", "15", // 小时，24 小时格式，有前导零

		"a", "pm", // 小写的上午和下午值
		"A", "PM", // 小写的上午和下午值

		"i", "04", // 有前导零的分钟数
		"s", "05", // 秒数，有前导零
	}
	replacer := strings.NewReplacer(patterns...)
	format = replacer.Replace(format)

	t := time.Now()
	if len(ts) > 0 {
		t = ts[0]
	}
	return t.Format(format)
}

func StrToTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",

		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return t
		}
	}

	return time.Now()
}

func UnixTimeToStr(unixTime int64, format string) string {
	return TimeFormat(format, time.Unix(unixTime, 0))
}

func StrTimeToUnixTime(timeStr string) int64 {
	return StrToTime(timeStr).Unix()
}

// MonthDayNum t 所在时间的月份总天数
func MonthDayNum(t time.Time) int {
	isLeapYear := isLeap(t.Year())
	month := t.Month()
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.February:
		if isLeapYear {
			return 29
		}
		return 28
	default:
		return 30
	}
}

func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	hours := diff.Hours()
	if hours < 1.0 {
		return fmt.Sprintf("约 %.0f 分钟前", diff.Minutes())
	}

	if hours < 24.0 {
		return fmt.Sprintf("约 %.0f 小时前", hours)
	}

	if hours < 72.0 {
		return fmt.Sprintf("约 %.0f 天前", hours/24.0)
	}

	// 同一年，不用年份
	if now.Year() == t.Year() {
		return t.Format("01-02 15:04")
	}

	return t.Format("2006-01-02")
}

// 是否闰年
func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func DateSub(date1, date2 time.Time) int {
	//计算相差天数
	return int(math.Ceil(date1.Sub(date2).Seconds() / 86400))
}
