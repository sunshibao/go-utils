package javaDes

import (
	"crypto/des"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"log"
)

var secretKey []byte // 加解密密钥

func init() {
	key, err := sha1prng([]byte("1dsdw234sew"), 128)
	if err != nil {
		panic(err)
	}
	secretKey = generateKey(key)
}

// DESEncrypt des加密
func DESEncrypt(src []byte) ([]byte, error) {
	cipher, err := des.NewCipher(secretKey[:8])
	if err != nil {
		return nil, err
	}

	length := (len(src) + des.BlockSize) / des.BlockSize
	plain := make([]byte, length*des.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}

	// 分组分块加密
	encrypted := make([]byte, len(plain))
	for bs, be := 0, cipher.BlockSize(); bs < len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	// base64编码
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(encrypted)))
	base64.StdEncoding.Encode(buf, encrypted)

	return buf, nil
}

// DESDecrypt des解密
func DESDecrypt(dst []byte) ([]byte, error) {
	// base64解码
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(dst)))
	n, err := base64.StdEncoding.Decode(buf, dst)
	data := buf[:n]
	if err != nil {
		return nil, err
	}

	cipher, err := des.NewCipher(secretKey[:8])
	if err != nil {
		log.Printf("DES解密器创建失败! err: %v", err)
		return nil, err
	}
	decrypted := make([]byte, len(data))
	for bs, be := 0, cipher.BlockSize(); bs < len(data); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim], nil
}

func sha1prng(keyBytes []byte, encryptLength int) ([]byte, error) {
	s := sha(sha(keyBytes))
	length := encryptLength / 8
	if length > len(s) {
		return nil, errors.New("invalid length")
	}

	return s[0:length], nil
}

func sha(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

func generateKey(key []byte) []byte {
	k := make([]byte, 16)
	copy(k, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			k[j] ^= key[i]
		}
	}

	return k
}
