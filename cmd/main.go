package main

import (
	"context"
	"github.com/Milefer7/LAN-chat/app/controller"
	"github.com/Milefer7/LAN-chat/internal"
	"github.com/Milefer7/LAN-chat/router"
	"github.com/gin-gonic/gin"
)

func main() {
	//
	ctx, cancel := context.WithCancel(context.Background())
	// 初始化路由
	r := gin.Default()
	router.InitRouter(r, cancel, ctx)
	// 启动服务端
	controller.StartServer(r)

	// 开启设备发现 获取在线用户
	go internal.ListenForBroadcastMessages(ctx)
	// 清除超时用户
	go internal.RemoveStaleUsers(ctx)
}
