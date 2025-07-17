package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClientGenerator is a generator for REST clients
type ClientGenerator struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClientGenerator creates a new client generator
func NewClientGenerator(baseURL string) *ClientGenerator {
	return &ClientGenerator{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// Request represents a REST request
type Request struct {
	Method string
	Path   string
	Body   interface{}
	Query  map[string]string
	Header map[string]string
}

// Response represents a REST response
type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// Do executes a REST request
func (c *ClientGenerator) Do(req *Request) (*Response, error) {
	// Prepare URL
	url := fmt.Sprintf("%s%s", c.BaseURL, req.Path)

	// Prepare body
	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest(req.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range req.Header {
		httpReq.Header.Set(k, v)
	}

	// Set query parameters
	q := httpReq.URL.Query()
	for k, v := range req.Query {
		q.Add(k, v)
	}
	httpReq.URL.RawQuery = q.Encode()

	// Execute request
	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create response
	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Body:       respBody,
		Header:     httpResp.Header,
	}

	return resp, nil
}
