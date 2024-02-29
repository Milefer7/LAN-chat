package main

import (
	"context"
	"log"

	"github.com/Milefer7/LAN-chat/app/controller"
	"github.com/Milefer7/LAN-chat/app/model/broadcast"
	"github.com/Milefer7/LAN-chat/internal"
	"github.com/Milefer7/LAN-chat/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个context，用于取消goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// 开启设备发现 获取在线用户
	go internal.ListenForBroadcastMessages(ctx)
	// 清除超时用户
	go internal.RemoveStaleUsers(ctx)
	// 发送心跳
	go func() {
		log.Printf("等待信号发送心跳")
		// 每25秒发一次心跳
		// RemoveStaleUsers是每5秒清除一次超时用户
		// 超时时间为30秒
		for {
			select {
			case <-broadcast.StartHeartbeat:
				go internal.SendHeartbeat(ctx)
			}
		}
	}()
	// 初始化路由
	r := gin.Default()
	router.InitRouter(r, cancel, ctx)
	// 启动服务端
	controller.StartServer(r)
}
