package core

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"utils/crypto/gsm-example/lib/smCrypto"
)

type Context struct {
	*gin.Context
}

type HandlerFunc func(ctx *Context)

func Handler(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{c}
		h(ctx)
	}
}

func (c *Context) Mix(code int, obj interface{}) {
	//cryptoType := c.GetHeader("crypto-type")
	if c.GetHeader("raw") == "real" || code != 200 {
		c.JSON(code, obj)
		return
	}

	//对数据进行加密
	b, _ := json.Marshal(obj)
	cryptBuf, err := smCrypto.Sm4Crypto.EcbEncodeBase64(string(b))
	if err != nil || cryptBuf == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "返回数据加密错误",
		})
		return
	}

	sign, err := smCrypto.Sm2Crypto.SignHex(cryptBuf)
	if err != nil || cryptBuf == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "签名错误",
		})
		return
	}
	c.Header("sign", sign)
	c.String(code, "%s", cryptBuf)
}
