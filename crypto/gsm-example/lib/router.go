package lib

import (
	"github.com/gin-gonic/gin"
	"github.com/sunshibao/go-utils/crypto/gsm-example/lib/core"
	"github.com/sunshibao/go-utils/crypto/gsm-example/lib/middleware"
)

func Register(g *gin.Engine) *gin.Engine {

	g.POST("/hello", core.Handler(middleware.CryptoCheck()), core.Handler(func(ctx *core.Context) {
		ctx.Mix(200, gin.H{})
	}))
	return g
}
