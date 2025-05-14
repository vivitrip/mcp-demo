package main

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestMCPServer(t *testing.T) {
	// 动态编译 server.go
	serverExecutable := "./server"
	buildCmd := exec.Command("go", "build", "server.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build server: %v", err)
	}

	// 创建一个基于 stdio 的MCP客户端
	mcpClient, err := client.NewStdioMCPClient(
		serverExecutable, // 使用动态编译后的可执行文件路径
		[]string{},
	)
	if err != nil {
		t.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Log("初始化 mcp 客户端...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Client Demo",
		Version: "1.0.0",
	}

	// 初始化MCP客户端并连接到服务器
	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		t.Fatalf("Failed to initialize MCP client: %v", err)
	}
	t.Logf("初始化成功，服务器信息: %s %s", initResult.ServerInfo.Name, initResult.ServerInfo.Version)

	// 调用工具
	t.Log("调用工具: hello_world")
	toolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	toolRequest.Params.Name = "hello_world"
	toolRequest.Params.Arguments = map[string]any{
		"name": "Vi_error",
	}

	// 调用工具并验证结果
	result, err := mcpClient.CallTool(ctx, toolRequest)
	if err != nil {
		t.Fatalf("Failed to call tool: %v", err)
	}

	expected := "Hello, Vi_error!, This is from your go mcp server"
	actual := result.Content[0].(mcp.TextContent).Text
	if actual != expected {
		t.Errorf("Unexpected tool result: got %q, want %q", actual, expected)
	}
}
