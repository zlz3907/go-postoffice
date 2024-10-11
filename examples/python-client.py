import asyncio
import websockets

async def hello():
    uri = "ws://localhost:7502/ws"
    async with websockets.connect(uri) as websocket:
        await websocket.send("Hello, server!")
        print(f"Sent: Hello, server!")

        while True:
            try:
                message = await websocket.recv()
                print(f"Received: {message}")
            except websockets.exceptions.ConnectionClosed:
                print("Connection closed")
                break

asyncio.get_event_loop().run_until_complete(hello())