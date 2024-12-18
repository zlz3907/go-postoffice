package socket

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
)

// PostOffice 结构体用于管理 WebSocket 连接
type PostOffice struct {
	upgrader        websocket.Upgrader
	clients         sync.Map // key: channelId, value: *websocket.Conn
	connectionCount int32
	maxConnections  int32
	schema          *gojsonschema.Schema
}

// NewPostOffice 创建一个新的 PostOffice 实例
func NewPostOffice(maxConnections int, schemaPath string) (*PostOffice, error) {
	po := &PostOffice{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源
			},
		},
		maxConnections: int32(maxConnections),
	}

	if schemaPath != "" {
		schemaBytes, err := os.ReadFile(schemaPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read schema file: %v", err)
		}

		schemaLoader := gojsonschema.NewStringLoader(string(schemaBytes))
		schema, err := gojsonschema.NewSchema(schemaLoader)
		if err != nil {
			return nil, fmt.Errorf("invalid schema: %v", err)
		}
		po.schema = schema
	}

	return po, nil
}

// validateConnection 验证连接的有效性
func (po *PostOffice) validateConnection(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	// if !strings.HasPrefix(authHeader, "Bearer ") {
	// 	return false
	// }
	token := strings.TrimPrefix(authHeader, "Bearer ")
	// TODO: 在这里添加更具体的token验证逻辑
	return true || token != ""
}

// HandleConnection 处理 WebSocket 连接的升级和消息接收
func (po *PostOffice) HandleConnection(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&po.connectionCount) >= po.maxConnections {
		http.Error(w, "Connection limit reached", http.StatusServiceUnavailable)
		return
	}

	// 验证连接
	if !po.validateConnection(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	// 创建一个带缓冲的通道用于消息处理
	msgChan := make(chan []byte, 100)
	
	// 启动一个 goroutine 专门处理消息发送
	go func() {
		for message := range msgChan {
			po.messageDelivery(message)
		}
	}()

	// 主循环读取消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			close(msgChan)
			return
		}
		// 将消息发送到通道，非阻塞方式
		select {
		case msgChan <- message:
			// 消息成功加入队列
		default:
			// 队列已满，可以选择丢弃消息或者记录日志
			fmt.Println("Message queue is full, message dropped")
		}
	}
}

// messageDelivery 处理消息分发
func (po *PostOffice) messageDelivery(msg []byte) {
	// 如果存在 schema，则进行验证
	if po.schema != nil {
		result, err := po.schema.Validate(gojsonschema.NewBytesLoader(msg))
		if err != nil {
			fmt.Printf("Error validating message: %v\n", err)
			return
		}
		if !result.Valid() {
			fmt.Println("Invalid message:")
			for _, desc := range result.Errors() {
				fmt.Printf("- %s\n", desc)
			}
			return
		}
	}

	// 消息处理逻辑
	msgJson := gjson.ParseBytes(msg)
	fmt.Printf("Received message: %s\n", string(msg))

	to := msgJson.Get("to")
	if to.Exists() {
		if to.IsArray() {
			targets := to.Array()
			fmt.Printf("Targets: %v\n", targets)
			for _, target := range targets {
				po.delivery(target.String(), msg)
			}
		} else {
			target := to.String()
			fmt.Printf("Single target: %s\n", target)
			po.delivery(target, msg)
		}
	} else {
		fmt.Println("No 'to' field in the message")
	}
}

// delivery 向指定的目标发送消息
func (po *PostOffice) delivery(targetChannelId string, msg []byte) {
	if connInterface, ok := po.clients.Load(targetChannelId); ok {
		if conn, ok := connInterface.(*websocket.Conn); ok {
			// 使用 WriteControl 发送一个 ping 来检查连接状态
			err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
			if err != nil {
				fmt.Printf("Connection check failed for %s: %v\n", targetChannelId, err)
				po.clients.Delete(targetChannelId)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, msg)
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
