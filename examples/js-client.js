const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:7502/ws');

ws.on('open', function open() {
  console.log('Connected to the server');
  sendMessage();
});

ws.on('message', function incoming(data) {
  console.log('Received:', data.toString());
});

ws.on('close', function close() {
  console.log('Disconnected from the server');
});

ws.on('error', function error(err) {
  console.error('WebSocket error: ', err);
});

function sendMessage() {
  const message = {
    from: 'js-client',
    to: 'server',
    subject: 'Hello',
    content: 'How are you?',
    type: 'msg'
  };
  ws.send(JSON.stringify(message));
}

// Send a message every 5 seconds
setInterval(sendMessage, 5000);