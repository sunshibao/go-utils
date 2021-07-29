package moment

import (
	"time"
)

const (
	DateFormat          = "2006-01-02"
	DateTerseFormat     = "20060102"
	DateTimeFormat      = "2006-01-02 15:04:05"
	DateTimeTerseFormat = "20060102150405"
)

// FormatDate 格式时间为 YYYY-MM-DD 形式
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// FormatDateTerse 格式时间为 YYYYMMDD 形式
func FormatDateTerse(t time.Time) string {
	return t.Format(DateTerseFormat)
}

// FormatDateTime 格式时间为 YYYY-MM-DD HH:II:SS 形式
func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

// FormatDateTimeTerse 格式时间为 YYYYMMDDHHIISS 形式
func FormatDateTimeTerse(t time.Time) string {
	return t.Format(DateTimeTerseFormat)
}

// ParseDate 解析 YYYY-MM-DD 形式的时间字符串
func ParseDate(s string) (time.Time, error) {
	return time.ParseInLocation(DateFormat, s, time.Local)
}

// 解析 YYYYMMDD 形式的时间字符串
func ParseDateTerse(s string) (time.Time, error) {
	return time.ParseInLocation(DateTerseFormat, s, time.Local)
}

// ParseDateTime 解析 YYYY-MM-DD HH:II:SS 形式的时间字符串
func ParseDateTime(s string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormat, s, time.Local)
}

// ParseDateTimeTerse 解析 YYYYMMDDHHIISS 形式的时间字符串
func ParseDateTimeTerse(s string) (time.Time, error) {
	return time.ParseInLocation(DateTimeTerseFormat, s, time.Local)
}

// 解析 YYYY-MM-DD 形式的时间字符串（忽略错误）
func MustParseDate(s string) time.Time {
	t, _ := ParseDate(s)
	return t
}

// MustParseDateTerse 解析 YYYYMMDD 形式的时间字符串（忽略错误）
func MustParseDateTerse(s string) time.Time {
	t, _ := ParseDateTerse(s)
	return t
}

// MustParseDateTime 解析 YYYY-MM-DD HH:II:SS 形式的时间字符串（忽略错误）
func MustParseDateTime(s string) time.Time {
	t, _ := ParseDateTime(s)
	return t
}

// MustParseDateTimeTerse 解析 YYYYMMDDHHIISS 形式的时间字符串（忽略错误）
func MustParseDateTimeTerse(s string) time.Time {
	t, _ := ParseDateTimeTerse(s)
	return t
}
