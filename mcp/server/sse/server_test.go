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
	// 启动 server.go 的 MCP 服务器
	cmd := exec.Command("go", "run", "server.go")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	//等待服务器启动
	time.Sleep(2 * time.Second)

	serverURL := "http://localhost:8090/sse"
	t.Logf("测试连接到 SSE 服务器: %s", serverURL)

	// 创建一个基于 SSE 的 MCP 客户端
	mcpClient, err := client.NewSSEMCPClient(serverURL)
	if err != nil {
		t.Fatalf("Failed to create SSE MCP client: %v", err)
	}
	defer mcpClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Log("start mcp 客户端...")
	if err := mcpClient.Start(ctx); err != nil {
		t.Fatalf("Failed to start mcpClient: %v", err)
	}
	t.Log("Client started")

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
	t.Logf("MCP Server Response is : %s", actual)

	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			t.Logf("Failed to kill server process: %v", err)
		}
	}()
}
