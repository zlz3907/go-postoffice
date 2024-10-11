# GO-POSTOFFICE Message Protocol

[中文](message-protocol_CN.md) | English

This document describes the message protocol used in the GO-POSTOFFICE project.

## Message Structure

Messages in GO-POSTOFFICE follow a JSON structure inspired by email protocols. Here's the general structure of a message:

```json
{
  "from": "sender_id",
  "to": "recipient_id",
  "subject": "Message subject",
  "content": "Message content",
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

## Field Descriptions

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| from | string | The identifier of the message sender |
| to | string or array of strings | The identifier(s) of the message recipient(s) |
| subject | string | The subject or title of the message |
| content | any | The main body of the message |
| type | string | The type of the message (e.g., "log", "heartbeat", "msg") |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| cc | string or array of strings | Carbon copy recipients |
| contentType | integer | An identifier for the type of content in the message |
| charset | string | The character encoding of the message content |
| level | integer | The priority or importance level of the message |
| tags | array of strings | Categories or labels associated with the message |
| attachments | array | Any files or additional data attached to the message |
| references | string | Identifier for message threading or grouping |
| inReplyTo | string | Identifier of the message this is in reply to |
| subjectId | string | Identifier for grouping messages by subject or conversation |
| createTime | integer | Unix timestamp of when the message was created |
| lastUpdateTime | integer | Unix timestamp of when the message was last updated |
| state | integer | Current state or status of the message |
| token | string | Authentication or validation token |
| fromTag | string | Category or type of the sender |

## Message Types

- `log`: Used for logging or system messages.
- `heartbeat`: Used for connection keep-alive messages.
- `msg`: Used for regular user-to-user or application messages.

## Examples

### Simple Message

```json
{
  "from": "user123",
  "to": "user456",
  "subject": "Hello",
  "content": "How are you?",
  "type": "msg"
}
```

### Complex Message with Optional Fields

```json
{
  "from": "system",
  "to": ["user1", "user2"],
  "subject": "System Update",
  "content": {
    "text": "The system will be down for maintenance.",
    "downtime": "2 hours"
  },
  "type": "log",
  "level": 2,
  "tags": ["maintenance", "system"],
  "createTime": 1623456789,
  "state": 1
}
```

## Validation

Messages are validated against a JSON schema to ensure they conform to the expected structure. The schema is defined in the `message_schema.json` file in the project root.

For detailed information about the schema and validation process, please refer to the `message_schema.json` file.