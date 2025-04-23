module github.com/TBXark/mcp-proxy

go 1.23.0

toolchain go1.23.7

//mcp-go 0.23.0 has a bug that causes the mcp-go to be unable to connect to the stdio mcp server.
//replace github.com/mark3labs/mcp-go v0.23.0 => github.com/tbxark-fork/mcp-go v0.0.0-20250423055820-6f258c8d2d40

require (
	github.com/TBXark/confstore v0.0.0-20250312091006-41b7721fb8c8
	github.com/mark3labs/mcp-go v0.22.0
	golang.org/x/sync v0.13.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)
