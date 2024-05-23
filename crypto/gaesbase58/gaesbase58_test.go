package gaesbase58

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"testing"
)

// 加密
func TestEncryptAesBase58(t *testing.T) {

	var (
		aesKey   = "1q^rfW@dnHGIu06Y" //AES KEY
		themeId  = 232                //产品ID
		deadline = "240530"           //过期时间
		num1     = 5                  //生成个数
		num2     = 1                  //生成次数
	)

	bas58ThemeId, err := IntToBase58(themeId, 3) //给了3个长度 最大58的3次方
	if err != nil {
		fmt.Println(err)
	}
	tempPlaintext, err := IntToBase58(gconv.Int(deadline), 4) //给了4个长度，最大58的4次方
	if err != nil {
		fmt.Println(err)
	}
	base58Num2, err := IntToBase58(num2, 1) //给了1个长度，最大58
	if err != nil {
		fmt.Println(err)
	}
	//生成兑换码
	for i := 0; i < num1; i++ {
		base58Num1, err := IntToBase58(i, 4) //给了4个长度，最大58的4次方
		//最终明文
		fmt.Println(fmt.Sprintf("base58后的明文:%s%s%s%s", bas58ThemeId, tempPlaintext, base58Num1, base58Num2))
		ciphertext, err := EncryptAesBase58(fmt.Sprintf("%s%s%s%s", bas58ThemeId, tempPlaintext, base58Num1, base58Num2), aesKey)
		if err != nil {
			continue
		}
		fmt.Printf("最终的密文:num:%d,plainText:%s \n", i, ciphertext)
	}

}

// 解密
func TestDecryptAesBase58(t *testing.T) {
	var (
		aesKey     = "1q^rfW@dnHGIu06Y"     //AES KEY
		ciphertext = "2W6aEX9Qsjc2pNfGw9Cw" //密文
	)

	plainText, err := DecryptAesBase58(ciphertext, aesKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("解密后的base58明文", plainText)

	newThemeId := Base58ToInt(plainText[0:3])
	newValidate := Base58ToInt(plainText[3:7])
	num1 := Base58ToInt(plainText[7:11])
	num2 := Base58ToInt(plainText[11:12])

	fmt.Printf("最终明文：themeId:%d,date:%d,num1:%d,num2:%d \n", newThemeId, newValidate, num1, num2)

}
