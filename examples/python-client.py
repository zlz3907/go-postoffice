import asyncio
import websockets
import json

async def hello():
    clientID = "python-client-001"  # 设置客户端ID
    uri = f"ws://localhost:7502/?clientID={clientID}"
    headers = {
        "Authorization": "Bearer your_token_here"
    }
    async with websockets.connect(uri, extra_headers=headers) as websocket:
        while True:
            message = {
                "from": "python-client",
                "to": "server",
                "subject": "Hello",
                "content": "How are you?",
                "type": "msg"
            }
            await websocket.send(json.dumps(message))
            print(f"Sent: {message}")

            try:
                response = await asyncio.wait_for(websocket.recv(), timeout=5.0)
                print(f"Received: {response}")
            except asyncio.TimeoutError:
                print("No response received within 5 seconds")

            await asyncio.sleep(5)

asyncio.get_event_loop().run_until_complete(hello())
