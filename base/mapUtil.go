package base

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	replaceCharReg, _ = regexp.Compile(`[\-\.\_\s]+`)
)

type M[K string, V any] map[K]V

//将struct类型数据转换成map类型数据，转不了，不要再试了。也没必要转，直接将struct转成map岂不更好。

//将map数据转成post提交所需格式的字符串
func MapToString(mm map[string]interface{}) string {
	//样例：name=65.56&sex=25&age=true&
	//将最后面的字符"&"去掉，有两种方式：

	//方式1：
	//str := ""//总接收变量
	//temp := ""//临时接收拼接的变量
	//for key,value := range mm {
	//	zhi := ValueToString(value)
	//	temp = key + "=" + zhi + "&"
	//	str += temp
	//}
	////由于map的key值顺序是不固定的，所以无法判断最后一个key是什么值，所以只能用最后出现的"&"符号来定位。
	//index := strings.LastIndex(str,"&")
	//result := str[0:index]
	//return result

	//方式2：
	str := ""      //总接收变量
	temp := ""     //临时接收拼接的变量
	i := 0         //自增变量
	cnt := len(mm) //map元素个数
	for key, value := range mm {
		zhi := ValueToString(value)
		if i < (cnt - 1) {
			temp = key + "=" + zhi + "&"
		} else {
			temp = key + "=" + zhi
		}
		str += temp
		i++ //临时变量自增
	}
	return str
}

// MapCopy does a shallow copy from map <data> to <copy> for most commonly used map type
// map[string]interface{}.
func MapCopy(data map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return
}

// MapContains checks whether map <data> contains <key>.
func MapContains(data map[string]interface{}, key string) (ok bool) {
	if len(data) == 0 {
		return
	}
	_, ok = data[key]
	return
}

// MapDelete deletes all <keys> from map <data>.
func MapDelete(data map[string]interface{}, keys ...string) {
	if len(data) == 0 {
		return
	}
	for _, key := range keys {
		delete(data, key)
	}
}

// MapMerge merges all map from <src> to map <dst>.
func MapMerge(dst map[string]interface{}, src ...map[string]interface{}) {
	if dst == nil {
		return
	}
	for _, m := range src {
		for k, v := range m {
			dst[k] = v
		}
	}
}

// MapMergeCopy creates and returns a new map which merges all map from <src>.
func MapMergeCopy(src ...map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{})
	for _, m := range src {
		for k, v := range m {
			copy[k] = v
		}
	}
	return
}

// MapReplace returns a copy of <origin>,
// which is replaced by a map in unordered way, case-sensitively.
func MapReplace(origin string, replaces map[string]string) string {
	for k, v := range replaces {
		origin = strings.Replace(origin, k, v, -1)
	}
	return origin
}

// MapPossibleItemByKey tries to find the possible key-value pair for given key ignoring cases and symbols.
//
// Note that this function might be of low performance.
func MapPossibleItemByKey(data map[string]interface{}, key string) (foundKey string, foundValue interface{}) {
	if len(data) == 0 {
		return
	}
	if v, ok := data[key]; ok {
		return key, v
	}
	// Loop checking.
	for k, v := range data {
		if EqualFoldWithoutChars(k, key) {
			return k, v
		}
	}
	return "", nil
}

// EqualFoldWithoutChars checks string <s1> and <s2> equal case-insensitively,
// with/without chars '-'/'_'/'.'/' '.
func EqualFoldWithoutChars(s1, s2 string) bool {
	return strings.EqualFold(
		replaceCharReg.ReplaceAllString(s1, ""),
		replaceCharReg.ReplaceAllString(s2, ""),
	)
}

// MapContainsPossibleKey checks if the given <key> is contained in given map <data>.
// It checks the key ignoring cases and symbols.
//
// Note that this function might be of low performance.
func MapContainsPossibleKey(data map[string]interface{}, key string) bool {
	if k, _ := MapPossibleItemByKey(data, key); k != "" {
		return true
	}
	return false
}

// MapOmitEmpty deletes all empty values from given map.
func MapOmitEmpty(data map[string]interface{}) {
	if len(data) == 0 {
		return
	}
	for k, v := range data {
		if IsEmpty(v) {
			delete(data, k)
		}
	}
}

// apiString is used for type assert api for String().
type apiString interface {
	String() string
}

// apiInterfaces is used for type assert api for Interfaces.
type apiInterfaces interface {
	Interfaces() []interface{}
}

// apiMapStrAny is the interface support for converting struct parameter to map.
type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// IsEmpty checks whether given <value> empty.
// It returns true if <value> is in: 0, nil, false, "", len(slice/map/chan) == 0,
// or else it returns false.
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch value := value.(type) {
	case int:
		return value == 0
	case int8:
		return value == 0
	case int16:
		return value == 0
	case int32:
		return value == 0
	case int64:
		return value == 0
	case uint:
		return value == 0
	case uint8:
		return value == 0
	case uint16:
		return value == 0
	case uint32:
		return value == 0
	case uint64:
		return value == 0
	case float32:
		return value == 0
	case float64:
		return value == 0
	case bool:
		return value == false
	case string:
		return value == ""
	case []byte:
		return len(value) == 0
	case []rune:
		return len(value) == 0
	default:
		// Common interfaces checks.
		if f, ok := value.(apiString); ok {
			return f.String() == ""
		}
		if f, ok := value.(apiInterfaces); ok {
			return len(f.Interfaces()) == 0
		}
		if f, ok := value.(apiMapStrAny); ok {
			return len(f.MapStrAny()) == 0
		}
		// Finally using reflect.
		var rv reflect.Value
		if v, ok := value.(reflect.Value); ok {
			rv = v
		} else {
			rv = reflect.ValueOf(value)
		}
		switch rv.Kind() {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			return rv.Len() == 0

		case reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return true
			}
		}
	}
	return false
}

// IsNil checks whether given <value> is nil.
// Note that it might use reflect feature which affects performance a little bit.
func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}
	var rv reflect.Value
	if v, ok := value.(reflect.Value); ok {
		rv = v
	} else {
		rv = reflect.ValueOf(value)
	}
	switch rv.Kind() {
	case reflect.Chan,
		reflect.Map,
		reflect.Slice,
		reflect.Func,
		reflect.Ptr,
		reflect.Interface,
		reflect.UnsafePointer:
		return rv.IsNil()
	}
	return false
}

//interface值转string
func ValueToString(i interface{}) string {
	//fmt.Println(i)//打印参数值

	str := ""
	//用这种方法判断，就省去了reflect.TypeOf(i)反射的判断，如：
	//obj := reflect.TypeOf(i)
	//if obj.Kind() == reflect.Int {}
	switch idata := i.(type) {
	case string:
		str = idata
	case int:
		str = strconv.Itoa(idata)
	case int8:
		str = strconv.Itoa(int(idata))
	case int16:
		str = strconv.Itoa(int(idata))
	case int32:
		str = strconv.Itoa(int(idata))
	case int64:
		str = strconv.FormatInt(idata, 10)
	case uint:
		str = strconv.Itoa(int(idata))
	case uint8:
		str = strconv.Itoa(int(idata))
	case uint16:
		str = strconv.Itoa(int(idata))
	case uint32:
		str = strconv.FormatInt(int64(idata), 10)
	case uint64:
		str = strconv.FormatInt(int64(idata), 10)
	case float32:
		str = strconv.FormatFloat(float64(idata), 'f', -1, 32)
	case float64:
		str = strconv.FormatFloat(idata, 'f', -1, 64)
	case bool:
		str = strconv.FormatBool(idata)
	case []byte:
		str = string(idata)
	default:
		str = "error" //未知类型时，返回error字符串
	}

	return str
}
