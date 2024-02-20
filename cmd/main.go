package main

import (
	"github.com/Milefer7/LAN-chat/app/controller"
	"github.com/Milefer7/LAN-chat/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化路由
	r := gin.Default()
	router.InitRouter(r)
	// 启动服务端
	controller.StartServer(r)
}
