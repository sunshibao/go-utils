package crypto

import "encoding/base64"

//加密
func Base64Encode(src []byte) string {

	return base64.StdEncoding.EncodeToString(src)
}

//解密
func Base64Decode(s string) ([]byte, error) {

	return base64.StdEncoding.DecodeString(s)
}
