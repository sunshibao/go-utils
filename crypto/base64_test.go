package crypto

import (
	"encoding/base64"
	"fmt"
	"testing"
)

//base64加解密测试
func TestBase64(t *testing.T) {
	//base64加密
	byteStr := []byte("张三是好人")
	str := base64.StdEncoding.EncodeToString(byteStr)
	fmt.Println("base64加密后：", str) //5YiY6Ziz5piv5aW95Lq6

	//base64解密
	deByte, deErr := base64.StdEncoding.DecodeString(str)
	if deErr != nil {
		fmt.Println(deErr)
	}
	fmt.Println("base64解密后：", string(deByte)) //张三是好人

}
