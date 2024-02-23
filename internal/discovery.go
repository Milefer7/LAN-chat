package internal

import (
	"context"
	"encoding/json"
	"fmt"
	model "github.com/Milefer7/LAN-chat/app/model/broadcast"
	"log"
	"net"
	"os"
	"time"
)

// 定义Done

// CreateUdpConn 创建一个udp连接用于发送广播消息
func CreateUdpConn(data model.BroadcastMsg, ctx context.Context) (err error) {
	// 创建udp连接
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 167),
		Port: 52017,
	})
	if err != nil {
		fmt.Println("udp连接错误:", err)
		return err
	}
	defer conn.Close()

	// 将data转化为json格式
	msgBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("转化为msgBytes错误：", err)
		return err
	}
	// 发送广播消息
	_, err = conn.Write(msgBytes)
	if err != nil {
		fmt.Println("发送广播错误：", err)
		return err
	}
	// 心跳消息
	go sendHeartbeat(conn, data, ctx)
	return nil
}

// 监听UDP广播消息,发现在线用户
func ListenForBroadcastMessages(ctx context.Context) {
	// 解析udp地址
	addr, err := net.ResolveUDPAddr("udp", ":52017")
	if err != nil {
		log.Fatal(err)
	}
	// 创建udp连接
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 设置为非阻塞模式
	if err := conn.SetReadBuffer(1024); err != nil {
		log.Fatal(err)
	}

	// 无限循环监听udp消息
	for {
		select {
		case <-ctx.Done():
			log.Println("接收到取消信号，停止监听UDP广播消息,发现在线用户")
			return
		default:
			buffer := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(1 * time.Second)) // 设置短暂的读取超时
			n, src, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if os.IsTimeout(err) {
					// 忽略超时错误，继续监听
					continue
				}
				log.Println("读取UDP消息错误：", err)
				continue
			}

			var msg model.BroadcastMsg
			if err := json.Unmarshal(buffer[:n], &msg); err != nil {
				log.Println("解析UDP消息错误：", err)
				continue
			}

			switch msg.Type {
			case "join":
				// 新增一个在线用户
				creatOnlineUsers(msg, src)
			case "heartbeat":
				updateOnlineUsers(msg, src)
			}
		}
	}
}

// 新增一个在线用户
func creatOnlineUsers(msg model.BroadcastMsg, src *net.UDPAddr) {
	// 加锁，保证线程安全
	model.OnlineUsersMutex.Lock()
	defer model.OnlineUsersMutex.Unlock()

	model.OnlineUsers = append(model.OnlineUsers, model.User{
		Host:          src.IP.String(),
		UserName:      msg.UserName,
		Fingerprint:   msg.Fingerprint,
		LastHeartbeat: time.Now(), // 设置最后一次心跳时间
	})
}

// 更新心跳消息处理逻辑
func updateOnlineUsers(msg model.BroadcastMsg, src *net.UDPAddr) {
	model.OnlineUsersMutex.Lock()
	defer model.OnlineUsersMutex.Unlock()

	exists := false
	for i, user := range model.OnlineUsers {
		if user.Fingerprint == msg.Fingerprint {
			exists = true
			model.OnlineUsers[i].LastHeartbeat = time.Now() // 更新最后一次心跳时间
			break
		}
	}
	if !exists {
		model.OnlineUsers = append(model.OnlineUsers, model.User{
			Host:          src.IP.String(),
			UserName:      msg.UserName,
			Fingerprint:   msg.Fingerprint,
			LastHeartbeat: time.Now(), // 设置最后一次心跳时间
		})
	}
}

// 实现定时任务以移除离线用户
func RemoveStaleUsers(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second) // 每5秒执行一次检查
	defer ticker.Stop()                       // 确保定时器被停止，避免泄露

	for {
		select {
		case <-ctx.Done():
			// 上下文被取消，函数应该停止执行
			log.Println("接收到取消信号，停止移除离线用户")
			return
		case <-ticker.C:
			// 定时器触发，执行移除过时用户的操作
			model.OnlineUsersMutex.Lock()
			now := time.Now()
			aliveUsers := model.OnlineUsers[:0] // 创建一个新的slice用于存储仍然在线的用户
			for _, user := range model.OnlineUsers {
				if now.Sub(user.LastHeartbeat) <= model.HeartbeatTimeout {
					aliveUsers = append(aliveUsers, user)
				}
			}
			model.OnlineUsers = aliveUsers
			model.OnlineUsersMutex.Unlock()
		}
	}
}

// 发送心跳消息
func sendHeartbeat(conn *net.UDPConn, data model.BroadcastMsg, ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // 每30秒发送一次心跳
	data.Type = "heartbeat"
	for {
		select {
		case <-ticker.C:
			// 将data转化为json格式
			msgBytes, err := json.Marshal(data)
			if err != nil {
				log.Printf("转化心跳消息为msgBytes错误：%v", err)
				continue
			}
			// 发送心跳消息
			if _, err := conn.Write(msgBytes); err != nil {
				log.Printf("发送心跳错误：%v", err)
				continue
			}
		case <-ctx.Done():
			log.Println("停止发送心跳消息")
			// 结束循环
			return
		}
	}
}
