package lara

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	accessKeyID     string
	accessKeySecret string
	baseURL         string
	httpClient      *http.Client
	sdkName         string
	sdkVersion      string
}

func newClient(accessKeyID, accessKeySecret, baseURL string) *Client {
	return &Client{
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		baseURL:         strings.TrimRight(baseURL, "/"),
		httpClient:      &http.Client{},
		sdkName:         "lara-go",
		sdkVersion:      "1.1.0",
	}
}

func (c *Client) request(method, path string, params map[string]string, body interface{}, files map[string]io.Reader, headers map[string]string) ([]byte, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	reqURL := c.baseURL + path
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		reqURL += "?" + values.Encode()
	}

	var bodyReader io.Reader
	var contentMD5, contentType string

	if len(files) > 0 {
		// Multipart form
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		for field, r := range files {
			fw, err := w.CreateFormFile(field, field)
			if err != nil {
				return nil, fmt.Errorf("failed to create form file: %w", err)
			}
			if _, err := io.Copy(fw, r); err != nil {
				return nil, fmt.Errorf("failed to copy file data: %w", err)
			}
		}
		if body != nil {
			jsonBytes, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body for multipart: %w", err)
			}
			_ = w.WriteField("json", string(jsonBytes))
		}
		w.Close()
		bodyBytes := b.Bytes()
		bodyReader = bytes.NewReader(bodyBytes)
		contentType = w.FormDataContentType()
		hash := md5.Sum(bodyBytes)
		contentMD5 = fmt.Sprintf("%x", hash)
	} else if body != nil {

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
		contentType = "application/json"
		hash := md5.Sum(bodyBytes)
		contentMD5 = fmt.Sprintf("%x", hash)
	}

	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-HTTP-Method-Override", method)
	req.Header.Set("Date", c.httpDate())
	req.Header.Set("X-Lara-SDK-Name", c.sdkName)
	req.Header.Set("X-Lara-SDK-Version", c.sdkVersion)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if contentMD5 != "" {
		req.Header.Set("Content-MD5", contentMD5)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if c.accessKeyID != "" && c.accessKeySecret != "" {
		signature := c.sign(method, path, contentMD5, contentType, req.Header.Get("Date"))
		req.Header.Set("Authorization", fmt.Sprintf("Lara %s:%s", c.accessKeyID, signature))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return nil, &LaraTimeoutError{Message: err.Error()}
		}
		return nil, &LaraConnectionError{Message: err.Error()}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	return respBody, nil
}

func (c *Client) httpDate() string {
	return time.Now().UTC().Format(http.TimeFormat)
}

// sign creates the HMAC signature for authentication
func (c *Client) sign(method, path, contentMD5, contentType, date string) string {
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		strings.ToUpper(method),
		path,
		contentMD5,
		contentType,
		date,
	)

	mac := hmac.New(sha256.New, []byte(c.accessKeySecret))
	mac.Write([]byte(stringToSign))
	signature := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(signature)
}

func parseAPIError(statusCode int, body []byte) *LaraError {
	var apiError struct {
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error.Type != "" {
		return &LaraError{
			Status:  statusCode,
			Type:    apiError.Error.Type,
			Message: apiError.Error.Message,
		}
	}
	return &LaraError{
		Status:  statusCode,
		Type:    "UnknownError",
		Message: "An unknown error occurred",
	}
}

func processStreamLine(line []byte, callback func([]byte) error) error {
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return nil
	}

	var resp struct {
		Status  int             `json:"status"`
		Data    json.RawMessage `json:"data"`
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil
	}

	if resp.Status < 200 || resp.Status >= 300 {
		errBody := json.RawMessage(line)
		if len(resp.Data) > 0 {
			errBody = resp.Data
		}
		return parseAPIError(resp.Status, errBody)
	}

	var content json.RawMessage
	if len(resp.Data) > 0 {
		var d struct {
			Content json.RawMessage `json:"content"`
		}
		if err := json.Unmarshal(resp.Data, &d); err == nil {
			content = d.Content
		}
	} else if len(resp.Content) > 0 {
		content = resp.Content
	}

	if len(content) > 0 {
		return callback(content)
	}
	return nil
}

func (c *Client) handleContent(respBytes []byte, result interface{}) error {
	var apiResponse struct {
		Content json.RawMessage `json:"content"`
	}

	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	return json.Unmarshal(apiResponse.Content, result)
}

func (c *Client) Get(path string, params map[string]string, headers map[string]string, result interface{}) error {
	respBytes, err := c.request("GET", path, params, nil, nil, headers)
	if err != nil {
		return err
	}
	return c.handleContent(respBytes, result)
}

func (c *Client) Post(path string, body interface{}, files map[string]io.Reader, headers map[string]string, result interface{}) error {
	respBytes, err := c.request("POST", path, nil, body, files, headers)
	if err != nil {
		return err
	}
	return c.handleContent(respBytes, result)
}

func (c *Client) Put(path string, body interface{}, files map[string]io.Reader, headers map[string]string, result interface{}) error {
	respBytes, err := c.request("PUT", path, nil, body, files, headers)
	if err != nil {
		return err
	}
	return c.handleContent(respBytes, result)
}

func (c *Client) Delete(path string, body interface{}, headers map[string]string, result interface{}) error {
	respBytes, err := c.request("DELETE", path, nil, body, nil, headers)
	if err != nil {
		return err
	}
	return c.handleContent(respBytes, result)
}

func (c *Client) GetRaw(path string, params map[string]string, headers map[string]string) ([]byte, error) {
	return c.request("GET", path, params, nil, nil, headers)
}

// PostAndGetStream makes a POST request and returns a channel for streaming results
func (c *Client) PostAndGetStream(path string, body interface{}, headers map[string]string, callback func([]byte) error) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	reqURL := c.baseURL + path

	var bodyReader io.Reader
	var contentMD5, contentType string

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
		contentType = "application/json"
		hash := md5.Sum(bodyBytes)
		contentMD5 = fmt.Sprintf("%x", hash)
	}

	req, err := http.NewRequest("POST", reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-HTTP-Method-Override", "POST")
	req.Header.Set("Date", c.httpDate())
	req.Header.Set("X-Lara-SDK-Name", c.sdkName)
	req.Header.Set("X-Lara-SDK-Version", c.sdkVersion)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if contentMD5 != "" {
		req.Header.Set("Content-MD5", contentMD5)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if c.accessKeyID != "" && c.accessKeySecret != "" {
		signature := c.sign("POST", path, contentMD5, contentType, req.Header.Get("Date"))
		req.Header.Set("Authorization", fmt.Sprintf("Lara %s:%s", c.accessKeyID, signature))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return &LaraTimeoutError{Message: err.Error()}
		}
		return &LaraConnectionError{Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return parseAPIError(resp.StatusCode, respBody)
	}

	// Read stream line by line (NDJSON format)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if err := processStreamLine(scanner.Bytes(), callback); err != nil {
			return err
		}
	}
	return scanner.Err()
}
