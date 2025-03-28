package socket

import (
	"fmt" // 添加这一行
	"net/http"
	"os"
	"strconv"
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
	clients         sync.Map // key: channelId, value: *Client
	connectionCount int32
	maxConnections  int32
	schema          *gojsonschema.Schema
	clientTimes     sync.Map // key: clientID, value: int64 (创建时间)
}

// Client 结构体用于存储客户端连接及其互斥锁
type Client struct {
	conn      *websocket.Conn
	mutex     sync.Mutex
	createdAt int64  // 记录连接创建时间
	ip        string // 客户端IP地址
	userAgent string // 客户端User-Agent
	os        string // 操作系统信息
	browser   string // 浏览器信息
	device    string // 设备信息
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

// 添加一个全局函数来格式化时间
func currentTime() string {
	return time.Now().Format("2006-01-02 15:04:05") // 获取当前系统时间并格式化
}

// 添加一个函数来获取客户端真实IP
func getClientIP(r *http.Request) string {
	// 按优先级获取IP地址
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip != "" {
			// 如果包含多个IP，取第一个
			ip = strings.Split(ip, ",")[0]
		}
	}
	if ip == "" {
		ip = r.RemoteAddr
		// 如果包含端口号，去掉端口号
		if strings.Contains(ip, ":") {
			ip = strings.Split(ip, ":")[0]
		}
	}
	return ip
}

// 添加一个函数来解析User-Agent
func parseUserAgent(userAgent string) (os, browser, device string) {
	// 默认值
	os = "Unknown"
	browser = "Unknown"
	device = "Unknown"

	// 转换为小写以便匹配
	ua := strings.ToLower(userAgent)

	// 操作系统检测
	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os x"):
		os = "MacOS"
	case strings.Contains(ua, "linux"):
		os = "Linux"
	case strings.Contains(ua, "android"):
		os = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		os = "iOS"
	}

	// 浏览器检测
	switch {
	case strings.Contains(ua, "chrome"):
		browser = "Chrome"
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "safari"):
		browser = "Safari"
	case strings.Contains(ua, "edge"):
		browser = "Edge"
	case strings.Contains(ua, "opera"):
		browser = "Opera"
	}

	// 设备类型检测
	switch {
	case strings.Contains(ua, "mobile"):
		device = "Mobile"
	case strings.Contains(ua, "tablet"):
		device = "Tablet"
	case strings.Contains(ua, "ipad"):
		device = "iPad"
	case strings.Contains(ua, "iphone"):
		device = "iPhone"
	default:
		device = "Desktop"
	}

	return os, browser, device
}

// HandleConnection 处理 WebSocket 连接的升级和消息接收
func (po *PostOffice) HandleConnection(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in HandleConnection: %v\n", r)
		}
	}()

	// 获取客户端信息
	clientIP := getClientIP(r)
	userAgent := r.Header.Get("User-Agent")
	os, browser, device := parseUserAgent(userAgent)

	if atomic.LoadInt32(&po.connectionCount) >= po.maxConnections {
		http.Error(w, "Connection limit reached", http.StatusServiceUnavailable)
		return
	}

	// 验证连接
	if !po.validateConnection(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	clientID := r.URL.Query().Get("clientID")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	// 获取时间戳
	tsStr := r.URL.Query().Get("startTime")
	var ts int64
	if tsStr != "" {
		ts, _ = strconv.ParseInt(tsStr, 10, 64)
	} else {
		ts = time.Now().Unix()
	}

	// 检查是否存在旧连接
	if oldClient, ok := po.clients.Load(clientID); ok {
		// 获取旧连接的创建时间
		oldTs, _ := po.clientTimes.Load(clientID)
		if oldTs != nil && oldTs.(int64) >= ts {
			// 如果旧连接的时间戳大于等于新连接，拒绝新连接
			http.Error(w, "Connection rejected: newer connection exists", http.StatusConflict)
			return
		}

		// 关闭旧连接
		oldClient.(*Client).conn.Close()
		po.clients.Delete(clientID)
		po.clientTimes.Delete(clientID)
		fmt.Printf("%s Old client connection closed: %s (old ts: %d, ip: %s, os: %s, browser: %s, device: %s)\n",
			currentTime(), clientID, oldTs.(int64), oldClient.(*Client).ip, oldClient.(*Client).os, oldClient.(*Client).browser, oldClient.(*Client).device)
	}

	conn, err := po.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}

	atomic.AddInt32(&po.connectionCount, 1)
	defer atomic.AddInt32(&po.connectionCount, -1)

	// 存储新连接和时间戳
	po.clients.Store(clientID, &Client{
		conn:      conn,
		createdAt: ts,
		ip:        clientIP,
		userAgent: userAgent,
		os:        os,
		browser:   browser,
		device:    device,
	})
	po.clientTimes.Store(clientID, ts)
	fmt.Printf("%s New client connected: %s (ts: %d, ip: %s, os: %s, browser: %s, device: %s)\n",
		currentTime(), clientID, ts, clientIP, os, browser, device)

	// 创建消息通道
	msgChan := make(chan []byte, 100)

	// 启动 goroutine 处理消息
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic in message processing: %v\n", r)
			}
		}()

		for message := range msgChan {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Recovered from panic in messageDelivery: %v\n", r)
					}
				}()
				po.messageDelivery(message)
			}()
		}
	}()

	defer func() {
		close(msgChan) // 确保在关闭连接之前关闭消息通道
		conn.Close()
		po.clients.Delete(clientID)
		po.clientTimes.Delete(clientID)
		fmt.Printf("%s Client disconnected: %s\n", currentTime(), clientID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Connection error for %s: %v\n", clientID, err)
			}
			close(msgChan) // 关闭消息通道
			return
		}

		select {
		case msgChan <- message:
		default:
			// 队列满时静默丢弃消息
		}
	}
}

// messageDelivery 处理消息分发
func (po *PostOffice) messageDelivery(msg []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in messageDelivery: %v\n", r)
		}
	}()

	// 如果存在 schema，则进行验证
	if po.schema != nil {
		result, err := po.schema.Validate(gojsonschema.NewBytesLoader(msg))
		if err != nil {
			fmt.Printf("Error validating message: %v\n", err)
			return
		}
		if !result.Valid() {
			fmt.Println("Invalid message schema")
			return
		}
	}

	// 消息处理逻辑
	msgJson := gjson.ParseBytes(msg)

	// 只打印关键字段，不打印整个消息内容
	from := msgJson.Get("from").String()
	to := msgJson.Get("to")
	msgType := msgJson.Get("type").String()

	if to.Exists() {
		if to.IsArray() {
			targets := to.Array()
			// 只打印消息的基本信息
			if msgType != "heartbeat" { // 心跳消息不打印
				fmt.Printf("%s Broadcasting message from %s to %d targets, type: %s\n",
					currentTime(), from, len(targets), msgType)
			}
			for _, target := range targets {
				po.delivery(target.String(), msg)
			}
		} else {
			target := to.String()
			if msgType != "heartbeat" { // 心跳消息不打印
				fmt.Printf("%s Sending message from %s to %s, type: %s\n",
					currentTime(), from, target, msgType)
			}
			po.delivery(target, msg)
		}
	}
}

// delivery 向指定的目标发送消息
func (po *PostOffice) delivery(targetChannelId string, msg []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in delivery to %s: %v\n", targetChannelId, r)
		}
	}()

	if clientInterface, ok := po.clients.Load(targetChannelId); ok {
		if client, ok := clientInterface.(*Client); ok {
			lockChan := make(chan struct{})
			go func() {
				client.mutex.Lock()
				close(lockChan)
			}()

			select {
			case <-lockChan:
				defer client.mutex.Unlock()
			case <-time.After(5 * time.Second):
				fmt.Printf("%s Lock timeout for client %s\n", currentTime(), targetChannelId)
				return
			}

			writeTimeout := time.Now().Add(5 * time.Second)
			if err := client.conn.SetWriteDeadline(writeTimeout); err != nil {
				return
			}

			// 检查连接状态
			if err := client.conn.WriteControl(websocket.PingMessage, []byte{}, writeTimeout); err != nil {
				fmt.Printf("WriteControl -- error %s \n", targetChannelId)
				po.clients.Delete(targetChannelId)
				return
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Printf("WriteMessage -- error %s \n", targetChannelId)
				po.clients.Delete(targetChannelId)
			}

			// 重置写入超时
			client.conn.SetWriteDeadline(time.Time{})
		}
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
