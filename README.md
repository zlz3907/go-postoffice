# WebSocket Server for Real-time Communication

[中文版](README_CN.md)

This project implements a WebSocket server designed for real-time communication, suitable for applications requiring instant messaging, live updates, or any scenario demanding low-latency data exchange between clients and server.

## Project Structure

The main components of this project are:

1. `main.go`: The entry point of the application. It initializes the environment, starts the WebSocket server, and manages the program lifecycle.

2. `go.work`: Defines the Go workspace and specifies the modules used in the project.

3. `.env/config-dev.json` and `.env/config-zhycit.json`: Configuration files for different environments (development and production). They contain settings for Redis, WebSocket port, and maximum connections.

4. `ai.zhycit.com/socket` package: (Not shown in the provided files) This package likely contains the implementation of the `PostOffice` struct, which handles WebSocket connections.

## Key Features

- Environment-specific configuration loading
- WebSocket server with configurable maximum connections
- Integration with Redis for data persistence (configuration present, implementation not shown)
- Graceful shutdown mechanism

## Message Protocol Schema

The following table describes the schema for the message protocol used in this WebSocket server:

| Field          | Type    | Description                                | Required |
|----------------|---------|--------------------------------------------| -------- |
| from           | string  | Message source                             | Yes      |
| to             | string/array | Message recipient(s)                  | Yes      |
| subject        | string  | Message subject                            | Yes      |
| content        | string  | Message content (base64 encoded)           | Yes      |
| type           | integer | Message type (1, 2, 3, or 4)               | Yes      |
| cc             | string/array | Carbon copy recipient(s)              | No       |
| contentType    | integer | Content type of the message                | No       |
| charset        | string  | Character encoding                         | No       |
| level          | integer | Message priority (default: 0)              | No       |
| tags           | array   | Tags associated with the message           | No       |
| attachments    | array   | Attachments                                | No       |
| references     | string  | Related topic ID                           | No       |
| inReplyTo      | string  | ID of the message being replied to         | No       |
| subjectId      | string  | Subject ID                                 | No       |
| createTime     | integer | Creation time (Unix timestamp)             | No       |
| lastUpdateTime | integer | Last update time (Unix timestamp)          | No       |
| state          | integer | Message sending state                      | No       |
| token          | string  | Authentication token                       | No       |
| fromTag        | string  | Source tag (e.g., QQ, APP, TAB)            | No       |

## Building an Executable

To build this project into an executable file, follow these steps:

1. Ensure you have Go installed on your system (version 1.23.1 or later as specified in `go.work`).

2. Open a terminal and navigate to the project root directory.

3. Run the following command to build the executable:

   ```
   go build -o websocket-server main.go
   ```

   This will create an executable named `websocket-server` (or `websocket-server.exe` on Windows) in the current directory.

4. (Optional) To build for a specific platform, you can use the `GOOS` and `GOARCH` environment variables. For example, to build for Windows:

   ```
   GOOS=windows GOARCH=amd64 go build -o websocket-server.exe main.go
   ```

5. The resulting executable can be run directly on the target system without needing Go installed.

Remember to include the necessary configuration files (`.env` folder) when deploying the executable to ensure proper functionality in different environments.
