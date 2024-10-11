# GO-POSTOFFICE Configuration

[中文](configuration_CN.md) | English

This document describes the configuration options for the GO-POSTOFFICE project.

## Configuration File

The configuration file is located in the `.env` directory. There are two main configuration files:

- `config-dev.json`: Used for development environment
- `config-zhycit.json`: Used for production environment

## Configuration Options

Here are the main configuration options:

### WebSocket Server

- `socketPort`: The port number on which the WebSocket server will listen. 
  - Type: integer
  - Default: 7502

- `maxWebSocketConnections`: The maximum number of concurrent WebSocket connections allowed.
  - Type: integer
  - Default: 20000 (for production)

### SSL Configuration (Optional)

- `sslPort`: The port number for SSL connections.
  - Type: integer
  - Default: 7503

- `sslCertPath`: The path to the SSL certificate file.
  - Type: string
  - Example: "./cert/cert.pem"

- `sslKeyPath`: The path to the SSL key file.
  - Type: string
  - Example: "./cert/key.pem"

### Data Source

- `dataSource`: Configuration for external data sources.
  - Type: object

  #### Redis Configuration

  - `redis`: Configuration for Redis connection.
    - Type: object
    - Properties:
      - `gnas-ai`: The configuration for a specific Redis instance.
        - Type: object
        - Properties:
          - `uri`: The URI for connecting to the Redis instance.
            - Type: string
            - Example: "127.0.0.1:6379" (for development)

## Example Configuration

Here's an example of a complete configuration file:

```json
{
  "dataSource": {
    "redis": {
      "gnas-ai": {
        "uri": "127.0.0.1:6379"
      }
    }
  },
  "socketPort": 7502,
  "sslPort": 7503,
  "sslCertPath": "./cert/cert.pem",
  "sslKeyPath": "./cert/key.pem",
  "maxWebSocketConnections": 2
}
```

## Changing Configuration

To change the configuration:

1. Copy the appropriate configuration file (e.g., `config-dev.json` for development) to `.env/config-zhycit.json`.
2. Modify the values in `config-zhycit.json` as needed.
3. Restart the GO-POSTOFFICE server for the changes to take effect.

Remember to never commit sensitive information (like production database credentials) to version control. Use environment variables or secure secret management systems for sensitive data in production environments.