package internal

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestListenForBroadcastMessages(t *testing.T) {
	// 创建一个context，用于取消goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 在另一个goroutine中运行ListenForBroadcastMessages函数
	go ListenForBroadcastMessages(ctx)

	// 等待一段时间，确保ListenForBroadcastMessages函数已经开始运行
	time.Sleep(time.Second)

	// 创建一个UDP连接，用于发送数据到多播地址
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 167),
		Port: 52017,
	})
	if err != nil {
		t.Fatalf("创建UDP连接失败: %v", err)
	}
	defer conn.Close()
}

func TestCreateUdpConn(t *testing.T) {
	// 创建udp连接
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(224, 0, 0, 167),
		Port: 5353,
	})
	if err != nil {
		t.Fatalf("udp连接错误: %v", err)
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭udp连接错误:", err)
		}
	}(conn)

	// 发送消息
	_, err = conn.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("发送广播错误: %v", err)
	}
}
