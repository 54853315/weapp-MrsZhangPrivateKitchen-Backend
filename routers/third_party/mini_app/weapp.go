package mini_app

import (
	"FoodBackend/pkg/cache"
	"FoodBackend/pkg/setting"
	"fmt"
	"github.com/medivhzhan/weapp/v2"
	"log"
	"time"
)

var (
	err                       error
	appId, appSecret, appCode string
)

const (
	gcTime time.Duration = 7200
)

func init() {
	sec, err := setting.Cfg.GetSection("wechat")
	if err != nil {
		log.Fatal(2, "Fail to get section 'wechat': %v", err)
	}

	appId = sec.Key("APP_ID").String()
	appSecret = sec.Key("APP_SECRET").String()
}

func getAccessToken() string {
	c := cache.NewCache(cache.DefaultExpiration, gcTime)
	appCodeCacheName := "wechat_app_code"

	if cacheValue, cacheExists := c.Get(appCodeCacheName); !cacheExists {
		fmt.Println("数据不存在，去获取")
		res, err := weapp.GetAccessToken(appId, appSecret)
		if err != nil {
			log.Println(err)
		}
		c.Set(appCodeCacheName, res.AccessToken, time.Duration(res.ExpiresIn))
		return res.AccessToken
	} else {
		fmt.Println("数据已存在")
		return cacheValue.(string)
	}
}
