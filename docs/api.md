# GO-POSTOFFICE API Documentation

[中文](api_CN.md) | English

This document describes the API for the GO-POSTOFFICE project.

## WebSocket Connection

To establish a WebSocket connection with the GO-POSTOFFICE server:

```
ws://localhost:7502/
```

Replace `localhost` with the appropriate host if connecting to a remote server.

### Connection Parameters

- `clientID`: A unique identifier for the client. This should be included as a query parameter in the WebSocket URL.

Example:
```
ws://localhost:7502/?clientID=user123
```

### Authentication

Authentication is done using a Bearer token in the WebSocket connection request headers.

Example:
```
Authorization: Bearer your_token_here
```

## Message Format

Messages sent and received through the WebSocket connection should follow this JSON format:

```json
{
  "from": "sender_id",
  "to": "recipient_id",
  "subject": "Message subject",
  "content": "Message content",
  "type": "msg"
}
```

### Message Fields

- `from`: The ID of the sender (required)
- `to`: The ID of the recipient or an array of recipient IDs (required)
- `subject`: The subject of the message (required)
- `content`: The content of the message (required)
- `type`: The type of the message (required, one of: "log", "heartbeat", "msg")

Additional optional fields:

- `cc`: Carbon copy recipients (string or array of strings)
- `contentType`: Integer representing the content type
- `charset`: String specifying the character encoding
- `level`: Integer representing the message priority level (default: 0)
- `tags`: Array of strings for message categorization
- `attachments`: Array of attachment objects
- `references`: String for message threading
- `inReplyTo`: String identifying the message this is a reply to
- `subjectId`: String for grouping related messages
- `createTime`: Integer timestamp of message creation
- `lastUpdateTime`: Integer timestamp of last message update
- `state`: Integer representing the message state
- `token`: String for additional authentication or validation
- `fromTag`: String for categorizing the sender

## Sending Messages

To send a message, simply send a JSON object in the above format through the established WebSocket connection.

## Receiving Messages

Messages from the server will be received through the WebSocket connection in the same JSON format described above.

## Error Handling

If an error occurs, the server will send an error message in the following format:

```json
{
  "error": "Error message description"
}
```

## Examples

### Connecting to the WebSocket server

```javascript
const ws = new WebSocket('ws://localhost:7502/?clientID=user123');
ws.onopen = () => {
  console.log('Connected to the server');
};
```

### Sending a message

```javascript
ws.send(JSON.stringify({
  from: "user123",
  to: "user456",
  subject: "Hello",
  content: "How are you?",
  type: "msg"
}));
```

### Receiving a message

```javascript
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received message:', message);
};
```

For more detailed examples in various programming languages, please refer to the `examples` directory in the project repository.