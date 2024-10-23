package main

import (
	"encoding/json"
	"log"
	"net/url"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	From    string      `json:"from"`
	To      string      `json:"to"`
	Subject string      `json:"subject"`
	Content interface{} `json:"content"`
	Type    string      `json:"type"`
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	clientID := "go-client-001" // 设置客户端ID
	u := url.URL{Scheme: "ws", Host: "localhost:7502", Path: "/", RawQuery: "clientID=" + clientID}
	log.Printf("connecting to %s", u.String())

	// Create custom header with token
	header := http.Header{}
	header.Add("Authorization", "Bearer your_token_here")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			message := Message{
				From:    "go-client",
				To:      "server",
				Subject: "Hello",
				Content: "How are you?",
				Type:    "msg",
			}
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				log.Println("json marshal:", err)
				return
			}
			err = c.WriteMessage(websocket.TextMessage, jsonMessage)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
