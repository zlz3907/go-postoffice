# GO-POSTOFFICE 消息协议

[English](message-protocol.md) | 中文

本文档描述了 GO-POSTOFFICE 项目中使用的消息协议。

## 消息结构

GO-POSTOFFICE 中的消息遵循受电子邮件协议启发的 JSON 结构。以下是消息的一般结构：

```json
{
  "from": "sender_id",
  "to": "recipient_id",
  "subject": "消息主题",
  "content": "消息内容",
  "type": "msg",
  "cc": ["cc_recipient1", "cc_recipient2"],
  "contentType": 0,
  "charset": "UTF-8",
  "level": 0,
  "tags": ["tag1", "tag2"],
  "attachments": [],
  "references": "reference_id",
  "inReplyTo": "original_message_id",
  "subjectId": "subject_thread_id",
  "createTime": 1623456789,
  "lastUpdateTime": 1623456789,
  "state": 0,
  "token": "authentication_token",
  "fromTag": "sender_category"
}
```

## 字段描述

### 必填字段

- `from` (字符串)：消息发送者的标识符。
- `to` (字符串或字符串数组)：消息接收者的标识符。
- `subject` (字符串)：消息的主题或标题。
- `content` (任意类型)：消息的主体��可以是字符串或更复杂的对象。
- `type` (字符串)：消息的类型。必须是以下之一："log"、"heartbeat" 或 "msg"。

### 可选字段

- `cc` (字符串或字符串数组)：抄送接收者。
- `contentType` (整数)：消息内容类型的标识符。
- `charset` (字符串)：消息内容的字符编码。
- `level` (整数)：消息的优先级或重要性级别。默认为 0。
- `tags` (字符串数组)：与消息相关的类别或标签。
- `attachments` (数组)：附加到消息的任何文件或额外数据。
- `references` (字符串)：用于消息线程或分组的标识符。
- `inReplyTo` (字符串)：此消息回复的原始消息的标识符。
- `subjectId` (字符串)：用于按主题或对话分组消息的标识符。
- `createTime` (整数)：消息创建时的 Unix 时间戳。
- `lastUpdateTime` (整数)：消息最后更新时的 Unix 时间戳。
- `state` (整数)：消息的当前状态。
- `token` (字符串)：认证或验证令牌。
- `fromTag` (字符串)：发送者的类别或类型。

## 消息类型

- `log`：用于日志或系统消息。
- `heartbeat`：用于连接保活消息。
- `msg`：用于常规的用户对用户或应用程序消息。

## 示例

### 简单消息

```json
{
  "from": "user123",
  "to": "user456",
  "subject": "你好",
  "content": "你好吗？",
  "type": "msg"
}
```

### 带有可选字段的复杂消息

```json
{
  "from": "system",
  "to": ["user1", "user2"],
  "subject": "系统更新",
  "content": {
    "text": "系统将进行维护停机。",
    "downtime": "2小时"
  },
  "type": "log",
  "level": 2,
  "tags": ["维护", "系统"],
  "createTime": 1623456789,
  "state": 1
}
```

## 验证

消息会根据 JSON schema 进行验证，以确保它们符合预期的结构。该 schema 定义在项目根目录的 `message_schema.json` 文件中。

有关 schema 和验证过程的详细信息，请参阅 `message_schema.json` 文件。