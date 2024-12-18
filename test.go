//go:build test
// +build test

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"ai.zhycit.com/socket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func startServer() *socket.PostOffice {
	// 创建 PostOffice 实例
	po, err := socket.NewPostOffice(100, "") // 最大连接数100，不使用schema验证
	if err != nil {
		log.Fatal("Failed to create PostOffice:", err)
	}

	// 设置WebSocket路由
	http.HandleFunc("/ws", po.HandleConnection)

	// 启动HTTP服务器
	go func() {
		log.Printf("Starting server on %s", *addr)
		if err := http.ListenAndServe(*addr, nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	return po
}

func runClient(clientID string, targetID string, done chan struct{}) {
	// 构建WebSocket URL
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	queryParams := u.Query()
	queryParams.Set("clientID", clientID)
	u.RawQuery = queryParams.Encode()

	// 连接WebSocket服务器
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// 创建接收消息的goroutine
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("Client %s received: %s", clientID, message)
		}
	}()

	// 发送测试消息
	message := fmt.Sprintf(`{"to":"%s","content":"Message from %s","timestamp":"%s"}`, 
		targetID, 
		clientID,
		time.Now().Format(time.RFC3339))

	// 连续发送5条消息
	for i := 0; i < 5; i++ {
		err := c.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("write:", err)
			return
		}
		log.Printf("Client %s sent message %d", clientID, i+1)
		time.Sleep(time.Millisecond * 100) // 短暂延迟，模拟连续发送
	}

	time.Sleep(time.Second * 2) // 等待接收响应
	done <- struct{}{}
}

func main() {
	flag.Parse()

	// 启动服务器
	startServer()
	time.Sleep(time.Second) // 等待服务器启动

	// 创建用于等待客户端完成的通道
	done := make(chan struct{})

	// 启动两个测试客户端
	go runClient("client1", "client2", done)
	go runClient("client2", "client1", done)

	// 等待两个客户端完成
	<-done
	<-done

	// 等待中断信号以优雅地关闭服务器
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Println("Test completed")
}
