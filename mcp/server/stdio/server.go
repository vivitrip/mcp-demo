package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"MCP Server with Stdio",
		"1.0.0",
	)

	// Add greetTool
	greetTool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	// Add greetTool handler
	s.AddTool(greetTool, helloHandler)

	//Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}
	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!, This is from your go mcp server", name)), nil
}
