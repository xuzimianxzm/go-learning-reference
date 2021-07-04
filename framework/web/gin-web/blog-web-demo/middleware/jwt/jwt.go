package jwt

import (
	"github.com/xuzimian/blog-web-demo/pkg/error"
	"github.com/xuzimian/blog-web-demo/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = error.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = error.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = error.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = error.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != error.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  error.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
