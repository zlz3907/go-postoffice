# GO-POSTOFFICE 配置

[English](configuration.md) | 中文

本文档描述了 GO-POSTOFFICE 项目的配置选项。

## 配置文件

配置文件位于 `.env` 目录中。有两个主要的配置文件：

- `config-dev.json`：用于开发环境
- `config-zhycit.json`：用于生产环境

## 配置选项

以下是主要的配置选项：

### WebSocket 服务器

- `socketPort`：WebSocket 服务器将监听的端口号。
  - 类型：整数
  - 默认值：7502

- `maxWebSocketConnections`：允许的最大并发 WebSocket 连接数。
  - 类型：整数
  - 默认值：20000（生产环境）

### SSL 配置（可选）

- `sslPort`：SSL 连接的端口号。
  - 类型：整数
  - 默认值：7503

- `sslCertPath`：SSL 证书文件的路径。
  - 类型：字符串
  - 示例："./cert/cert.pem"

- `sslKeyPath`：SSL 密钥文件的路径。
  - 类型：字符串
  - 示例："./cert/key.pem"

### 数据源

- `dataSource`：外部数据源的配置。
  - 类型：对象

  #### Redis 配置

  - `redis`：Redis 连接的配置。
    - 类型：对象
    - 属性：
      - `gnas-ai`：特定 Redis 实例的配置。
        - 类型：对象
        - 属性：
          - `uri`：连接到 Redis 实例的 URI。
            - 类型：字符串
            - 示例："127.0.0.1:6379"（开发环境）

## 配置示例

以下是一个完整的配置文件示例：

```json
{
  "dataSource": {
    "redis": {
      "gnas-ai": {
        "uri": "127.0.0.1:6379"
      }
    }
  },
  "socketPort": 7502,
  "sslPort": 7503,
  "sslCertPath": "./cert/cert.pem",
  "sslKeyPath": "./cert/key.pem",
  "maxWebSocketConnections": 2
}
```

## 更改配置

要更改配置：

1. 复制适当的配置文件（例如，开发环境使用 `config-dev.json`）到 `.env/config-zhycit.json`。
2. 根据需要修改 `config-zhycit.json` 中的值。
3. 重启 GO-POSTOFFICE 服务器以使更改生效。

请记住，切勿将敏感信息（如生产数据库凭据）提交到版本控制系统。在生产环境中，请使用环境变量或安全的密钥管理系统来处理敏感数据。