package internal

import (
	"context"
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
