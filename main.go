package main

import (
	"FoodBackend/pkg/setting"
	"FoodBackend/routers"
	"strconv"
)

func main() {
	r := routers.InitRouter()
	r.Static("/uploads", "./uploads")
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	port := setting.HTTPPort
	r.Run(":" + strconv.Itoa(port)) // 监听并在 0.0.0.0:8080 上启动服务
}
