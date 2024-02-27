package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	model "github.com/Milefer7/LAN-chat/app/model/broadcast"
	"github.com/Milefer7/LAN-chat/app/model/communication"
)

func MockBroadcastMsg(i int) model.BroadcastMsg {
	return model.BroadcastMsg{
		UserName:    fmt.Sprintf("TestUser%d", i),
		DeviceType:  fmt.Sprintf("TestDevice%d", i),
		Fingerprint: fmt.Sprintf("Fingerprint%d", i),
		Port:        52017 + i, // 确保端口在合法范围内
	}
}

func TestWs(t *testing.T) {
	// 定义一个在线用户列表
	// 定义一个在线用户列表
	onlineUsers := []model.User{
		{
			UserName:      "User1",
			Fingerprint:   "Fingerprint1",
			Host:          "Host1",
			LastHeartbeat: time.Now(),
		},
		{
			UserName:      "User2",
			Fingerprint:   "Fingerprint2",
			Host:          "Host2",
			LastHeartbeat: time.Now(),
		},
		{
			UserName:      "User3",
			Fingerprint:   "Fingerprint3",
			Host:          "Host3",
			LastHeartbeat: time.Now(),
		},
		{
			UserName:      "User4",
			Fingerprint:   "Fingerprint4",
			Host:          "Host4",
			LastHeartbeat: time.Now(),
		},
		{
			UserName:      "User5",
			Fingerprint:   "Fingerprint5",
			Host:          "Host5",
			LastHeartbeat: time.Now(),
		},
	}
	requestMsg := communication.RequestMessage{
		Data:        communication.Data{},
		Content:     []communication.ContentReq{},
		OnlineUsers: communication.OnlineUsers{onlineUsers},
	}
	users, _ := json.Marshal(requestMsg)
	// 打印看看
	log.Println("OnlineUsers:", string(users))
}

// 测试广播消息的发送和接收（看看防火墙会不会拦截）
func TestBroadcastAndReceive(t *testing.T) {
	// 创建接收端，加入多播组
	recvAddr := &net.UDPAddr{
		IP:   net.ParseIP("224.0.0.167"),
		Port: 5353,
	}
	recvConn, err := net.ListenMulticastUDP("udp", nil, recvAddr)
	if err != nil {
		t.Errorf("ListenMulticastUDP error: %v", err)
		return
	}
	defer recvConn.Close()

	// 设置接收端读取超时
	err = recvConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		t.Errorf("SetReadDeadline error: %v", err)
		return
	}

	// 等待一小段时间确保接收端准备就绪
	log.Println(1)
	time.Sleep(time.Second * 3)
	log.Println(2)

	// 创建发送端连接
	conn, err := net.Dial("udp", "224.0.0.167:5353")
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer conn.Close()

	// 创建并发送广播消息
	msg := model.BroadcastMsg{
		UserName:    "TestUser",
		DeviceType:  "TestDevice",
		Fingerprint: "Fingerprint",
		Port:        5353,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	_, err = conn.Write(msgBytes)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// 读取接收端的数据
	buffer := make([]byte, 1024)
	n, _, err := recvConn.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("ReadFromUDP error: %v", err)
	}

	// 解析接收到的消息
	var recvMsg model.BroadcastMsg
	err = json.Unmarshal(buffer[:n], &recvMsg)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// 检查接收到的消息是否与发送的消息匹配
	if recvMsg.UserName != msg.UserName || recvMsg.DeviceType != msg.DeviceType || recvMsg.Fingerprint != msg.Fingerprint || recvMsg.Port != msg.Port {
		t.Fatalf("Received message does not match sent message")
	}
}

const BUF_SIZE int = 8192

func TestBroadcast(t *testing.T) {
	// sender
	conn, err := net.Dial("udp", "224.1.1.2:9190")
	if err != nil {
		checkError(err)
	}

	go func() {
		for {
			conn.Write([]byte("hello, world!"))
			time.Sleep(1 * time.Second)
		}
	}()

	// receiver
	gaddr, _ := net.ResolveUDPAddr("udp", "224.1.1.2:9190")
	listener, err := net.ListenMulticastUDP("udp", nil, gaddr)
	if err != nil {
		checkError(err)
	}

	// listener.SetReadBuffer(maxDatagramSize)

	message := make([]byte, BUF_SIZE)

	for {
		n, src, _ := listener.ReadFromUDP(message)
		fmt.Println(src, ": ", string(message[:n]))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func TestSender(t *testing.T) {
	udpAddress := "224.1.1.2:9190"
	message := "hello, world!"
	conn, err := net.Dial("udp", udpAddress)
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}

	done := time.After(10 * time.Second)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				conn.Write([]byte(message))
				log.Println("Sent message:" + message)
			}
		}
	}()
}

func TestReceiver(t *testing.T) {
	udpAddress := "224.1.1.2:9190"
	gaddr, _ := net.ResolveUDPAddr("udp", udpAddress)
	listener, err := net.ListenMulticastUDP("udp", nil, gaddr)
	if err != nil {
		t.Fatalf("ListenMulticastUDP error: %v", err)
	}

	message := make([]byte, BUF_SIZE)
	done := time.After(10 * time.Second)

	for {
		select {
		case <-done:
			return
		default:
			n, src, _ := listener.ReadFromUDP(message)
			fmt.Println(src, ": ", string(message[:n]))
		}
	}
}
