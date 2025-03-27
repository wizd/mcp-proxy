# MCP Proxy Server

An MCP proxy server that aggregates and serves multiple MCP resource servers through a single interface.

## Features

- **Proxy Multiple MCP Clients**: Connects to multiple MCP resource servers and aggregates their tools and capabilities.
- **SSE Support**: Provides an SSE (Server-Sent Events) server for real-time updates.
- **Flexible Configuration**: Supports multiple client types (`stdio` and `sse`) with customizable settings.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/TBXark/mcp-proxy.git
   cd mcp-proxy
   ```

2. Build the project:
   ```bash
   go build -o mcp-proxy main.go
   ```

## Configuration

The server is configured using a JSON file. Below is an example configuration:

```json
{
  "server": {
    "baseURL": "http://localhost:8080",
    "addr": ":8080",
    "name": "MCP Proxy",
    "version": "1.0.0"
  },
  "clients": {
    "fetch": {
      "type": "stdio",
      "config": {
        "command": "uvx",
        "env": [],
        "args": ["mcp-server-fetch"]
      }
    },
    "amap": {
      "type": "sse",
      "config": {
        "baseURL": "https://router.mcp.so/sse/xxxxx",
      }
    }
  }
}
```

- **Server Configuration**:
  - `baseURL`: The base URL for the SSE server.
  - `addr`: The address the server listens on.
  - `name`: The name of the server.
  - `version`: The version of the server.

- **Clients Configuration**:
  - `type`: The type of the client (`stdio` or `sse`).
  - `config`: The specific configuration for the client type.

## Usage

1. Run the server:
   ```bash
   ./mcp-proxy -config path/to/config.json
   ```

2. The server will start and aggregate the tools and capabilities of the configured MCP clients.


### Thanks

This project was inspired by the [adamwattis/mcp-proxy-server](https://github.com/adamwattis/mcp-proxy-server) project



## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.