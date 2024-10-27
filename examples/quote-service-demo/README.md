# Quote Service Demo

This demo showcases a simple quote service implementation using the GO-POSTOFFICE WebSocket server.

![Global Architecture Diagram](../../docs/imgs/global_architecture_diagram_en.png)

## Overview

The Quote Service Demo demonstrates how to create a basic service that provides random quotes to clients connected via WebSocket. This example illustrates the core concepts of using GO-POSTOFFICE for real-time communication.

## Features

- WebSocket-based communication
- Random quote generation and delivery
- Simple client-server interaction

## Prerequisites

- Go 1.23.1 or higher
- GO-POSTOFFICE server running

## Setup and Running

1. Ensure the GO-POSTOFFICE server is running. If not, start it by following the instructions in the main README.

2. Navigate to the quote service demo directory:
   ```
   cd examples/quote-service-demo
   ```

3. Install dependencies and run the quote service:
   ```
   yarn install
   yarn dev
   ```

4. The service will start and connect to the GO-POSTOFFICE server.

## How It Works

1. The quote service connects to the GO-POSTOFFICE server as a client.
2. It registers itself to receive messages addressed to "/service/quote".
3. When a client requests a quote, it sends a message to "/service/quote".
4. The quote service receives the request, generates a random quote, and sends it back to the client.

## Testing

You can test the quote service using a WebSocket client or by implementing a simple client application that connects to the GO-POSTOFFICE server and requests quotes.

## Customization

Feel free to modify the `quotes` slice in the code to add your own quotes or implement a different quote generation mechanism.

## Next Steps

- Implement error handling and logging
- Add more sophisticated quote selection logic
- Create a front-end client to interact with the quote service

This demo serves as a starting point for building more complex services using GO-POSTOFFICE. Explore the main documentation for advanced features and best practices.
