# GO-POSTOFFICE

[中文版](README_CN.md)

GO-POSTOFFICE is a high-performance WebSocket server implementation based on Go, designed using the post office concept model. This project primarily handles connection management, security authentication, and message (mail) distribution and delivery.

![Go-Postoffice Communication Structure](docs/imgs/global_architecture_diagram_en.png)

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Configuration](#configuration)
5. [API Documentation](#api-documentation)
6. [Message Protocol](#message-protocol)
7. [Client Examples](#client-examples)
8. [Contributing](#contributing)
9. [License](#license)

## Features

1. **High-performance Concurrent Processing**: Utilizes Go's goroutines and channels for efficient concurrent connection management.
2. **Flexible Message Routing**: Based on the post office model, supports point-to-point and broadcast message distribution for efficient message delivery.
3. **Security Authentication**: Integrates token authentication to ensure connection security.
4. **Configurable Message Validation**: Supports optional JSON Schema validation to ensure message format correctness.
5. **Environmental Adaptability**: Supports multi-environment configuration for easy deployment in different scenarios.
6. **Graceful Service Management**: Implements graceful startup and shutdown mechanisms to ensure service stability.
7. **Scalability**: Modular design for easy feature expansion and customization.
8. **Real-time Communication**: Full-duplex communication based on WebSocket, supporting real-time data exchange.

### Advantages of the Post Office Model

- **Decoupling**: Complete separation of senders and receivers, increasing system flexibility.
- **Reliability**: Message persistence and retry mechanisms ensure reliable message delivery.
- **Extensibility**: Easy to add new message types and processing logic.
- **Load Balancing**: Multiple "post office" instances can be implemented to increase system throughput.

## Quick Start Guide for Enterprise AI Chatbot Integration
For a detailed guide on how to quickly integrate an enterprise-level AI chatbot, please refer to our [Enterprise AI Chatbot Integration Guide](docs/enterprise_ai_chatbot_integration_guide.md).

## Installation

Ensure that Go (version 1.23.1 or higher) is installed on your system.

1. Clone the repository:
   ```
   git clone https://github.com/zlz3907/go-postoffice.git
   ```

2. Enter the project directory:
   ```
   cd go-postoffice
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

## Quick Start

1. Configure the environment:
   Copy `.env/config-dev.json` to `.env/config-zhycit.json` and modify the configuration as needed.

2. Run the server:
   ```
   go run main.go
   ```

3. Build the executable:

   For Linux:
   ```
   env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.env=zhycit" -o dist/go-postoffice-linux
   ```

   For macOS:
   ```
   env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.env=zhycit" -o dist/go-postoffice-macos
   ```

   For Windows:
   ```
   env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.env=zhycit" -o dist/go-postoffice-windows.exe
   ```

   Note: Replace `zhycit` with your desired environment name if different.

4. Run the built executable:

   For Linux/macOS:
   ```
   ./dist/go-postoffice-linux   # or go-postoffice-macos
   ```

   For Windows:
   ```
   .\dist\go-postoffice-windows.exe
   ```

## Configuration

Main configuration items include:

- `socketPort`: WebSocket server port
- `maxWebSocketConnections`: Maximum number of connections
- `dataSource`: Data source configuration (e.g., Redis)

For detailed configuration instructions, please refer to the [Configuration Documentation](docs/configuration.md).

## API Documentation

For API usage instructions, please refer to the [API Documentation](docs/api.md).

## Message Protocol

For message format and field descriptions, please refer to the [Message Protocol Documentation](docs/message-protocol.md).

## Client Examples

- [Go Client Example](examples/go-client.go)
- [Java Client Example](examples/JavaClient.java)
- [JavaScript Client Example](examples/js-client.js)
- [Python Client Example](examples/python-client.py)

## Contributing

We welcome contributions of any form. Please read the [Contributing Guidelines](CONTRIBUTING.md) to learn how to participate in project development.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.
