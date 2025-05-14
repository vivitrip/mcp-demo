package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
)

func main() {
	s := server.NewMCPServer(
		"MCP Server with SSE",
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

	//Start the sse server
	port := ":8090"
	baseUrl := "http://localhost" + port + "/"
	log.Printf("baseUrl is : %s", baseUrl)
	sseServer := server.NewSSEServer(s, server.WithBaseURL(baseUrl))
	log.Printf("SSE server listening on : %s", port)
	if err := sseServer.Start(port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}
	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!, This is from your go mcp server", name)), nil
}
