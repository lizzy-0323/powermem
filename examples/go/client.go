package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) (*APIResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp APIResponse
	decoder := json.NewDecoder(bytes.NewReader(respBody))
	decoder.UseNumber()
	if err := decoder.Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &apiResp, fmt.Errorf("API error (status %d): %s", resp.StatusCode, apiResp.Message)
	}

	return &apiResp, nil
}

func (c *Client) CreateMemory(req CreateMemoryRequest) ([]Memory, error) {
	resp, err := c.doRequest("POST", "/api/v1/memories", req)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var memories []Memory
	if err := json.Unmarshal(dataBytes, &memories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memories: %w", err)
	}

	return memories, nil
}

func (c *Client) ListMemories(req ListMemoriesRequest) (*MemoryListData, error) {
	params := url.Values{}
	if req.UserID != "" {
		params.Add("user_id", req.UserID)
	}
	if req.AgentID != "" {
		params.Add("agent_id", req.AgentID)
	}
	if req.Limit > 0 {
		params.Add("limit", strconv.Itoa(req.Limit))
	}
	if req.Offset > 0 {
		params.Add("offset", strconv.Itoa(req.Offset))
	}
	if req.SortBy != "" {
		params.Add("sort_by", req.SortBy)
	}
	if req.Order != "" {
		params.Add("order", req.Order)
	}

	path := "/api/v1/memories"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var memoryList MemoryListData
	if err := json.Unmarshal(dataBytes, &memoryList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memory list: %w", err)
	}

	return &memoryList, nil
}

func (c *Client) SearchMemories(req SearchMemoryRequest) (*SearchData, error) {
	resp, err := c.doRequest("POST", "/api/v1/memories/search", req)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var searchData SearchData
	if err := json.Unmarshal(dataBytes, &searchData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return &searchData, nil
}

func (c *Client) UpdateMemory(memoryID MemoryID, req UpdateMemoryRequest) (*Memory, error) {
	path := fmt.Sprintf("/api/v1/memories/%s", memoryID.String())
	resp, err := c.doRequest("PUT", path, req)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var memory Memory
	if err := json.Unmarshal(dataBytes, &memory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memory: %w", err)
	}

	return &memory, nil
}

func (c *Client) DeleteMemory(memoryID MemoryID, userID, agentID string) error {
	params := url.Values{}
	if userID != "" {
		params.Add("user_id", userID)
	}
	if agentID != "" {
		params.Add("agent_id", agentID)
	}

	path := fmt.Sprintf("/api/v1/memories/%s", memoryID.String())
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	_, err := c.doRequest("DELETE", path, nil)
	return err
}
