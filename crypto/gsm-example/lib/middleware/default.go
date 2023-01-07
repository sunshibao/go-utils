package middleware

import (
	"bytes"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"utils/crypto/gsm-example/lib/core"
	"utils/crypto/gsm-example/lib/smCrypto"
)

//解密客户端参数
func CryptoCheck() core.HandlerFunc {
	return func(c *core.Context) {
		if c.GetHeader("raw") == "real" {
			c.Next()
			return
		}
		c.Request.Header.Set("crypto-type", "sm")
		sign := c.GetHeader("sign")

		//处理query参数
		urlQueryData := c.Request.URL.RawQuery
		if urlQueryData != "" {
			//数据验签
			if !smCrypto.Sm2Crypto.VerifyHex(urlQueryData, sign) {
				c.Mix(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "请求数据错误!",
				})
				c.Abort()
				return
			}
			//数据解密
			decrypted, err := smCrypto.Sm4Crypto.EcbDecodeBase64(urlQueryData)
			if err != nil || len(decrypted) == 0 {
				c.Mix(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "数据解密失败!",
				})
				c.Abort()
				return
			}

			c.Request.URL.RawQuery = decrypted
		}

		//处理body参数
		buf, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Mix(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "数据读取失败",
			})
			c.Abort()
			return
		}
		if len(buf) != 0 {
			//数据验签
			if !smCrypto.Sm2Crypto.VerifyHex(string(buf), sign) {
				c.Mix(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "请求数据错误!",
				})
				c.Abort()
				return
			}
			//数据解密
			bodyByte, err := base64.StdEncoding.DecodeString(string(buf))
			if err != nil {
				c.Mix(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "数据解密失败!",
				})
				c.Abort()
				return
			}

			decrypted, err := smCrypto.Sm4Crypto.EcbDecode(bodyByte)
			if err != nil || len(decrypted) == 0 {
				c.Mix(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "数据解密失败!",
				})
				c.Abort()
				return
			}

			c.Request.Header.Set("content-type", "application/json")
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(decrypted))
		}

		c.Next()
	}
}
