package main

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestAdapter(t *testing.T) {
	// 启动 rest/server.go 的 HTTP 服务器
	restCmd := exec.Command("go", "run", "rest/server.go")
	if err := restCmd.Start(); err != nil {
		t.Fatalf("Failed to start REST server: %v", err)
	}
	defer func() {
		if err := restCmd.Process.Kill(); err != nil {
			t.Logf("Failed to kill REST server process: %v", err)
		}
		// 等待进程完全退出
		_, _ = restCmd.Process.Wait()
	}()

	// 等待 REST 服务器启动
	time.Sleep(2 * time.Second)

	// 启动 adapter.go 的 MCP 服务器
	cmd := exec.Command("go", "run", "adapter.go")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start adapter server: %v", err)
	}
	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			t.Logf("Failed to kill adapter server process: %v", err)
		}
		// 等待进程完全退出
		_, _ = cmd.Process.Wait()
	}()

	// 等待 MCP 服务器启动
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

	t.Log("启动 MCP 客户端...")
	if err := mcpClient.Start(ctx); err != nil {
		t.Fatalf("Failed to start mcpClient: %v", err)
	}
	t.Log("客户端已启动")

	t.Log("初始化 MCP 客户端...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Adapter Test Client",
		Version: "1.0.0",
	}

	// 初始化 MCP 客户端并连接到服务器
	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		t.Fatalf("Failed to initialize MCP client: %v", err)
	}
	t.Logf("初始化成功，服务器信息: %s %s", initResult.ServerInfo.Name, initResult.ServerInfo.Version)

	// 调用���具
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

	expected := "Hello, Vi_error! This is from your rest server"
	actual := result.Content[0].(mcp.TextContent).Text
	if actual != expected {
		t.Errorf("Unexpected tool result: got %q, want %q", actual, expected)
	}
	t.Logf("MCP Server Response is : %s", actual)
}
