package middleware

import (
	"FoodBackend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"FoodBackend/pkg/e"
	"FoodBackend/pkg/util"
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
		UserId = models.GetUserIdByToken(token)
		c.Next()
	}
}
