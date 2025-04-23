package main

import (
	"encoding/json"
	"time"
)

type StdioMCPClientConfig struct {
	Command string            `json:"command"`
	Env     map[string]string `json:"env"`
	Args    []string          `json:"args"`
}

type SSEMCPClientConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type StreamableMCPClientConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Timeout time.Duration     `json:"timeout"`
}

type MCPClientType string

const (
	MCPClientTypeStdio      MCPClientType = "stdio"
	MCPClientTypeSSE        MCPClientType = "sse"
	MCPClientTypeStreamable MCPClientType = "streamable"
)

type MCPClientConfig struct {
	Type           MCPClientType   `json:"type"`
	Config         json.RawMessage `json:"config"`
	PanicIfInvalid bool            `json:"panicIfInvalid"`
	LogEnabled     bool            `json:"logEnabled"`
	AuthTokens     []string        `json:"authTokens"`
}

type SSEServerConfig struct {
	BaseURL          string   `json:"baseURL"`
	Addr             string   `json:"addr"`
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	GlobalAuthTokens []string `json:"globalAuthTokens"`
}

type Config struct {
	Server  SSEServerConfig            `json:"server"`
	Clients map[string]MCPClientConfig `json:"clients"`
}
