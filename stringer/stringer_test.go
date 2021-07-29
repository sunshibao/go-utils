package stringer

import (
	"fmt"
	"testing"
	"time"
)

func TestReverse(t *testing.T) {
	now := time.Now()
	str := "这是一个字符串反正程序，非常好用。"
	reverse := Reverse(str)
	fmt.Println(reverse)
	elapsed := time.Since(now)
	fmt.Println("app run time", elapsed)
}
