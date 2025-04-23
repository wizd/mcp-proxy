# MCP Proxy Server

An MCP proxy server that aggregates and serves multiple MCP resource servers through a single HTTP server.

## Features

- **Proxy Multiple MCP Clients**: Connects to multiple MCP resource servers and aggregates their tools and capabilities.
- **SSE Support**: Provides an SSE (Server-Sent Events) server for real-time updates.
- **Flexible Configuration**: Supports multiple client types (`stdio` and `sse`) with customizable settings.

## Installation

### Build from Source

 ```bash
git clone https://github.com/TBXark/mcp-proxy.git
cd mcp-proxy
go build -o mcp-proxy main.go
./mcp-proxy --config path/to/config.json
```

### Install by go

```bash
go install github.com/TBXark/mcp-proxy@latest
````

### Docker

> The Docker image supports two MCP calling methods by default: `npx` and `uvx`.
```bash
docker run -d -p 9090:9090 -v /path/to/config.json:/config/config.json ghcr.io/tbxark/mcp-proxy:latest
# or 
docker run -d -p 9090:9090 ghcr.io/tbxark/mcp-proxy:latest --config https://example.com/path/to/config.json
```

## Configuration

The server is configured using a JSON file. Below is an example configuration:
> This is the format for the new version's configuration. The old version's configuration will be automatically converted to the new format's configuration when it is loaded.

```jsonc
{
  "mcpProxy": {
    "baseURL": "http://localhost:9090",
    "addr": ":9090",
    "name": "MCP Proxy",
    "version": "1.0.0",
    "options": {
      "panicIfInvalid": false,
      "logEnabled": false,
      "authTokens": [
        "AdminToken"
      ]
    }
  },
  "mcpServers": {
    "fetch": {
      "command": "npx",
      "args": [
        "-y",
        "fetch-mcp"
      ],
      "env": {
      },
      "options": {
        "panicIfInvalid": true,
        "logEnabled": true,
        "authTokens": [
          "HelloWorld"
        ]
      }
    },
    "exampleServer": {
      "url": "https://example.com/mcp-sse",
      "headers":  {
        "Authorization": "Bearer example-token"
      }
    }
  }
}
```

### **`options`**
Common options for `mcpProxy` and `mcpServers`.

- `panicIfInvalid`: If true, the server will panic if the client is invalid.
- `logEnabled`: If true, the server will log the client's requests.
- `authTokens`: A list of authentication tokens for the client. The `Authorization` header will be checked against this list. 

> In the new configuration, the `authTokens` of `mcpProxy` is not a global authentication token, but rather the default authentication token for `mcpProxy`. When `authTokens` is set in `mcpServers`, the value of `authTokens` in `mcpServers` will be used instead of the value in `mcpProxy`. In other words, the `authTokens` of `mcpProxy` serves as a default value and is only applied when `authTokens` is not set in `mcpServers`.
> Other fields are the same.

### **`mcpProxy`**
Proxy HTTP server configuration
- `baseURL`: The public accessible URL of the server. This is used to generate the URLs for the clients.
- `addr`: The address the server listens on.
- `name`: The name of the server.
- `version`: The version of the server.
- `options`: Global options for the server, When `options.authTokens` is set, It will be the global authentication token for all clients. 

### **`mcpServers`**
MCP server configuration, Adopt the same configuration format as other MCP Clients.

For stdio mcp servers, the `command` field is required.
- `command`: The command to run the MCP client.
- `args`: The arguments to pass to the command.
- `env`: The environment variables to set for the command.
- `options`: Options specific to the client.

For sse mcp servers, the `url` field is required. 
- `url`: The URL of the MCP client.
- `headers`: The headers to send with the request to the MCP client.

For http streaming mcp servers, Not supported yet.

## Usage

```
Usage of mcp-proxy:
  -config string
        path to config file or a http(s) url (default "config.json")
  -help
        print help and exit
  -version
        print version and exit
```
1. The server will start and aggregate the tools and capabilities of the configured MCP clients.
2. You can access the server at `http(s)://{baseURL}/{clientName}/sse`. (e.g., `https://my-mcp.example.com/fetch/sse`, based on the example configuration)
3. If your MCP client does not support custom request headers., you can change the key in `clients` such as `fetch` to `fetch/{apiKey}`, and then access it via `fetch/{apiKey}`.

## Thanks

This project was inspired by the [adamwattis/mcp-proxy-server](https://github.com/adamwattis/mcp-proxy-server) project

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.