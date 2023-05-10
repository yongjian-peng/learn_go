package goutils

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfInt32(condition bool, trueVal, falseVal int32) int32 {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfInt64(condition bool, trueVal, falseVal int64) int64 {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfInt(condition bool, trueVal, falseVal int) int {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfBool(condition bool, trueVal, falseVal bool) bool {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfString(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

func Max(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func Min(x, y int) int {
	if x < 0 || y < 0 {
		return 0
	}
	return int(math.Min(float64(x), float64(y)))
}

// GetZodiac 获取生肖
func GetZodiac(year int) (zodiac string) {
	if year <= 0 {
		zodiac = "-1"
	}
	start := 1901
	x := (start - year) % 12
	if x == 1 || x == -11 {
		zodiac = "鼠"
	}
	if x == 0 {
		zodiac = "牛"
	}
	if x == 11 || x == -1 {
		zodiac = "虎"
	}
	if x == 10 || x == -2 {
		zodiac = "兔"
	}
	if x == 9 || x == -3 {
		zodiac = "龙"
	}
	if x == 8 || x == -4 {
		zodiac = "蛇"
	}
	if x == 7 || x == -5 {
		zodiac = "马"
	}
	if x == 6 || x == -6 {
		zodiac = "羊"
	}
	if x == 5 || x == -7 {
		zodiac = "猴"
	}
	if x == 4 || x == -8 {
		zodiac = "鸡"
	}
	if x == 3 || x == -9 {
		zodiac = "狗"
	}
	if x == 2 || x == -10 {
		zodiac = "猪"
	}
	return
}

// GetAgeByYear y, m, d := GetTimeFromStrDate("1993-08-20")
//fmt.Println(GetAge(y),GetConstellation(m, d),GetZodiac(y))
func GetAgeByYear(year int) (age int) {
	if year <= 0 {
		age = -1
	}
	nowyear := time.Now().Year()
	age = nowyear - year
	return
}

// GetConstellation 获取星座
func GetConstellation(month, day int) (star string) {
	if month <= 0 || month >= 13 {
		star = "-1"
	}
	if day <= 0 || day >= 32 {
		star = "-1"
	}
	if (month == 1 && day >= 20) || (month == 2 && day <= 18) {
		star = "水瓶座"
	}
	if (month == 2 && day >= 19) || (month == 3 && day <= 20) {
		star = "双鱼座"
	}
	if (month == 3 && day >= 21) || (month == 4 && day <= 19) {
		star = "白羊座"
	}
	if (month == 4 && day >= 20) || (month == 5 && day <= 20) {
		star = "金牛座"
	}
	if (month == 5 && day >= 21) || (month == 6 && day <= 21) {
		star = "双子座"
	}
	if (month == 6 && day >= 22) || (month == 7 && day <= 22) {
		star = "巨蟹座"
	}
	if (month == 7 && day >= 23) || (month == 8 && day <= 22) {
		star = "狮子座"
	}
	if (month == 8 && day >= 23) || (month == 9 && day <= 22) {
		star = "处女座"
	}
	if (month == 9 && day >= 23) || (month == 10 && day <= 22) {
		star = "天秤座"
	}
	if (month == 10 && day >= 23) || (month == 11 && day <= 21) {
		star = "天蝎座"
	}
	if (month == 11 && day >= 22) || (month == 12 && day <= 21) {
		star = "射手座"
	}
	if (month == 12 && day >= 22) || (month == 1 && day <= 19) {
		star = "魔蝎座"
	}

	return star
}

// GetAgeBirthDay 获取年龄
func GetAgeBirthDay(birthday string) (string, int32) {
	year, month, day := GetDateTimeFromStrDate(birthday)
	var age int32
	if year != 0 && month != 0 && day != 0 {
		//获取年龄
		age = int32(GetAgeByYear(year))
	} else {
		age = 22
		nowYear := time.Now().Year()
		bYear := int32(nowYear) - age
		birthday = fmt.Sprintf("%s-01-01", strconv.Itoa(int(bYear)))
	}
	return birthday, age
}

func MinInt32(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func MaxInt32(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func GetMinMaxUid(one, two int32) string {
	min := MinInt32(one, two)
	max := MaxInt32(one, two)
	return fmt.Sprintf("%d_%d", min, max)
}

// Abs 获取绝对值
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func GetRandSource() *rand.Rand {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r
}

// Random 获取随机数
func Random(source *rand.Rand, min, max int) int {
	randNum := source.Intn(max-min) + min
	return randNum
}

// GetFriendShowTime /获取友好时间
func GetFriendShowTime(curSecond int32) string {
	if curSecond > 86400 {
		return fmt.Sprintf("%d天", curSecond/86400)
	} else if curSecond > 3600 {
		return fmt.Sprintf("%d小时", curSecond/3600)
	} else if curSecond > 60 {
		return fmt.Sprintf("%d分钟", curSecond/60)
	} else {
		return fmt.Sprintf("%d秒", curSecond)
	}
}

// GetDistance 返回值的单位为米
func GetDistance(lat1, lng1, lat2, lng2 float64) string {
	radius := float64(6371000) // 6378137
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	distance := dist * radius
	if distance > 1000 {
		return fmt.Sprintf("%skm", strconv.FormatFloat(distance/1000, 'f', 0, 64))
	}
	return fmt.Sprintf("%sm", strconv.FormatFloat(distance, 'f', 0, 64))
}

func Dump(result interface{}) {
	fmt.Printf("%#v \n", result)
}
