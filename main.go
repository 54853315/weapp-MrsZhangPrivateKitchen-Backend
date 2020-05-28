package main

import (
	"FoodBackend/routers"
)

func main() {
	r := routers.InitRouter()
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
