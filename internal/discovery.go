package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Milefer7/LAN-chat/app/model/broadcast"
)

// 定义Done

// CreateUdpConn 创建一个udp连接用于发送广播消息
func CreateUdpConn(data *broadcast.BroadcastMsg) (err error) {
	// 打印广播消息
	// log.Printf("查看广播消息:%+v\n", *data)

	// 创建udp连接
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 167),
		Port: 5353,
	})
	if err != nil {
		fmt.Println("udp连接错误:", err)
		return err
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭udp连接错误:", err)
		}
	}(conn)

	// 将data转化为json格式
	msgBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("转化为msgBytes错误：", err)
		return err
	}

	broadcast.LocalBroadcastMsg = *data
	// 发送心跳消息
	broadcast.StartHeartbeat <- true
	// 发送广播消息
	_, err = conn.Write(msgBytes)
	if err != nil {
		fmt.Println("发送广播错误：", err)
		return err
	}
	return nil
}

// ListenForBroadcastMessages 监听UDP广播消息,发现在线用户
func ListenForBroadcastMessages(ctx context.Context) {
	// 解析udp地址
	multicastAddr := &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 167),
		Port: 5353,
	}

	// 创建一个多播UDP连接并监听
	conn, err := net.ListenMulticastUDP("udp", nil, multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 设置读取缓冲区大小
	if err := conn.SetReadBuffer(1024); err != nil {
		log.Fatal(err)
	}

	// 输出日志，开始监听UDP广播消息
	log.Println("ListenForBroadcastMessages开启")

	// 无限循环监听udp消息
	for {
		select {
		case <-ctx.Done():
			log.Println("接收到取消信号，停止监听UDP广播")
			return
		default:
			buffer := make([]byte, 1024)
			// 设置读取超时
			err := conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			if err != nil {
				return
			}
			n, src, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if os.IsTimeout(err) {
					// 忽略超时错误，继续监听
					// log.Println("*******************监听循环次数分割线*************************")
					continue
				}
				log.Println("读取UDP消息错误：", err)
				continue
			}
			var msg broadcast.BroadcastMsg
			if err := json.Unmarshal(buffer[:n], &msg); err != nil {
				log.Println("解析UDP消息错误：", err)
				continue
			}
			log.Println("接收到UDP消息：", msg)
			// 根据消息类型处理
			switch msg.Type {
			case "join":
				creatOnlineUsers(msg, src)
			case "heartbeat":
				updateOnlineUsers(msg, src)
			}
		}
	}
}

// 新增一个在线用户
func creatOnlineUsers(msg broadcast.BroadcastMsg, src *net.UDPAddr) {
	// 加锁，保证同一个时间只有应该线程可以访问
	broadcast.OnlineUsersMutex.Lock()
	defer broadcast.OnlineUsersMutex.Unlock()

	broadcast.OnlineUsers = append(broadcast.OnlineUsers, broadcast.User{
		Host:          src.IP.String(),
		UserName:      msg.UserName,
		Fingerprint:   msg.Fingerprint,
		LastHeartbeat: time.Now(), // 设置最后一次心跳时间
	})
}

// 更新心跳消息处理逻辑
func updateOnlineUsers(msg broadcast.BroadcastMsg, src *net.UDPAddr) {
	broadcast.OnlineUsersMutex.Lock()
	defer broadcast.OnlineUsersMutex.Unlock()

	// 遍历在线用户列表，更新最后一次心跳时间
	for i, user := range broadcast.OnlineUsers {
		if user.Fingerprint == msg.Fingerprint {
			broadcast.OnlineUsers[i].LastHeartbeat = time.Now() // 更新最后一次心跳时间
			break
		}
	}
}

// RemoveStaleUsers 实现定时任务以移除离线用户
func RemoveStaleUsers(ctx context.Context) {
	log.Println("RemoveStaleUsers开启")
	ticker := time.NewTicker(5 * time.Second) // 每5秒执行一次检查
	defer ticker.Stop()                       // 确保定时器被停止，避免泄露

	for {
		select {
		case <-ctx.Done():
			// 上下文被取消，函数应该停止执行
			log.Println("接收到取消信号，停止移除离线用户")
			broadcast.OnlineUsers = nil
			return
		case <-ticker.C:
			// 定时器触发，执行移除过时用户的操作
			broadcast.OnlineUsersMutex.Lock()
			now := time.Now()
			aliveUsers := broadcast.OnlineUsers[:0] // 创建一个新的slice用于存储仍然在线的用户
			for _, user := range broadcast.OnlineUsers {
				if now.Sub(user.LastHeartbeat) <= broadcast.HeartbeatTimeout {
					aliveUsers = append(aliveUsers, user)
				}
			}
			broadcast.OnlineUsers = aliveUsers
			broadcast.OnlineUsersMutex.Unlock()
		}
	}
}

// 发送心跳消息
func SendHeartbeat(ctx context.Context) {
	log.Println("接收到信号，开启发送心跳")
	// 创建udp连接
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 167),
		Port: 5353,
	})
	if err != nil {
		log.Printf("udp连接错误：%v", err)
		return
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭udp连接错误:", err)
		}
	}(conn)
	ticker := time.NewTicker(25 * time.Second) // 每25秒发送一次心跳
	broadcast.LocalBroadcastMsg.Type = "heartbeat"
	for {
		select {
		case <-ticker.C:
			// 将data转化为json格式
			msgBytes, err := json.Marshal(broadcast.LocalBroadcastMsg)
			// log.Println("data:", data)
			if err != nil {
				log.Printf("转化心跳消息为msgBytes错误：%v", err)
				continue
			}
			// 发送心跳消息
			if _, err := conn.Write(msgBytes); err != nil {
				log.Printf("发送心跳错误：%v", err)
				continue
			}
			log.Println("已发送心跳消息：", broadcast.LocalBroadcastMsg)
			log.Println("在线用户：", broadcast.OnlineUsers)
		case <-ctx.Done():
			log.Println("停止发送心跳消息")
			// 结束循环
			return
		}
	}
}
