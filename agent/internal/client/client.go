package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient 封装了与后端 API 的交互
type APIClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewAPIClient 创建一个新的 API 客户端
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		httpClient: &http.Client{Timeout: 40 * time.Second}, // 超时要比长轮询的时间长一点
		baseURL:    baseURL,
	}
}

// Register 注册 Agent 到后端
func (c *APIClient) Register(hostname, ip, os string) (string, error) {
	// 这个结构体应该与后端 api/types.go 中的 RegisterAgentRequest 一致
	reqBody, _ := json.Marshal(map[string]string{
		"hostname":   hostname,
		"ip_address": ip,
		"os":         os,
	})

	resp, err := c.httpClient.Post(c.baseURL+"/api/v1/agent/register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	// 这个结构体应该与后端 api/types.go 中的 RegisterAgentResponse 一致
	var respBody struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", err
	}
	return respBody.AgentID, nil
}

// SendHeartbeat 发送心跳
func (c *APIClient) SendHeartbeat(agentID string) error {
	reqBody, _ := json.Marshal(map[string]string{"agent_id": agentID})
	resp, err := c.httpClient.Post(c.baseURL+"/api/v1/agent/heartbeat", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat failed with status: %s", resp.Status)
	}
	return nil
}

// FetchTasks 通过长轮询获取任务
func (c *APIClient) FetchTasks(ctx context.Context, agentID string) (*Task, error) {
	url := fmt.Sprintf("%s/api/v1/agent/tasks?agent_id=%s", c.baseURL, agentID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil // 超时，没有任务
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch tasks failed with status: %s", resp.Status)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

// PostResult 上报任务结果
func (c *APIClient) PostResult(result TaskResult) error {
	reqBody, _ := json.Marshal(result)
	resp, err := c.httpClient.Post(c.baseURL+"/api/v1/agent/tasks/results", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("post result failed with status: %s, body: %s", resp.Status, string(body))
	}
	return nil
}

// --- 在 client 包内部定义与后端 API 对应的类型 ---
// 这些结构体应该与后端 internal/core/engine/model.go 中的 Task 和 TaskResult 结构一致
type Task struct {
	ID         string `json:"ID"`
	WorkflowID string `json:"WorkflowID"`
	Type       string `json:"Type"`
	Command    string `json:"Command"`
}

type TaskResult struct {
	TaskID   string `json:"task_id"`
	AgentID  string `json:"agent_id"`
	Success  bool   `json:"success"`
	Output   string `json:"output"`
	Error    string `json:"error"`
	ExitCode int    `json:"exit_code"`
}
