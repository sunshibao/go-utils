package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

//PasswordEncode 将密码加密，此方法比md5要好，刘阳推荐
func PasswordEncode(password string) (string, error) {
	byteArr, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(byteArr), nil
}

//PasswordVerify 将密码验证，参数1：加密之后的密码，参数2：新密码
func PasswordVerify(pwd_encode string, new_pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwd_encode), []byte(new_pwd))
	if err != nil {
		return false
	}

	return true
}
