package main

import (
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ai.zhycit.com/socket"
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

func startWebSocketServer() error {

	maxConnections := 20000 // 默认值
	if maxConn, ok := gnasConfig["maxWebSocketConnections"].(float64); ok {
		maxConnections = int(maxConn)
	}

	// 如果要使用 schema 验证
	// postOffice, err := socket.NewPostOffice(maxConnections, "message_schema.json")

	// 如果不使用 schema 验证
	postOffice, err := socket.NewPostOffice(maxConnections, "")
	if err != nil {
		return fmt.Errorf("failed to create PostOffice: %v", err)
	}

	port := fmt.Sprintf(":%v", gnasConfig["socketPort"])

	// 创建一个新的 ServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", postOffice.HandleConnection)

	// 设置 HTTP 服务器
	httpServer := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// 启动 HTTP 服务器（WS）
	go func() {
		log.Printf("Starting WS server on port %s\n", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting WS server: %v", err)
		}
	}()

	// 检查是否配置了 sslPort
	if sslPort, ok := gnasConfig["sslPort"].(float64); ok {
		// 设置 HTTPS 服务器
		httpsServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", int(sslPort)),
			Handler: mux,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: "socket.zhycit.com",
			},
		}

		// 启动 HTTPS 服务器（WSS）
		go func() {
			if sslPort <= 0 || sslPort > 65535 {
				log.Println("Invalid SSL port, WSS server not started")
				return
			}

			certPath, certOk := gnasConfig["sslCertPath"].(string)
			keyPath, keyOk := gnasConfig["sslKeyPath"].(string)

			if !certOk || !keyOk {
				log.Println("SSL cert or key path not found in config, WSS server not started")
				return
			}

			log.Printf("Starting WSS server on port %d\n", int(sslPort))
			if err := httpsServer.ListenAndServeTLS(certPath, keyPath); err != nil && err != http.ErrServerClosed {
				log.Printf("Error starting WSS server: %v", err)
			}
		}()
	} else {
		log.Println("SSL port not configured, WSS server not started")
	}

	return nil
}

func main() {
	// 初始化系统
	initEnv(env)

	// 启动 WebSocket 服务器
	if err := startWebSocketServer(); err != nil {
		log.Fatalf("Failed to start WebSocket server: %v", err)
	}

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-sigChan

	log.Println("Shutting down WebSocket servers...")
	// 这里可以添加优雅关闭的逻辑

	// log.Println("WebSocket servers stopped. Exiting program.")
}
