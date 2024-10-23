import asyncio
import websockets
import json
import random
import string

# 生成随机字符串的函数
def random_string(length):
    return ''.join(random.choice(string.ascii_letters) for _ in range(length))

# 创建消息的函数
def create_message(msg_type, to, content):
    return {
        "from": "python-client",
        "to": to,
        "subject": f"Test message - {msg_type}",
        "content": content,
        "type": msg_type
    }

async def test_client():
    clientID = "python-test-client-001"  # 设置客户端ID
    uri = f"ws://{host}:{port}/?clientID={clientID}"
    headers = {
        "Authorization": "Bearer your_token_here"
    }
    async with websockets.connect(uri, extra_headers=headers) as websocket:
        print(f"Connected to {uri}")

        # 发送登录消息
        login_msg = create_message("login", "server", {"token": "your_token_here"})
        await websocket.send(json.dumps(login_msg))
        print(f"Sent login message: {login_msg}")

        # 接收登录响应
        response = await websocket.recv()
        print(f"Received login response: {response}")

        # 发送不同类型的消息
        message_types = ["msg", "log", "heartbeat"]
        for _ in range(5):  # 发送5次消息
            msg_type = random.choice(message_types)
            content = random_string(10)  # 生成随机内容
            message = create_message(msg_type, "server", content)
            
            await websocket.send(json.dumps(message))
            print(f"Sent message: {message}")

            # 接收响应
            try:
                response = await asyncio.wait_for(websocket.recv(), timeout=5.0)
                print(f"Received response: {response}")
            except asyncio.TimeoutError:
                print("No response received within 5 seconds")

            await asyncio.sleep(1)  # 等待1秒再发送下一条消息

        # 发送退出消息
        logout_msg = create_message("logout", "server", "Logging out")
        await websocket.send(json.dumps(logout_msg))
        print(f"Sent logout message: {logout_msg}")

        # 接收退出响应
        response = await websocket.recv()
        print(f"Received logout response: {response}")

if __name__ == "__main__":
    asyncio.get_event_loop().run_until_complete(test_client())
