package main

import (
	"FoodBackend/routers"
)

func main() {
	r := routers.InitRouter()
	r.Static("/uploads", "./uploads")
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
