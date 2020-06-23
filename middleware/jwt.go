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
		util.Log.Debug("token", token)
		if token == "" {
			code = e.UNAUTHORIZED
		} else {
			claims, err := util.ParseToken(token)
			util.Log.Debug("err", err)
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
