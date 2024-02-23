package controller

import (
	"context"
	"fmt"
	"github.com/Milefer7/LAN-chat/app/model/broadcast"
	"github.com/Milefer7/LAN-chat/internal"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net"
)

// StartServer 开启服务器
func StartServer(r *gin.Engine) {
	log.Println(fmt.Sprintf("Server is running at %s", "localhost:80"))
	if err := r.Run(":80"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

// Ws 处理websocket连接
func Ws(c *gin.Context, ctx context.Context) {
	// 调用升级器 upgrader 将HTTP请求升级到WebSocket连接将HTTP请求升级到WebSocket连接
	ws, err := internal.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// 关闭连接
	defer func(ws *websocket.Conn) {
		// Send a close frame
		message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
		if err := ws.WriteMessage(websocket.CloseMessage, message); err != nil {
			log.Println("Write close message error:", err)
		}
		err := ws.Close()
		if err != nil {
			log.Println("Close error:", err)
		}
	}(ws)

	// 检查客户端是否在本地, 如果是则记录
	remoteAddr := ws.RemoteAddr().String()
	host, _, err := net.SplitHostPort(remoteAddr)
	if host == "127.0.0.1" {
		internal.LocalClient = ws
		// 每隔3s向本地服务端返回在线用户
		go internal.SendOnlineUsers(internal.LocalClient, ctx)
	}

	// 无限循环处理客户端发来的信息
	go internal.HandleGet(ws, ctx)

}

// Broadcast 广播 join消息+返回在线用户
func Broadcast(c *gin.Context, ctx context.Context) {
	// 解析body
	flag := true
	var data broadcast.BroadcastMsg
	if err := c.BindJSON(&data); err != nil {
		log.Println("BindJSON error:", err)
		flag = false
		return
	}

	//// 创建一个UDP连接用于发送广播消息
	//err := internal.CreateUdpConn(data)
	//if err != nil {
	//	log.Println("CreateUdpConn error:", err)
	//	return
	//}
	//// 启动一个goroutine监听UDP广播消息
	//users := make([]model.User, 0)
	//done := make(chan bool)
	//go func() {
	//	// 监听广播消息
	//	pc, err := net.ListenPacket("udp", ":52017")
	//	if err != nil {
	//		log.Println(err)
	//		done <- true
	//		return
	//	}
	//	defer pc.Close()
	//
	//	// 设置超时时间
	//	err = pc.SetDeadline(time.Now().Add(10 * time.Second))
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//
	//	for {
	//		buffer := make([]byte, 1024)
	//		_, addr, err := pc.ReadFrom(buffer)
	//		if err != nil {
	//			break // 退出循环
	//		}
	//
	//		var msg model.BroadcastMsg
	//		if err := json.Unmarshal(buffer, &msg); err != nil {
	//			log.Println(err)
	//			continue // 忽略错误的消息
	//		}
	//
	//		host, _, _ := net.SplitHostPort(addr.String())
	//		user := model.User{
	//			Host:        host,
	//			UserName:    msg.UserName,
	//			Fingerprint: msg.Fingerprint,
	//		}
	//		users = append(users, user)
	//	}
	//
	//	done <- true
	//}()
	//
	//<-done // 等待监听结束

	// 发送广播消息
	err := internal.CreateUdpConn(data, ctx)
	if err != nil {
		log.Println("CreateUdpConn error:", err)
		flag = false
		return
	}

	//
	if flag {
		c.JSON(200, gin.H{
			"code":    "1",
			"message": "调用广播接口成功",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    "0",
			"message": "调用广播接口失败，详细信息请查看日志",
		})
	}
}

// 定义用户离开接口
func Leave(c *gin.Context, cancel context.CancelFunc) {
	type Request struct {
		Leave   bool   `json:"leave"`
		Message string `json:"message"`
	}
	var req Request
	if err := c.BindJSON(&req); err != nil {
		log.Println("BindJSON error:", err)
		return
	}
	if req.Leave {
		internal.LocalClient = nil
		// 通知goroutine退出
		cancel()
	}
	c.JSON(200, gin.H{
		"code":    "1",
		"message": "调用离开接口成功",
	})
}
