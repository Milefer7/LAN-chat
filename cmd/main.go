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
	//
	go func() {
		log.Println("StartHeartbeat开启")
		for {
			select {
			case <-broadcast.StartHeartbeat:
				log.Println("StartHeartbeat接收到消息")
				go internal.SendHeartbeat(broadcast.LocalBroadcastMsg, ctx)
			}
		}
	}()
	// 初始化路由
	r := gin.Default()
	router.InitRouter(r, cancel, ctx)
	// 启动服务端
	controller.StartServer(r)
}

//
//const BUF_SIZE int = 8192

//func main() {
//	log.Println("我是接收端")
//	// WLAN
//	iface, err := net.InterfaceByName("WLAN") // 举例使用"eth0"，实际上根据你的网络接口来
//	if err != nil {
//		log.Fatal(err)
//	}
//	// group addr
//	gaddr, _ := net.ResolveUDPAddr("udp", "224.1.1.2:9190")
//	listener, err := net.ListenMulticastUDP("udp", iface, gaddr)
//	if err != nil {
//		checkError(err)
//	}
//
//	// listener.SetReadBuffer(maxDatagramSize)
//
//	message := make([]byte, BUF_SIZE)
//	log.Println(1)
//	for {
//		log.Println(2)
//		n, src, _ := listener.ReadFromUDP(message)
//		log.Println(src, ": ", string(message[:n]))
//	}
//}
//
//func checkError(err error) {
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
//		os.Exit(1)
//	}
//}
