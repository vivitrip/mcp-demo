package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"io/ioutil"
	"log"
	"net/http"
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

	// 构造请求体
	reqBody, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 发送 HTTP 请求到 /greet
	resp, err := http.Post("http://localhost:8091/greet", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call /greet endpoint: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 response: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析响应体
	var greetResp struct {
		Message string `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&greetResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	// 返回包装后的结果
	return mcp.NewToolResultText(greetResp.Message), nil
}
