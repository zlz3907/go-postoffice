# GO-POSTOFFICE

[English](README.md)

GO-POSTOFFICE 是一个基于 Go 语言的高性能 WebSocket 服务器实现，采用邮局的概念模型设计。该项目主要负责连接管理、安全认证和消息（邮件）的分发投递。

## 目录

1. [特点](#特点)
2. [安装](#安装)
3. [快速开始](#快速开始)
4. [配置](#配置)
5. [API 文档](#api-文档)
6. [消息协议](#消息协议)
7. [客户端示例](#客户端示例)
8. [贡献指南](#贡献指南)
9. [许可证](#许可证)

## 特点

1. **高性能并发处理**：利用 Go 语言的 goroutine 和 channel 机制，实现高效的并发连接管理。

2. **灵活的消息路由**：基于邮局模型，支持点对点和广播消息分发，实现高效的消息投递。

3. **安全认证机制**：集成 token 认证，确保连接的安全性。

4. **可配置的消息验证**：支持可选的 JSON Schema 验证，确保消息格式的正确性。

5. **环境适应性**：支持多环境配置，便于在不同场景下部署。

6. **优雅的服务管理**：实现了优雅启动和关闭机制，确保服务的稳定性。

7. **可扩展性**：模块化设计，便于功能扩展和定制。

8. **实时通信**：基于 WebSocket 的全双工通信，支持实时数据交换。

### 邮局模型的优势

- **解耦性**：发送者和接收者完全分离，提高系统的灵活性。
- **可靠性**：消息持久化和重试机制确保消息的可靠投递。
- **扩展性**：易于添加新的消息类型和处理逻辑。
- **负载均衡**：可以实现多个"邮局"实例，提高系统的吞吐量。

## 安装

确保您的系统已安装 Go（版本 1.23.1 或更高）。

1. 克隆仓库：
   ```
   git clone https://github.com/your-username/GO-POSTOFFICE.git
   ```

2. 进入项目目录：
   ```
   cd GO-POSTOFFICE
   ```

3. 安装依赖：
   ```
   go mod tidy
   ```

## 快速开始

1. 配置环境：
   复制 `.env/config-dev.json` 到 `.env/config-zhycit.json` 并根据需要修改配置。

2. 运行服务器：
   ```
   go run main.go
   ```

3. 构建可执行文件：
   ```
   go build -o postoffice main.go
   ```

## 配置

主要配置项包括：

- `socketPort`: WebSocket 服务器端口
- `maxWebSocketConnections`: 最大连接数
- `dataSource`: 数据源配置（如 Redis）

详细配置说明请参考 [配置文档](docs/configuration.md)。

## API 文档

API 使用说明请参考 [API 文档](docs/api.md)。

## 消息协议

消息格式和字段说明请参考 [消息协议文档](docs/message-protocol.md)。

## 客户端示例

- [Python 客户端示例](examples/python-client.py)
- [JavaScript 客户端示例](examples/js-client.js)

## 贡献指南

我们欢迎任何形式的贡献。请阅读 [贡献指南](CONTRIBUTING.md) 了解如何参与项目开发。

## 许可证

本项目采用 MIT 许可证。详情请见 [LICENSE](LICENSE) 文件。