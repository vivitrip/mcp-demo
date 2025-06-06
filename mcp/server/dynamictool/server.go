package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
)

var s *server.MCPServer // 将变量s提升为全局变量

func main() {
	// 创建 MCP server
	s = server.NewMCPServer(
		"Demo one",
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

	addTool := mcp.NewTool("add_tool", mcp.WithDescription("add_tool"))
	s.AddTool(addTool, addToolHandler)

	deleteTool := mcp.NewTool("delete_tool", mcp.WithDescription("delete_tool"))
	s.AddTool(deleteTool, deleteToolHandler)

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
	name, ok := request.Params.Arguments.(map[string]any)["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}
	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!, This is from your go mcp server", name)), nil
}

// addToolHandler 方法，支持新增参数 paramList，并封装动态生成的工具逻辑
func addToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	toolName, ok := request.Params.Arguments.(map[string]any)["toolName"].(string)
	if !ok {
		return nil, errors.New("toolName must be a string")
	}

	toolDesc, ok := request.Params.Arguments.(map[string]any)["toolDesc"].(string)
	if !ok {
		return nil, errors.New("toolDesc must be a string")
	}

	paramListStr, ok := request.Params.Arguments.(map[string]any)["paramList"].(string)
	if !ok {
		return nil, errors.New("paramList must be a JSON string")
	}

	var paramList map[string]bool
	if err := json.Unmarshal([]byte(paramListStr), &paramList); err != nil {
		return nil, fmt.Errorf("failed to parse paramList: %v", err)
	}

	// 定义参数
	toolOptions := []mcp.ToolOption{mcp.WithDescription(toolDesc)}
	for paramName, isRequired := range paramList {
		if isRequired {
			toolOptions = append(toolOptions, mcp.WithString(paramName, mcp.Required()))
		} else {
			toolOptions = append(toolOptions, mcp.WithString(paramName))
		}
	}

	// 创建新工具
	newTool := mcp.NewTool(toolName, toolOptions...)
	// 注册工具
	s.AddTool(newTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 校验必填参数
		for paramName, isRequired := range paramList {
			if isRequired {
				argsMap, ok := req.Params.Arguments.(map[string]any)
				if !ok {
					return nil, errors.New("invalid arguments format")
				}
				if _, exists := argsMap[paramName]; !exists {
					return nil, fmt.Errorf("missing required parameter: %s", paramName)
				}
			}
		}
		// 增加你需要的逻辑
		response := "this answer is from your go mcp server"
		return mcp.NewToolResultText(fmt.Sprintf("Tool %s executed successfully! Response: %s", toolName, response)), nil
	})
	return mcp.NewToolResultText(fmt.Sprintf("Tool %s with description '%s' and parameters added successfully!", toolName, toolDesc)), nil
}

func deleteToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	toolName, ok := request.Params.Arguments.(map[string]any)["toolName"].(string)
	if !ok {
		return nil, errors.New("toolName must be a string")
	}
	// 从 MCPServer 中删除工具
	s.DeleteTools(toolName)

	return mcp.NewToolResultText(fmt.Sprintf("Tool %s removed successfully!", toolName)), nil
}
