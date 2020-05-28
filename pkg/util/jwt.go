package util

import (
	"FoodBackend/pkg/setting"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	//Password string `json:"password"`
	jwt.StandardClaims
}

func GetToken(c *gin.Context) string {
	header := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(header) == 2 {
		return header[1]
	} else {
		return ""
	}
}

func GenerateToken(username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ieo",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //crypto.Hash
	token, err := tokenClaims.SignedString(jwtSecret)                //生成签名字符串

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	//解析鉴权
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		//验证基于时间的声明exp, iat, nbf
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			//如果没有任何声明在令牌中，仍然会被认为是有效的，并且对于时区偏差没有计算方法!
			return claims, nil
		}
	}

	return nil, err
}
