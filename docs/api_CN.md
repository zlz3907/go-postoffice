# GO-POSTOFFICE API 文档

[English](api.md) | 中文

本文档描述了 GO-POSTOFFICE 项目的 API。

## WebSocket 连接

要与 GO-POSTOFFICE 服务器建立 WebSocket 连接：

```
ws://localhost:7502/ws
```

如果连接到远程服务器，请将 `localhost` 替换为适当的主机地址。

### 连接参数

- `clientID`：客户端的唯一标识符。这应该作为查询参数包含在 WebSocket URL 中。

示例：
```
ws://localhost:7502/ws?clientID=user123
```

### 认证

认证是通过在 WebSocket 连接请求头中使用 Bearer 令牌完成的。

示例：
```
Authorization: Bearer your_token_here
```

## 消息格式

通过 WebSocket 连接发送和接收的消息应遵循以下 JSON 格式：

```json
{
  "from": "sender_id",
  "to": "recipient_id",
  "subject": "消息主题",
  "content": "消息内容",
  "type": "msg"
}
```

### 消息字段

- `from`：发送者的 ID（必填）
- `to`：接收者的 ID 或接收者 ID 数组（必填）
- `subject`：消息的主题（必填）
- `content`：消息的内容（必填）
- `type`：消息的类型（必填，可选值："log"、"heartbeat"、"msg"）

其他可选字段：

- `cc`：抄送接收者（字符串或字符串数组）
- `contentType`：表示内容类型的整数
- `charset`：指定字符编码的字符串
- `level`：表示消息优先级的整数（默认：0）
- `tags`：用于消息分类的字符串数组
- `attachments`：附件对象数组
- `references`：用于消息线程的字符串
- `inReplyTo`：标识此消息回复的消息的字符串
- `subjectId`：用于分组相关消息的字符串
- `createTime`：消息创建的整数时间戳
- `lastUpdateTime`：最后更新消息的整数时间戳
- `state`：表示消息状态的整数
- `token`：用于额外认证或验证的字符串
- `fromTag`：用于分类发送者的字符串

## 发送消息

要发送消息，只需通过已建立的 WebSocket 连接发送上述格式的 JSON 对象。

## 接收消息

来自服务器的消息将通过 WebSocket 连接以上述相同的 JSON 格式接收。

## 错误处理

如果发生错误，服务器将发送以下格式的错误消息：

```json
{
  "error": "错误消息描述"
}
```

## 示例

### 连接到 WebSocket 服务器

```javascript
const ws = new WebSocket('ws://localhost:7502/ws?clientID=user123');
ws.onopen = () => {
  console.log('已连接到服务器');
};
```

### 发送消息

```javascript
ws.send(JSON.stringify({
  from: "user123",
  to: "user456",
  subject: "你好",
  content: "你好吗？",
  type: "msg"
}));
```

### 接收消息

```javascript
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('收到消息:', message);
};
```

有关各种编程语言的更详细示例，请参阅项目仓库中的 `examples` 目录。