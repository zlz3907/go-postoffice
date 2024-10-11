const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:7502/ws');

ws.on('open', function open() {
  console.log('Connected to the server');
  ws.send('Hello, server!');
});

ws.on('message', function incoming(data) {
  console.log('Received: %s', data);
});

ws.on('close', function close() {
  console.log('Disconnected from the server');
});

ws.on('error', function error(err) {
  console.error('WebSocket error: ', err);
});

// Keep the script running
setInterval(() => {
  ws.send('Ping');
}, 30000);