# Enterprise AI Chatbot Integration Guide (Socket Client-based Integration)

## 1. Introduction
   - Purpose of this guide
     This guide aims to help developers quickly integrate enterprise-level AI chatbots with WebSocket servers.
   
   - Overview of the integration process
     The integration process includes implementing a Socket client application, establishing secure connections, and handling message exchange between AI chatbots and servers.

## 2. System Requirements
   - Operating System
     - Windows 10 or higher
     - macOS 10.12 or higher
     - Linux (Ubuntu 18.04 LTS or higher recommended)
   
   - Hardware Requirements
     - CPU: Dual-core processor or higher
     - Memory: At least 4GB RAM
     - Storage: At least 10GB available disk space
     - Network Bandwidth: Recommended at least 10Mbps upload and download speed
   
   - Network Environment
     - Stable internet connection
     - Network settings that allow WebSocket communication

## 3. Preparation
   - Obtaining necessary credentials
     - WebSocket server URL, for enterprise version `GO-POSTOFFICE` server, the URL is: `wss://socket.zhycit.com/`. Users can set up their own server, refer to the documentation: [GO-POSTOFFICE Documentation](https://github.com/zhycit/GO-POSTOFFICE)
     - Authorization token (required for enterprise version)

       Use the following API to obtain the authorization token:
       ```
       curl --location 'https://ai.zhycit.com/wecom/token?appid=<APPID>&secret=<SECRET>'
       ```
       
       Example of returned result:
       ```json
       {
           "expires": 1732256675685,
           "token": "af71********6922",
           "ttl": 2589694
       }
       ```
       
       Here, the `token` field is the authorization token, `expires` is the expiration timestamp (end time), and `ttl` is the remaining time of the token (in seconds).
       
       Note: Please keep your appid and secret safe and do not disclose them to others. When using in practice, please use your own appid and secret.
     - Client ID (used to identify each connection)

       The client ID is a unique identifier for each connection, used to distinguish different client connections. Each client can provide corresponding functional services, such as weather queries, web searches, knowledge Q&A, data analysis, etc.

       Client IDs can be defined by users, for example: `user_123456`, `user_654321`, etc. If the client provides specific services, it can be defined as `/service/<company>/weather_1`, `/service/<company>/search_1`, etc.

   - Choosing programming language and WebSocket client library

     Choose a suitable programming language and WebSocket client library based on your technology stack. Common choices include:
     - Python: `websockets` or `socket.io-client`
     - Java: `Java-WebSocket` or `Tyrus`
     - JavaScript: Native WebSocket API or `socket.io-client`

## 4. Implementing Socket Client
   - Establishing WebSocket connection
     Use the selected WebSocket library to connect to the server. Example (using Python):
     ```python
     import websockets
     import asyncio

     async def connect():
         uri = "ws://your-server-url:port"
         async with websockets.connect(uri) as websocket:
             # Handle connection logic
     ```

   - Implementing authentication

     If the server requires authentication, send necessary credentials when establishing the connection. Example:
     ```python
     headers = {
         "Authorization": "Bearer your-token-here",
         "Client-ID": "your-client-id"
     }
     async with websockets.connect(uri, extra_headers=headers) as websocket:
         # Post-connection logic
     ```

   - Message sending and receiving

     Implement functionality to send messages to the server and receive server messages. Example:
     ```python
     async def message_handler(websocket):
         while True:
             message = await websocket.recv()
             # Handle received message
             
             response = "AI-generated reply"
             await websocket.send(response)
     ```

## 5. Message Format and Protocol
   - Defining message structure

     JSON structure specification for enterprise server intercommunication messages (reference for self-built server programs):
     ```json
     {
       "from": "sender_id",
       "to": "recipient_id",
       "subject": "Message subject",
       "content": "User input or AI reply",
       "type": "msg",
       "createTime": 1623456789
     }
     ```

    Protocol reference: [GO-POSTOFFICE Message Protocol](message-protocol.md)

   - Implementing message serialization and deserialization
     Use JSON library to handle message encoding and decoding.

## 6. Error Handling and Reconnection Mechanism
   - Implementing error catching
     Catch and handle potential network errors or server response issues.

   - Adding reconnection logic
     Automatically attempt to reconnect when the connection is lost. Example:
     ```python
     async def connect_with_retry():
         while True:
             try:
                 await connect()
             except websockets.exceptions.ConnectionClosed:
                 print("Connection lost, retrying in 5 seconds...")
                 await asyncio.sleep(5)
     ```

## 7. Complete Example Code (Python example implementing a weather query client)

Here's a complete Python example demonstrating how to implement a weather query service client:

```python
import asyncio
import websockets
import json
import requests
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Configuration
SOCKET_URL = "wss://socket.zhycit.com/"
TOKEN_URL = "https://ai.zhycit.com/wecom/token"
APPID = os.getenv("APPID")
SECRET = os.getenv("SECRET")
WEATHER_API_KEY = os.getenv("WEATHER_API_KEY")
CLIENT_ID = "/service/zhycit/weather_1"  # Weather service client ID

# Get token
def get_token():
    response = requests.get(f"{TOKEN_URL}?appid={APPID}&secret={SECRET}")
    data = response.json()
    return data["token"]

# Simulate weather query
def get_weather(city):
    # This should use an actual weather API, this is just a simulation
    return f"Weather in {city}: Sunny, 25°C"

# Handle received messages
async def handle_message(websocket, message):
    data = json.loads(message)
    if data["type"] == "msg" and "weather" in data["content"].lower():
        city = data["content"].split("weather")[0].strip()
        weather = get_weather(city)
        response = {
            "from": CLIENT_ID,
            "to": data["from"],  # Reply to the chatbot that sent the message
            "subject": "Weather Query Result",
            "content": weather,
            "type": "msg"
        }
        await websocket.send(json.dumps(response))

# Main WebSocket client logic
async def weather_service():
    token = get_token()
    uri = f"{SOCKET_URL}?clientID={CLIENT_ID}"
    headers = {"Authorization": f"Bearer {token}"}

    async with websockets.connect(uri, extra_headers=headers) as websocket:
        print(f"Connected to {uri}")

        while True:
            try:
                message = await websocket.recv()
                print(f"Received: {message}")
                await handle_message(websocket, message)
            except websockets.exceptions.ConnectionClosed:
                print("Connection closed. Reconnecting...")
                break

# Run the client
async def main():
    while True:
        try:
            await weather_service()
        except Exception as e:
            print(f"Error occurred: {e}")
        print("Reconnecting in 5 seconds...")
        await asyncio.sleep(5)

if __name__ == "__main__":
    asyncio.get_event_loop().run_until_complete(main())
```

This example code implements the following features:

1. Loads configuration information from environment variables.
2. Implements a function to obtain an authorization token.
3. Creates a simulated weather query function (in a real application, you should use an actual weather API).
4. Implements message handling logic, querying and replying with weather information when a message containing the "weather" keyword is received.
5. Establishes a WebSocket connection and automatically reconnects when the connection is lost.

Usage instructions:

1. Ensure you have installed the required Python libraries: `pip install websockets requests python-dotenv`
2. Create a `.env` file with the following content:
   ```
   APPID=your_appid
   SECRET=your_secret
   WEATHER_API_KEY=your_weather_api_key
   ```
3. Run the script: `python weather_service_client.py`

This client will connect to the WebSocket server, listen for weather query requests, and reply with weather information. Chatbots can send messages containing the "weather" keyword to this service, for example, "London weather", and the service will reply with the corresponding weather information.

Note: This example uses `/service/zhycit/weather_1` as the client ID. In a real application, you need to ensure that this ID is unique and matches your system design.

## 8. Security Considerations
   - Using encrypted connections
     Ensure the use of WSS (WebSocket Secure) for encrypted communication.

   - Implementing message verification
     Verify the integrity and source of received messages.

   - Protecting sensitive information
     Do not log sensitive information, such as authorization tokens, in logs.

## 9. Performance Optimization
   - Implementing heartbeat mechanism
     Regularly send heartbeat messages to keep the connection active.

   - Message queue
     Use message queues to manage message flow in high-load situations.

By following this guide, you should be able to successfully integrate AI chatbots into a WebSocket-based communication system. Remember to adjust implementation details according to the actual server API and requirements.
