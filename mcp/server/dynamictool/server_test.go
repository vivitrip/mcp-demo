package main

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestDynamicToolServer(t *testing.T) {
	//// 启动 server.go 的 MCP 服务器
	cmd := exec.Command("go", "run", "server.go")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// 等待服务器启动
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
	t.Log("客户端启动成功")

	t.Log("初始化 MCP 客户端...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Client Demo",
		Version: "1.0.0",
	}

	// 初始化 MCP 客户端并连接到服务器
	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		t.Fatalf("Failed to initialize MCP client: %v", err)
	}
	t.Logf("初始化成功，服务器信息: %s %s", initResult.ServerInfo.Name, initResult.ServerInfo.Version)

	// 获取并打印 MCP Server 所有可用的工具
	toolsResult, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}
	for _, tool := range toolsResult.Tools {
		t.Logf("可用工具: %s - %s", tool.Name, tool.Description)
	}

	// 测试 hello_tool
	t.Log("调用工具: hello_world")
	helloRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	helloRequest.Params.Name = "hello_world"
	helloRequest.Params.Arguments = map[string]any{
		"name": "Vi_error",
	}

	helloResult, err := mcpClient.CallTool(ctx, helloRequest)
	if err != nil {
		t.Fatalf("Failed to call hello_world tool: %v", err)
	}

	expectedHello := "Hello, Vi_error!, This is from your go mcp server"
	actualHello := helloResult.Content[0].(mcp.TextContent).Text
	if actualHello != expectedHello {
		t.Errorf("Unexpected hello_world result: got %q, want %q", actualHello, expectedHello)
	}
	t.Logf("hello_world 工具响应: %s", actualHello)

	// 测试 add_tool
	t.Log("调用工具: add_tool")
	addToolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	addToolRequest.Params.Name = "add_tool"
	addToolRequest.Params.Arguments = map[string]any{
		"toolName": "dynamic_tool",
		"toolDesc": "A dynamically added tool",
		"paramList": `{
			"param1": true,
			"param2": false
		}`,
	}

	addToolResult, err := mcpClient.CallTool(ctx, addToolRequest)
	if err != nil {
		t.Fatalf("Failed to call add_tool: %v", err)
	}
	t.Logf("add_tool 工具响应: %s", addToolResult.Content[0].(mcp.TextContent).Text)

	// 测试动态添加的工具
	t.Log("调用动态添加的工具: dynamic_tool")
	dynamicToolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	dynamicToolRequest.Params.Name = "dynamic_tool"
	dynamicToolRequest.Params.Arguments = map[string]any{
		"param1": "value1",
	}

	dynamicToolResult, err := mcpClient.CallTool(ctx, dynamicToolRequest)
	if err != nil {
		t.Fatalf("Failed to call dynamic_tool: %v", err)
	}
	t.Logf("dynamic_tool 工具响应: %s", dynamicToolResult.Content[0].(mcp.TextContent).Text)

	// 测试 delete_tool
	t.Log("调用工具: delete_tool")
	deleteToolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	deleteToolRequest.Params.Name = "delete_tool"
	deleteToolRequest.Params.Arguments = map[string]any{
		"toolName": "dynamic_tool",
	}

	deleteToolResult, err := mcpClient.CallTool(ctx, deleteToolRequest)
	if err != nil {
		t.Fatalf("Failed to call delete_tool: %v", err)
	}
	t.Logf("delete_tool 工具响应: %s", deleteToolResult.Content[0].(mcp.TextContent).Text)

	// 验证删除的工具是否已被移除
	t.Log("验证动态工具是否已被删除")
	_, err = mcpClient.CallTool(ctx, dynamicToolRequest)
	if err == nil {
		t.Fatalf("Expected error when calling deleted tool, but got none")
	}
	t.Logf("动态工具已成功删除: %v", err)

	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			t.Logf("Failed to kill server process: %v", err)
		}
	}()
}
