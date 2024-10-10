package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"ai.zhycit.com/socket" // 导入 PostOffice 包
)

//go:embed .env
var envFolder embed.FS

var env = "dev"
var gnasConfig map[string]interface{}

func initEnv(env string) {
	envEntries, _ := fs.ReadDir(envFolder, ".env")
	for _, entry := range envEntries { // 读取环境配置文件
		if strings.Index(entry.Name(), "config-"+env) == 0 { // Fixed Yoda condition
			log.Println("当前运行环境: " + entry.Name())
			data, err := envFolder.ReadFile(".env/" + entry.Name())
			if nil != err {
				log.Fatalln("Error reading embedded file:", err)
			}

			parseErr := json.Unmarshal(data, &gnasConfig)
			if parseErr != nil {
				log.Println("Error unmarshalling JSON:", err)
			}
			break
		}
	}
}

func startWebSocketServer() {
	log.Println("Starting WebSocket server on port:", gnasConfig["socketPort"])

	maxConnections := 20000 // 默认值
	if maxConn, ok := gnasConfig["maxWebSocketConnections"].(float64); ok {
		maxConnections = int(maxConn)
	}

	// 如果要使用 schema 验证
	// postOffice, err := socket.NewPostOffice(maxConnections, "message_schema.json")

	// 如果不使用 schema 验证
	postOffice, err := socket.NewPostOffice(maxConnections, "")

	if err != nil {
		log.Fatalf("Failed to create PostOffice: %v", err)
	}

	http.HandleFunc("/ws", postOffice.HandleConnection)

	port := fmt.Sprintf(":%v", gnasConfig["socketPort"])
	// log.Printf("WebSocket server listening on port %s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("Error starting WebSocket server: %v", err)
		// 不使用 panic，而是让函数正常返回，这样可以触发 main 中的程序退出逻辑
	}
}

func main() {
	// 初始化系统
	initEnv(env)

	// 创建一个 channel 用于保持程序运行
	done := make(chan bool)

	// 启动 WebSocket 服务器
	go func() {
		startWebSocketServer()
		// 如果 startWebSocketServer 返回（比如发生错误），我们关闭 channel
		close(done)
	}()

	// log.Println("WebSocket server started on port:", gnasConfig["socketPort"])

	// 等待 done channel 关闭，这会阻塞主程序
	<-done

	log.Println("WebSocket server stopped. Exiting program.")
}
