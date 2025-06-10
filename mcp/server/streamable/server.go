package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"net/http"
)

func main() {
	// 以下二选一
	startMcpStateful()
	startMcpStateless()
}

func startMcpStateful() {
	mcpServer := server.NewMCPServer(
		"Stateful MCP Server with StreamableHTTP",
		"1.0.0",
	)
	greetTool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
	mcpServer.AddTool(greetTool, helloHandler)

	mux := http.NewServeMux()

	statefulServer := server.NewStreamableHTTPServer(mcpServer, server.WithEndpointPath("/mcp/stateful"))
	mux.Handle("/mcp/stateful", statefulServer)
	addr := ":8080"
	fmt.Printf("Starting MCP HTTP server at %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func startMcpStateless() {
	mcpServer := server.NewMCPServer(
		"Stateful MCP Server with StreamableHTTP",
		"1.0.0",
	)
	greetTool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
	mcpServer.AddTool(greetTool, helloHandler)

	mux := http.NewServeMux()

	statelessServer := server.NewStreamableHTTPServer(mcpServer, server.WithEndpointPath("/mcp/stateless"), server.WithStateLess(true))
	mux.Handle("/mcp/stateless", statelessServer)
	addr := ":8080"
	fmt.Printf("Starting MCP HTTP server at %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments.(map[string]any)["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}
	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!, This is from your go mcp server", name)), nil
}
