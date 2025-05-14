package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGreetHandler(t *testing.T) {
	// 启动测试服务器
	server := httptest.NewServer(http.HandlerFunc(greetHandler))
	defer server.Close() // 确保服务器在测试结束时关闭

	// 测试正常请求
	t.Run("Valid Request", func(t *testing.T) {
		reqBody := `{"name":"Alice"}`
		resp, err := http.Post(server.URL+"/greet", "application/json", bytes.NewBufferString(reqBody))
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var respBody GreetingResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		expectedMessage := "Hello, Alice!"
		if respBody.Message != expectedMessage {
			t.Errorf("expected message %q, got %q", expectedMessage, respBody.Message)
		}
	})

	// 测试非 POST 请求
	t.Run("Invalid Method", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/greet", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", resp.StatusCode)
		}
	})

	// 测试无效请求体
	t.Run("Invalid Body", func(t *testing.T) {
		reqBody := `{"invalid_json"}`
		resp, err := http.Post(server.URL+"/greet", "application/json", bytes.NewBufferString(reqBody))
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", resp.StatusCode)
		}
	})
}
