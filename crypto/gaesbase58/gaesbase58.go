package gaesbase58

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"io"
	"math/big"
)

// EncryptAesBase58 加密
func EncryptAesBase58(plaintext, key string) (ciphertext string, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	iv := []byte{12, 31, 78, 221, 198, 146, 65, 55, 98, 211, 165, 125, 33, 66} //14
	// Generate a random IV
	ivRandByte := make([]byte, 2)
	if _, err := io.ReadFull(rand.Reader, ivRandByte); err != nil {
		return "", err
	}
	iv = append(iv, ivRandByte...)

	stream := cipher.NewOFB(block, iv)
	encrypted := make([]byte, len(plaintext))
	stream.XORKeyStream(encrypted, []byte(plaintext))
	encrypted = append(encrypted, ivRandByte...)
	ciphertext = base58.Encode(encrypted)
	return ciphertext, nil
}

// DecryptAesBase58 解密
func DecryptAesBase58(ciphertext, key string) (string, error) {

	ciphertextByte := base58.Decode(ciphertext)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	iv := []byte{12, 31, 78, 221, 198, 146, 65, 55, 98, 211, 165, 125, 33, 66} //14
	// Generate a random IV
	ivRandByte := ciphertextByte[len(ciphertextByte)-2:]
	iv = append(iv, ivRandByte...)

	textByte := ciphertextByte[:len(ciphertextByte)-2]
	stream := cipher.NewOFB(block, iv)
	decrypted := make([]byte, len(textByte))
	stream.XORKeyStream(decrypted, textByte)
	return string(decrypted), nil
}

func IntToBase58(num int, length int) (string, error) {
	if length != 1 && length != 3 && length != 4 {
		return "", errors.New("length must be 1 or 3 or 4")
	}
	if length == 1 && num > 57 {
		return "", errors.New("length must be greater than or equal to 57")
	}
	if length == 3 && num > 190000 {
		return "", errors.New("length must be greater than or equal to 190000")
	}
	if length == 4 && num > 10000000 {
		return "", errors.New("length must be greater than or equal to 10000000")
	}
	b := make([]byte, 8)
	b[0] = byte(num >> 56)
	b[1] = byte(num >> 48)
	b[2] = byte(num >> 40)
	b[3] = byte(num >> 32)
	b[4] = byte(num >> 24)
	b[5] = byte(num >> 16)
	b[6] = byte(num >> 8)
	b[7] = byte(num)

	base58Str := base58.Encode(b)
	return base58Str[len(base58Str)-length:], nil
}

func Base58ToInt(base58Str string) int64 {

	b := base58.Decode(base58Str)

	// 创建一个big.Int和使用字节切片设置值
	intVal := new(big.Int)
	intVal.SetBytes(b)

	return intVal.Int64()
}
