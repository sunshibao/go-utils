package javaDes

import (
	"fmt"
	"testing"
)

func TestJavaDes(t *testing.T) {
	encrypt, _ := DESEncrypt([]byte("孙世宝"))
	fmt.Println(string(encrypt))

	decrypt, _ := DESDecrypt(encrypt)
	fmt.Println(string(decrypt))
}
