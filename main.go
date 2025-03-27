package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"time"
)

type StdioMCPClientConfig struct {
	Command string   `json:"command"`
	Env     []string `json:"env"`
	Args    []string `json:"args"`
}

type SSEMCPClientConfig struct {
	BaseURL string            `json:"baseURL"`
	Headers map[string]string `json:"headers"`
	Timeout int64             `json:"timeout"`
}

type MCPClientType string

const (
	MCPClientTypeStdio MCPClientType = "stdio"
	MCPClientTypeSSE   MCPClientType = "sse"
)

type MCPClientConfig struct {
	Type   MCPClientType   `json:"type"`
	Config json.RawMessage `json:"config"`
}
type SSEServerConfig struct {
	BaseURL string `json:"baseURL"`
	Addr    string `json:"addr"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Config struct {
	Server  SSEServerConfig            `json:"server"`
	Clients map[string]MCPClientConfig `json:"clients"`
}

func main() {
	conf := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	config, err := loadConfig(*conf)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	start(config)
}

func start(config *Config) {
	mcpServer := server.NewMCPServer(
		config.Server.Name,
		config.Server.Addr,
		server.WithResourceCapabilities(true, true),
	)

	var errorGroup errgroup.Group
	var clients []client.MCPClient
	info := mcp.Implementation{
		Name:    config.Server.Name,
		Version: config.Server.Version,
	}
	for name, clientConfig := range config.Clients {
		log.Printf("Connecting to %s", name)
		mcpClient, err := newMCPClient(clientConfig)
		if err != nil {
			log.Fatalf("Failed to create MCP client: %v", err)
		}
		clients = append(clients, mcpClient)
		errorGroup.Go(func() error {
			return addClient(info, mcpClient, mcpServer)
		})
	}
	defer func() {
		for _, c := range clients {
			_ = c.Close()
		}
	}()
	err := errorGroup.Wait()
	if err != nil {
		log.Fatalf("Failed to add clients: %v", err)
	}
	sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL(config.Server.BaseURL))
	log.Printf("Starting SSE server")
	log.Printf("SSE server listening on %s", config.Server.Addr)
	err = sseServer.Start(config.Server.Addr)
	if err != nil {
		log.Fatalf("Failed to start SSE server: %v", err)
	}
}

func loadConfig(filePath string) (*Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func parseMCPClientConfig(conf MCPClientConfig) (any, error) {
	switch conf.Type {
	case MCPClientTypeStdio:
		var config StdioMCPClientConfig
		err := json.Unmarshal(conf.Config, &config)
		if err != nil {
			return nil, err
		}
		return config, nil
	case MCPClientTypeSSE:
		var config SSEMCPClientConfig
		err := json.Unmarshal(conf.Config, &config)
		if err != nil {
			return nil, err
		}
		return config, nil
	default:
		return nil, errors.New("invalid client type")
	}
}

func newMCPClient(conf MCPClientConfig) (client.MCPClient, error) {
	clientInfo, pErr := parseMCPClientConfig(conf)
	if pErr != nil {
		return nil, pErr
	}
	switch v := clientInfo.(type) {
	case StdioMCPClientConfig:
		return client.NewStdioMCPClient(v.Command, v.Env, v.Args...)
	case SSEMCPClientConfig:
		var options []client.ClientOption
		if v.Timeout > 0 {
			options = append(options, client.WithSSEReadTimeout(time.Duration(v.Timeout)*time.Second))
		}
		if len(v.Headers) > 0 {
			options = append(options, client.WithHeaders(v.Headers))
		}
		return client.NewSSEMCPClient(v.BaseURL, options...)
	}
	return nil, errors.New("invalid client type")
}

func addClient(clientInfo mcp.Implementation, mcpClient client.MCPClient, mcpServer *server.MCPServer) error {
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = clientInfo
	_, err := mcpClient.Initialize(context.Background(), initRequest)
	if err != nil {
		return err
	}
	log.Printf("Successfully initialized MCP client")
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := mcpClient.ListTools(context.Background(), toolsRequest)
	if err != nil {
		return err
	}
	log.Printf("Successfully listed %d tools", len(tools.Tools))
	for _, tool := range tools.Tools {
		log.Printf("Adding tool %s", tool.Name)
		mcpServer.AddTool(tool, mcpClient.CallTool)
	}
	return nil
}
