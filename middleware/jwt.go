package middleware

import (
	"FoodBackend/models"
	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var UserId int
var AuthToken string

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		code = e.SUCCESS
		token := util.GetToken(c)
		if token == "" {
			code = e.UNAUTHORIZED
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		UserId = models.GetUserIdByToken(token)
		if UserId == 0 {
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
			c.Abort()
			return
		}

		AuthToken = token
		c.Next()
	}
}
