# 实时通信 WebSocket 服务器

本项目实现了一个专为实时通信设计的 WebSocket 服务器，适用于需要即时消息传递、实时更新或任何需要客户端和服务器之间低延迟数据交换的应用场景。

## 项目结构

本项目的主要组件包括：

1. `main.go`：应用程序的入口点。负责初始化环境、启动 WebSocket 服务器，并管理程序生命周期。

2. `go.work`：定义 Go 工作空间并指定项目中使用的模块。

3. `.env/config-dev.json` 和 `.env/config-zhycit.json`：不同环境（开发和生产）的配置文件。包含 Redis、WebSocket 端口和最大连接数等设置。

4. `ai.zhycit.com/socket` 包：（在提供的文件中未显示）这个包可能包含了 `PostOffice` 结构体的实现，用于处理 WebSocket 连接。

## 主要特性

- 支持特定环境的配置加载
- 可配置最大连接数的 WebSocket 服务器
- 集成 Redis 用于数据持久化（配置存在，实现未显示）
- 优雅的关闭机制

## 消息协议模式

以下表格描述了此 WebSocket 服务器使用的消息协议模式：

| 字段           | 类型    | 描述                                | 是否必需 |
|----------------|---------|-------------------------------------|----------|
| from           | string  | 消息来源                            | 是       |
| to             | array   | 消息接收者                          | 是       |
| subject        | string  | 消息主题                            | 是       |
| content        | string  | 消息内容（base64 编码）             | 是       |
| type           | integer | 消息类型（1, 2, 3, 或 4）           | 是       |
| cc             | array   | 抄送接收者                          | 否       |
| contentType    | integer | 消息内容类型                        | 否       |
| charset        | string  | 字符编码                            | 否       |
| level          | integer | 消息优先级（默认：0）               | 否       |
| tags           | array   | 与消息关联的标签                    | 否       |
| attachments    | array   | 附件                                | 否       |
| references     | string  | 相关主题 ID                         | 否       |
| inReplyTo      | string  | 被回复消息的 ID                     | 否       |
| subjectId      | string  | 主题 ID                             | 否       |
| createTime     | integer | 创建时间（Unix 时间戳）             | 否       |
| lastUpdateTime | integer | 最后更新时间（Unix 时间戳）         | 否       |
| state          | integer | 消息发送状态                        | 否       |
| token          | string  | 认证令牌                            | 否       |
| fromTag        | string  | 来源标签（例如：QQ, APP, TAB）      | 否       |

## 构建可执行文件

要将此项目构建成可执行文件，请按以下步骤操作：

1. 确保您的系统上安装了 Go（版本 1.23.1 或更高，如 `go.work` 中指定）。

2. 打开终端并导航到项目根目录。

3. 运行以下命令来构建可执行文件：

   ```
   go build -o websocket-server main.go
   ```

   这将在当前目录中创建一个名为 `websocket-server`（在 Windows 上为 `websocket-server.exe`）的可执行文件。

4. （可选）要为特定平台构建，可以使用 `GOOS` 和 `GOARCH` 环境变量。例如，要为 Windows 构建：

   ```
   GOOS=windows GOARCH=amd64 go build -o websocket-server.exe main.go
   ```

5. 生成的可执行文件可以直接在目标系统上运行，无需安装 Go。

部署可执行文件时，请记得包含必要的配置文件（`.env` 文件夹），以确保在不同环境中正常运行。