# 企业级AI聊天机器人接入指南（基于Socket客户端的接入）

## 1. 简介
   - 本指南的目的
     本指南旨在帮助开发者快速将企业级AI聊天机器人与WebSocket服务器集成。
   
   - 接入过程概述
     接入过程包括实现Socket客户端应用程序、建立安全连接，以及处理AI聊天机器人与服务器之间的消息交换。

## 2. 系统要求
   - 操作系统
     - Windows 10 或更高版本
     - macOS 10.12 或更高版本（不支持WeChat）
     - Linux（推荐使用Ubuntu 18.04 LTS或更高版本，不支持WeChat）
   
   - 配置要求
     - CPU: 双核处理器或更高
     - 内存: 至少4GB RAM
     - 存储: 至少10GB可用磁盘空间
     - 网络带宽: 建议至少10Mbps上行和下行速度
   
   - 网络环境
     - 稳定的互联网连接
     - 允许WebSocket通信的网络设置

## 3. 准备工作
   - 获取必要的凭证
     - WebSocket服务器的URL，企业版 `GO-POSTOFFICE` 服务器的URL为：`wss://socket.zhycit.com/`，用户可以自行搭建服务器，参与文档：[GO-POSTOFFICE 文档](https://github.com/zlz3907/GO-POSTOFFICE)
     - 授权令牌（企业版需要）

       使用以下API获取授权令牌：
       ```
       curl --location 'https://ai.zhycit.com/wecom/token?appid=<APPID>&secret=<SECRET>'
       ```
       
       返回结果示例：
       ```json
       {
           "expires": 1732256675685,
           "token": "af71********6922",
           "ttl": 2589694
       }
       ```
       
       其中，`token`字段即为授权令牌，`expires`为过期时间戳（截止时间），`ttl`为令牌剩余时间（单位：秒）。
       
       注意：请妥善保管您的appid和secret，不要泄露给他人。在实际使用时，请使用您自己的appid和secret。
     - 客户端ID（用于识别每个连接）

       客户端ID是每个连接的唯一标识符，用于区分不同的客户端连接。每一个客户端都可以提供对应的功能服务，例如天气查询、网络搜索、知识问答、数据分析等。

       客户端ID可以由用户自行定义，例如：`user_123456`，`user_654321`等。如果客户端提供特定的服务，可以定义为`/service/<company>/weather_1`，`/service/<company>/search_1`等。

   - 选择编程语言和WebSocket客户端库

     根据您的技术栈选择合适的编程语言和WebSocket客户端库。常见的选择包括：
     - Python: `websockets`或`socket.io-client`
     - Java: `Java-WebSocket`或`Tyrus`
     - JavaScript: 原生WebSocket API或`socket.io-client`

## 4. 实现Socket客户端
   - 建立WebSocket连接
     使用选定的WebSocket库连接到服务器。示例（使用Python）：
     ```python
     import websockets
     import asyncio

     async def connect():
         uri = "ws://your-server-url:port"
         async with websockets.connect(uri) as websocket:
             # 处理连接逻辑
     ```

   - 实现身份验证

     如果服务器要求身份验证，在建立连接时发送必要的凭证。示例：
     ```python
     headers = {
         "Authorization": "Bearer your-token-here",
         "Client-ID": "your-client-id"
     }
     async with websockets.connect(uri, extra_headers=headers) as websocket:
         # 连接后的逻辑
     ```

   - 消息发送和接收

     实现发送消息到服务器和接收服务器消息的功能。示例：
     ```python
     async def message_handler(websocket):
         while True:
             message = await websocket.recv()
             # 处理接收到的消息
             
             response = "AI生成的回复"
             await websocket.send(response)
     ```

## 5. 消息格式和协议
   - 定义消息结构

     企业版服务器互通消息的JSON结构规范（自建服务端程序可参考）：
     ```json
     {
       "from": "sender_id",
       "to": "recipient_id",
       "subject": "消息主题",
       "content": "用户输入或AI回复",
       "type": "msg",
       "createTime": 1623456789
     }
     ```

    协议参考：[GO-POSTOFFICE 消息协议](message-protocol_CN.md)

   - 实现消息序列化和反序列化
     使用JSON库处理消息的编码和解码。

## 6. 错误处理和重连机制
   - 实现错误捕获
     捕获并处理可能出现的网络错误或服务器响应问题。

   - 添加重连逻辑
     在连接断开时自动尝试重新连接。示例：
     ```python
     async def connect_with_retry():
         while True:
             try:
                 await connect()
             except websockets.exceptions.ConnectionClosed:
                 print("连接断开，5秒后重试...")
                 await asyncio.sleep(5)
     ```

## 7. 完整的示例代码（以Python为例，实现一个天气查询的客户端）

以下是一个完整的Python示例，展示了如何实现一个天气查询服务客户端：

```python
import asyncio
import websockets
import json
import requests
import os
from dotenv import load_dotenv

# 加载环境变量
load_dotenv()

# 配置
SOCKET_URL = "wss://socket.zhycit.com/"
TOKEN_URL = "https://ai.zhycit.com/wecom/token"
APPID = os.getenv("APPID")
SECRET = os.getenv("SECRET")
WEATHER_API_KEY = os.getenv("WEATHER_API_KEY")
CLIENT_ID = "/service/zhycit/weather_1"  # 天气服务客户端ID

# 获取token
def get_token():
    response = requests.get(f"{TOKEN_URL}?appid={APPID}&secret={SECRET}")
    data = response.json()
    return data["token"]

# 模拟天气查询
def get_weather(city):
    # 这里应该使用实际的天气API，这里只是一个模拟
    return f"{city}的天气：晴天，温度25°C"

# 处理接收到的消息
async def handle_message(websocket, message):
    data = json.loads(message)
    if data["type"] == "msg" and "天气" in data["content"]:
        city = data["content"].split("天气")[0]
        weather = get_weather(city)
        response = {
            "from": CLIENT_ID,
            "to": data["from"],  # 回复给发送消息的聊天机器人
            "subject": "天气查询结果",
            "content": weather,
            "type": "msg"
        }
        await websocket.send(json.dumps(response))

# 主要的WebSocket客户端逻辑
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

# 运行客户端
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

这个示例代码实现了以下功能：

1. 从环境变量加载配置信息。
2. 实现了获取授权token的函数。
3. 创建了一个模拟的天气查询函数（在实际应用中，您应该使用真实的天气API）。
4. 实现了消息处理逻辑，当接收到包含"天气"关键词的消息时，会查询相应城市的天气并回复。
5. 建立了WebSocket连接，并在连接断开时自动重连。

使用说明：

1. 确保已安装所需的Python库：`pip install websockets requests python-dotenv`
2. 创建一个`.env`文件，包含以下内容：
   ```
   APPID=your_appid
   SECRET=your_secret
   WEATHER_API_KEY=your_weather_api_key
   ```
3. 运行脚本：`python weather_service_client.py`

这个客户端将连接到WebSocket服务器，监听天气查询请求，并回复天气信息。聊天机器人可以发送包含"天气"关键词的消息到这个服务，例如"北京天气"，服务将回复相应的天气信息。

注意：这个示例使用了`/service/zhycit/weather_1`作为客户端ID。在实际应用中，您需要确保这个ID是唯一的，并且与您的系统设计相匹配。




## 8. 安全考虑
   - 使用加密连接
     确保使用WSS（WebSocket Secure）进行加密通信。

   - 实施消息验证
     验证接收到的消息的完整性和来源。

   - 保护敏感信息
     不要在日志中记录敏感信息，如授权令牌。

## 9. 性能优化
   - 实现心跳机制
     定期发送心跳消息以保持连接活跃。

   - 消息队列
     在高负载情况下使用消息队列来管理消息流。

通过遵循这个指南，您应该能够成功地将AI聊天机器人集成到基于WebSocket的通信系统中。记得根据实际的服务器API和需求调整实现细节。
