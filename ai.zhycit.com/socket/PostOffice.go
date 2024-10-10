package socket

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

// PostOffice 结构体用于管理 WebSocket 连接
type PostOffice struct {
	upgrader        websocket.Upgrader
	clients         sync.Map // key: channelId, value: *websocket.Conn
	connectionCount int32
	maxConnections  int32
}

// NewPostOffice 创建一个新的 PostOffice 实例
func NewPostOffice(maxConnections int) *PostOffice {
	return &PostOffice{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源
			},
		},
		maxConnections: int32(maxConnections),
	}
}

// HandleConnection 处理 WebSocket 连接的升级和消息接收
func (po *PostOffice) HandleConnection(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&po.connectionCount) >= po.maxConnections {
		http.Error(w, "Connection limit reached", http.StatusServiceUnavailable)
		return
	}

	// 从查询参数中获取客户端 ID
	clientID := r.URL.Query().Get("clientID")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	conn, err := po.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error during connection upgrade:", err)
		return
	}

	atomic.AddInt32(&po.connectionCount, 1)
	defer atomic.AddInt32(&po.connectionCount, -1)

	po.clients.Store(clientID, conn)
	fmt.Printf("Client connected: %s (Total: %d/%d)\n", clientID, atomic.LoadInt32(&po.connectionCount), po.maxConnections)

	defer func() {
		conn.Close()
		po.clients.Delete(clientID)
		fmt.Printf("Client disconnected: %s (Total: %d/%d)\n", clientID, atomic.LoadInt32(&po.connectionCount), po.maxConnections)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}
		po.messageDelivery(message)
	}
}

// messageDelivery 处理消息分发
func (po *PostOffice) messageDelivery(msg []byte) {
	msgJson := gjson.ParseBytes(msg)
	fmt.Printf("Received message: %s\n", string(msg))

	to := msgJson.Get("to")
	if to.Exists() {
		targets := to.Array()
		fmt.Printf("Targets: %v\n", targets)
		for _, target := range targets {
			po.delivery(target.String(), msg)
		}
	} else {
		fmt.Println("No 'to' field in the message")
	}
}

// delivery 向指定的目标发送消息
func (po *PostOffice) delivery(targetChannelId string, msg []byte) {
	if connInterface, ok := po.clients.Load(targetChannelId); ok {
		if conn, ok := connInterface.(*websocket.Conn); ok {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Printf("Error sending message to %s: %v\n", targetChannelId, err)
				po.clients.Delete(targetChannelId)
			} else {
				fmt.Printf("Message sent to %s\n", targetChannelId)
			}
		} else {
			fmt.Printf("Invalid connection type for %s\n", targetChannelId)
		}
	} else {
		fmt.Printf("Client %s not found\n", targetChannelId)
	}
}

// GetConnectionCount 返回当前连接数
func (po *PostOffice) GetConnectionCount() int {
	return int(atomic.LoadInt32(&po.connectionCount))
}

// GetMaxConnections 返回最大连接数
func (po *PostOffice) GetMaxConnections() int {
	return int(po.maxConnections)
}
