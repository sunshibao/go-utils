package base

import (
	"io"
	"reflect"
	"regexp"
	"strings"

	"math/rand"

	"time"

	"fmt"

	"unicode"

	"strconv"

	"github.com/coreos/etcd/pkg/idutil"
	"github.com/pborman/uuid"
)

//时间格式
const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
	TimeFormat     = "15:04:05"
)

var idGen = idutil.NewGenerator(uint16(rand.Uint32()>>16), time.Now())

func UUID() string {
	return strings.Replace(uuid.NewUUID().String(), "-", "", -1)
}

//IDGen 生成数字ID带日期
func IDGen() string {
	next := idGen.Next()
	return fmt.Sprintf("%s%d", time.Now().Format("060102"), next)
}

//ID 生成唯一数字ID
func ID() uint64 {
	return idGen.Next()
}

//RandCode IsEmail 生成随机字符串,len字符串长度,onlynumber 是否只包含数字
func RandCode(len int, onlynumber bool) string {
	if len < 1 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	ret := make([]byte, len)
	if onlynumber {
		chars := []byte("0123456789")
		for i := 0; i < len; i++ {
			ret[i] = chars[rand.Int31n(9)]
		}
	} else {
		chars := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
		for i := 0; i < len; i++ {
			ret[i] = chars[rand.Int31n(35)]
		}
	}
	return string(ret[:])
}

//IsEmail 验证是否是有效的邮箱
func IsEmail(emailaddress string) bool {
	regexpstring := `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}`
	ret, _ := regexp.MatchString(regexpstring, emailaddress)
	return ret
}

//IsPhone 验证是否是有效的手机号
func IsPhone(phonenumber string) bool {
	//regexpstring := `^1[3|4|5|7|8|9][0-9]\d{8}$`
	regexpstring := `^1\d{10}$`
	ret, _ := regexp.MatchString(regexpstring, phonenumber)
	return ret
}

//IsIdCard 验证是否有效的身份证号码
func IsIdCard(Identnumber string) bool {
	regexpstring := `^[1-9]\d{7}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}$|^[1-9]\d{5}[1-9]\d{3}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}([0-9]|X)$`
	ret, _ := regexp.MatchString(regexpstring, Identnumber)
	return ret
}

// IsChineseChar 判断是否有中文
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) ||
			(regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

// Decimal 浮点数保留指定的位数
// @value  float64 要转换的浮点数值
// @preciseDigits int 要保留的位数
func Decimal(value float64, preciseDigits int) float64 {
	precise := "%." + strconv.Itoa(preciseDigits) + "f"
	value, _ = strconv.ParseFloat(fmt.Sprintf(precise, value), 64)
	return value
}

//DateStartTimeString 日期起始时间
func DateStartTimeString(t time.Time) string {
	return t.Format(DateFormat) + " 00:00:00"
}

//DateEndTimeString 日期结束时间
func DateEndTimeString(t time.Time) string {
	return t.Format(DateFormat) + " 23:59:59"
}

// IsMark determines whether the rune is a marker
func IsMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r)
}

// IsNumeric 检查传入的字符串是否是数字 浮点型的也算数字。
func IsNumeric(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '-' && i == 0 {
			continue
		}
		if s[i] == '.' {
			if i > 0 && i < len(s)-1 {
				continue
			} else {
				return false
			}
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

//NumberM 数字取模，参数1为自然数，参数2为模
func NumberM(n int, m int) int {

	//如果参数大于等于0时
	if n >= 0 {
		return n % m
	}

	return NumberM(n+m, m) //这里必须要写返回才行，不然全得到的全是0
}

//将结构体的字段名称和字段值以字符串相连接，供sql语句使用，返回字段、值、bool
func StructFieldAndValueToString(i interface{}) (string, string, bool) {
	//获取struct的field
	fobj := reflect.TypeOf(i)

	//判断参数类型如果不是struct就返回，即使是struct对象的指针也返回，因为指针传过来也不知道它是什么类型的指针。只能知道它是ptr类型。
	if fobj.Kind() != reflect.Struct {
		return "参数不是struct类型", "", false
	}

	//声明变量，接收struct的field字符串
	var fieldStr string
	for n := 0; n < fobj.NumField(); n++ {
		if n < fobj.NumField()-1 {
			fieldStr += "`" + fobj.Field(n).Name + "`,"
		} else {
			fieldStr += "`" + fobj.Field(n).Name + "`"
		}
	}

	//获取struct的value
	vobj := reflect.ValueOf(i)
	//声明变量，接收struct的value字符串
	var valueStr string
	for m := 0; m < vobj.NumField(); m++ {
		if m < vobj.NumField()-1 {
			valueStr += "'" + fmt.Sprintf("%v", vobj.Field(m)) + "',"
		} else {
			valueStr += "'" + fmt.Sprintf("%v", vobj.Field(m)) + "'"
		}
	}
	return fieldStr, valueStr, true
}

//GetType 获取参数类型，以kind()返回，这样方便和其它类型判断
func GetType(i interface{}) reflect.Kind {
	return reflect.TypeOf(i).Kind()
}

//GetT 获取变量类型，返回系统类型的名称
func GetT(i interface{}) string { //struct
	obj := reflect.TypeOf(i)
	return obj.Kind().String()
}

//GetTT 获取变量类型，返回实际类型的名称(具体的实际类型，更细节)
func GetTT(i interface{}) string { //main.Fangchong
	return fmt.Sprintf("%T", i)
}

//CheckErr 用于检测错误，如出现错误会打印该错误并中止程序运行
func CheckErr(err error) {
	if err != nil {
		panic(err)
		//log.Fatal(err)//2019/08/20 13:05:41 strconv.Atoi: parsing "2c": invalid syntax
	}
}

//CheckErrEOF 用于检测错误，并且排除EOF，如出现错误并且不是EOF时会打印该错误并中止程序运行
func CheckErrEOF(err error) {
	if err != nil && err != io.EOF {
		panic(err)
		//log.Fatal(err)//2019/08/20 13:05:41 strconv.Atoi: parsing "2c": invalid syntax
	}
}
