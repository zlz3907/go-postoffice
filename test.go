//go:build test
// +build test

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsServerURL     = "ws://localhost:7502/ws"
	numConnections  = 3 // 尝试连接的总数
	messageInterval = 1 * time.Second
	authToken       = "your-auth-token-here" // 替换为实际的认证令牌
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type Message struct {
	From    string      `json:"from"`
	To      interface{} `json:"to"` // 可以是字符串或字符串数组
	Subject string      `json:"subject"`
	Content string      `json:"content"` // base64 编码的内容
	Type    int         `json:"type"`
}

var clients []*Client
var mutex sync.Mutex

func main() {
	// 尝试创建多个客户端连接
	for i := 0; i < numConnections; i++ {
		client, err := connectWebSocket()
		if err != nil {
			log.Printf("Failed to connect client %d: %v", i, err)
			continue
		}
		clients = append(clients, client)
		go readMessages(client)
		log.Printf("Client %s connected successfully", client.ID)
	}

	log.Printf("Successfully connected clients: %d", len(clients))

	// 每秒发送消息
	ticker := time.NewTicker(messageInterval)
	for range ticker.C {
		if len(clients) > 0 {
			sender := clients[rand.Intn(len(clients))]
			sendMessageToAll(sender)
		}
	}
}

func connectWebSocket() (*Client, error) {
	clientID := fmt.Sprintf("Client-%d", rand.Intn(1000))
	u, _ := url.Parse(wsServerURL)
	q := u.Query()
	q.Set("clientID", clientID)
	u.RawQuery = q.Encode()

	// 创建自定义的 header
	header := http.Header{}
	header.Add("Authorization", "Bearer "+authToken)

	// 使用自定义的 Dialer
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}

	conn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		return nil, err
	}
	return &Client{
		ID:   clientID,
		Conn: conn,
	}, nil
}

func readMessages(client *Client) {
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message for %s: %v", client.ID, err)
			removeClient(client)
			return
		}
		log.Printf("%s received: %s", client.ID, string(message))
	}
}

func sendMessageToAll(sender *Client) {
	mutex.Lock()
	defer mutex.Unlock()

	var recipients []string
	for _, client := range clients {
		if client.ID != sender.ID {
			recipients = append(recipients, client.ID)
		}
	}

	content := fmt.Sprintf("Hello from %s at %s", sender.ID, time.Now().Format(time.RFC3339))
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))

	message := Message{
		From:    sender.ID,
		To:      recipients,
		Subject: "Test Message",
		Content: encodedContent,
		Type:    1, // 假设 1 代表普通消息
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	log.Printf("Sending message: %s", string(jsonMessage))
	err = sender.Conn.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		log.Printf("Error sending message from %s: %v", sender.ID, err)
		removeClient(sender)
	} else {
		log.Printf("Message sent successfully from %s", sender.ID)
	}
}

func removeClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, c := range clients {
		if c.ID == client.ID {
			clients = append(clients[:i], clients[i+1:]...)
			client.Conn.Close()
			log.Printf("Client %s removed", client.ID)
			return
		}
	}
}
