package base

import (
	"fmt"
	"time"
)

var (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
	TimeFormat     = "15:04:05"
)

//TimeToString 时间 转 字符串
func TimeToString(tTime time.Time, tFormat string) (tString string) {
	if tTime.IsZero() {
		tTime = time.Now()
	}
	if tFormat == "" {
		tFormat = DateTimeFormat
	}
	tString = tTime.Format(tFormat)

	return tString
}

//TimeToTimeNum 时间 转 时间戳
func TimeToTimeNum(tTime time.Time) (tTimeNum int64) {
	if tTime.IsZero() {
		tTime = time.Now()
	}
	tTimeNum = tTime.Unix()

	return tTimeNum
}

//TimeNumToString 时间戳 转 字符串
func TimeNumToString(tTimeNum int64, tFormat string) (tSting string) {

	if tTimeNum == 0 {
		tTimeNum = 1663819810
	}
	if tFormat == "" {
		tFormat = DateTimeFormat
	}
	tSting = time.Unix(tTimeNum, 0).Format(tFormat)
	return tSting
}

//TimeNumToTime 时间戳 转 时间
func TimeNumToTime(tTimeNum int64) (tTime time.Time) {
	if tTimeNum == 0 {
		tTimeNum = 1663819810
	}
	tTime = time.Unix(tTimeNum, 0)
	return tTime
}

//StringToTime 字符串 转 时间
func StringToTime(tString, tFormat string) (tTime time.Time) {

	if tString == "" {
		tString = "2022-09-22 12:10:10"
	}
	if tFormat == "" {
		tFormat = DateTimeFormat
	}
	local, _ := time.LoadLocation("Local")
	tTime, _ = time.ParseInLocation(tFormat, tString, local)
	return tTime
}

//StringToTimeNum 字符串 转 时间戳
func StringToTimeNum(tString, tFormat string) (tTimeNum int64) {

	if tString == "" {
		tString = "2022-09-22 12:10:10"
	}
	if tFormat == "" {
		tFormat = DateTimeFormat
	}
	local, _ := time.LoadLocation("Local")
	tTime, _ := time.ParseInLocation(tFormat, tString, local)
	return tTime.Unix()
}

//WhatTime 创建时间对象
func WhatTime() {
	// 新建一个时间对象
	now := time.Now()
	// 打印具体时间
	fmt.Println(now)
	// 对时间对象进行分段(年月日)打印
	fmt.Println(now.Year())
	fmt.Println(now.Month())
	fmt.Println(now.Day())
	fmt.Println(now.Date())
	fmt.Println(now.Hour())
	fmt.Println(now.Minute())
	fmt.Println(now.Second())
	// 打印时间戳
	fmt.Println(now.Unix())     //毫秒
	fmt.Println(now.UnixNano()) //微秒
	// time.Unix()
	// 将时间戳转换为时间对象（这里时间戳是距离1970.1.1的毫秒数）
	ret := time.Unix(1564803667, 0)
	fmt.Println(ret)
	//直接获取时间对象对应的年月日
	fmt.Println(ret.Year())
	fmt.Println(ret.Day())
	fmt.Println("---------------", ret.Unix(), "--------------")
	// 创建一个定时器
	timer := time.Tick(time.Second)
	for t := range timer {
		fmt.Println(t, "hello") // 1秒钟执行一次
	}
}

//SubTime 求时间间隔
func SubTime() {
	now := time.Now() // 本地的时间
	fmt.Println(now)
	// 明天的这个时间
	// 按照指定格式取解析一个字符串格式的时间
	time.Parse(DateTimeFormat, "2019-08-04 14:41:50")
	// 按照东八区的时区和格式取解析一个字符串格式的时间
	// 根据字符串加载时区(给需要解析的字符串设置时区)
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("load loc failed, err:%v\n", err)
		return
	}
	// 按照指定时区解析时间
	timeObj, err := time.ParseInLocation(
		DateTimeFormat,
		"2021-12-06 11:05:50", loc)
	if err != nil {
		fmt.Printf("parse time failed, err:%v\n", err)
		return
	}
	fmt.Println(timeObj)
	// 时间对象相减
	td := now.Sub(timeObj)
	fmt.Println(td)
}

func GetTimeDifference() int64 {
	nowTime := time.Now()
	// 当天秒级时间戳
	nowTimeStamp := nowTime.Unix()

	nowTimeStr := nowTime.Format(DateFormat)

	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t2, _ := time.ParseInLocation(DateFormat, nowTimeStr, time.Local)
	// 第二天零点时间戳
	towTimeStamp := t2.AddDate(0, 0, 1).Unix()

	return towTimeStamp - nowTimeStamp
}
