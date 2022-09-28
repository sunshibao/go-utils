package base

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

/*
type Person struct {
	User string `json:"user"`
	Age int `json:"age"`
	Sex string `json:"sex"`
	Tels []Tel `json:"tels"`
}

type Tel struct {
	Money float64 `json:"money"`
	Num int `json:"num"`
}

func Http_Client_test()  {

	//post提交普通数据
	mm := make(map[string]interface{}, 0)
	mm["mobile"] = "13000000001"
	mm["password"] = "123456"

	//测试post提交
	url := "http://localhost:8090/user/login"

	result, err := HttpRequestPOST(url,mm)
	fmt.Println("是否错误：",err)//<nil>
	fmt.Println("返回值为：",result)//{"code":0,"data":{"info":{"uid":6660001}, ... }


	fmt.Println("===================================================")

	//post提交json数据
	pp := &Person{
		User: "刘阳",
		Age:  25,
		Sex:  "男",
		Tels:[]Tel{
			{
				Money:25.56,
				Num:3,
			},
			{
				Money:6.8,
				Num:2,
			},
		},
	}
	aa,_ := json.Marshal(pp)
	jsonstr := string(aa)

	url1 := "http://localhost:8081/json"
	result1, err1 := HttpRequestPOSTJSON(url1,jsonstr)
	fmt.Println("是否错误：",err1)
	fmt.Println("返回值为：",result1)
}
*/
//http请求操作===========================================

//http请求,get方式
func HttpRequestGET(url string) (string, error) {
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil) //建立一个请求
	//返回的err没有用
	//defer reqest.Body.Close()//设置关闭会导致错误，不要写。
	//设置头协议
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("Cookie", "设置cookie")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	response, resErr := client.Do(reqest) //提交请求
	if resErr != nil {
		return "", resErr
	}
	defer response.Body.Close()
	//cookies := response.Cookies() //遍历cookies
	//for _, cookie := range cookies {
	//	fmt.Println("cookie:", cookie)
	//}

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	return "", fmt.Errorf("服务器响应异常，状态：%s", response.Status)

}

//http请求,post方式,post内容为常规字段和内容。非json。
func HttpRequestPOST(url string, data map[string]interface{}) (string, error) {
	client := &http.Client{}
	//reqest, _ := http.NewRequest("POST",url,strings.NewReader("name=刘阳&sex=男&age=30"))
	//本来传参为string,为了方便传map，再由MapToString()方法转string
	str := HttpMapToString(data) //调用map转string函数
	reqest, _ := http.NewRequest("POST", url, strings.NewReader(str))
	//返回的err没有用
	//defer reqest.Body.Close()//设置关闭会导致错误，不要写。
	//设置头协议
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("Cookie", "设置cookie")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	response, resErr := client.Do(reqest) //提交请求
	if resErr != nil {
		return "", resErr
	}
	defer response.Body.Close()
	//cookies := response.Cookies()
	//for _, cookie := range cookies {
	//	fmt.Println("cookie:", cookie)
	//}

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	return "", fmt.Errorf("服务器响应异常，状态：%s", response.Status)
}

//http发起post请求，提交数据为：json字符串。前题是服务端设置了接收json类型的数据提交。
func HttpRequestPOSTJSON(url, jsonstr string) (string, error) {

	requestBody := fmt.Sprintf(`%s`, jsonstr)

	jsonByte := []byte(requestBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, resErr := client.Do(req) //提交请求

	if resErr != nil {
		return "", resErr
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	return "", fmt.Errorf("服务器响应异常，状态：%s", resp.Status)
}

//将map数据转成post提交所需格式的字符串
func HttpMapToString(mm map[string]interface{}) string {
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
		zhi := HttpValueToString(value)
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

//interface值转string
func HttpValueToString(i interface{}) string {
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
