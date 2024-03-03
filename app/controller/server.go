package controller

// import "C"
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Milefer7/LAN-chat/app/model/broadcast"
	"github.com/Milefer7/LAN-chat/internal"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StartServer 开启服务器
func StartServer(r *gin.Engine) {
	log.Println(fmt.Sprintf("Server is running at %s", "localhost:8080"))
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

// Ws 处理websocket连接
func Ws(c *gin.Context, ctx context.Context) {
	// 调用升级器 upgrader 将HTTP请求升级到WebSocket连接将HTTP请求升级到WebSocket连接
	ws, err := broadcast.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// 关闭连接
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Println("Close error:", err)
		}
	}(ws)

	// 检查客户端是否在本地, 如果是则记录
	remoteAddr := ws.RemoteAddr().String()
	host, _, err := net.SplitHostPort(remoteAddr)
	if host == "127.0.0.1" {
		broadcast.LocalClient = ws
		log.Println("本地客户端已连接")
	}
	log.Println("客户端连接成功: ", remoteAddr)
	// 无限循环处理客户端发来的信息
	internal.HandleGet(ws, ctx)
}

// Broadcast 广播 join消息+返回在线用户
func Broadcast(c *gin.Context) {
	// 解析body
	flag := true
	var data broadcast.BroadcastMsg
	if err := c.BindJSON(&data); err != nil {
		log.Println("发送广播 BindJSON error:", err)
		flag = false
	}
	// 发送广播消息
	err := internal.CreateUdpConn(&data)
	if err != nil {
		log.Println("发送广播 CreateUdpConn error:", err)
		flag = false
	}
	// 返回在线用户
	if flag {
		c.JSON(200, gin.H{
			"code":    1,
			"message": "发送广播成功",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "调用广播接口失败，详细信息请查看日志",
		})
	}
}

// Leave 定义用户离开接口
func Leave(c *gin.Context, cancel context.CancelFunc) {
	var req broadcast.LeaveRequest
	if err := c.BindJSON(&req); err != nil {
		log.Println("离开接口 BindJSON error:", err)
		return
	}
	if req.Leave {
		broadcast.LocalClient = nil
		// 通知goroutine退出
		defer cancel()
	}
	c.JSON(200, gin.H{
		"code":    "1",
		"message": "调用离开接口成功",
	})
}

func GetOnlineUsers(c *gin.Context) {
	// 返回在线用户
	c.JSON(200, gin.H{
		"code":        1,
		"message":     "获取在线用户成功",
		"onlineUsers": broadcast.OnlineUsers,
	})
}

func HttpSend(c *gin.Context) {
	// 结构体转转为json格式
	data, err := json.Marshal(broadcast.LocalBroadcastMsg)
	if err != nil {
		log.Println("结构体转json格式错误:", err)
		return
	}

	// 获取本机wlan的IPv4地址
	ip, err := internal.GetOutBoundIP()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("本机ip是：", ip)
	// 获取本机wlan的子网
	parts := strings.Split(ip, ".")
	baseIP := strings.Join(parts[:3], ".") + "."
	// fmt.Println(baseIP)

	var wg sync.WaitGroup

	// 定义一个切片来存储发现的设备IP地址
	var discoveredIPs []string

	// 遍历整个IP范围
	for i := 0; i <= 255; i++ {
		ip := fmt.Sprintf("%s%d", baseIP, i)
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			url := fmt.Sprintf("http://%s:8080/httpReceive", ip)

			// 发送POST请求
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
			if err != nil {
				// fmt.Printf("Failed to send request to %s: %v\n", ip, err)
				return
			}
			defer resp.Body.Close()

			// 如果连接成功，写入切片
			if resp.StatusCode == http.StatusOK {
				discoveredIPs = append(discoveredIPs, ip)
			}
		}(ip)
	}

	wg.Wait()

	// 打印发现的IP地址
	fmt.Println("Discovered IP addresses:")
	for _, ip := range discoveredIPs {
		fmt.Println(ip)
	}
	c.JSON(200, gin.H{
		"code":    1,
		"message": "发送http接口成功",
	})
}

func HttpReceive(c *gin.Context) {
	// 接收json
	var msg broadcast.BroadcastMsg
	if err := c.BindJSON(&msg); err != nil {
		log.Println("接收http接口 BindJSON error:", err)
		return
	}
	// 加锁，保证同一个时间只有应该线程可以访问
	broadcast.OnlineUsersMutex.Lock()
	defer broadcast.OnlineUsersMutex.Unlock()

	// log.Println("Host:", c.ClientIP())
	// 遍历在线用户，看是否为新用户
	for i, user := range broadcast.OnlineUsers {
		if user.Fingerprint == msg.Fingerprint {
			broadcast.OnlineUsers[i].LastHeartbeat = time.Now() // 设置最后一次心跳时间
			c.JSON(200, gin.H{
				"code":    1,
				"message": "接收http接口成功",
			})
			return
		}
	}
	// 如果是新用户，加入在线用户列表
	broadcast.OnlineUsers = append(broadcast.OnlineUsers, broadcast.User{
		Host:          c.ClientIP(),
		UserName:      msg.UserName,
		Fingerprint:   msg.Fingerprint,
		LastHeartbeat: time.Now(), // 设置最后一次心跳时间
	})
	c.JSON(200, gin.H{
		"code":    1,
		"message": "接收http接口成功",
	})
}
