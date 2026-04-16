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
	accessKey    *AccessKey
	token        string
	refreshToken string
	baseURL      string
	httpClient   *http.Client
	sdkName      string
	sdkVersion   string
}

type authResponse struct {
	Token string `json:"token"`
}

func newClient(auth interface{}, baseURL string) *Client {
	client := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{},
		sdkName:    "lara-go",
		sdkVersion: "1.2.0",
	}

	// Set authentication method based on type
	switch a := auth.(type) {
	case *AccessKey:
		client.accessKey = a
	case *Credentials:
		// Backward compatibility with deprecated Credentials
		client.accessKey = a.AccessKey
	case *AuthToken:
		// Use pre-existing token directly
		client.token = a.Token
		client.refreshToken = a.RefreshToken
	}
	// If auth is nil or unknown type, accessKey remains nil

	return client
}

// isTokenExpired checks if the current JWT token is expired or about to expire.
func (c *Client) isTokenExpired() bool {
	if c.token == "" {
		return true
	}

	parts := strings.Split(c.token, ".")
	if len(parts) != 3 {
		return true
	}

	payload := parts[1]
	if m := len(payload) % 4; m != 0 {
		payload += strings.Repeat("=", 4-m)
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return true
	}

	var claims struct {
		Exp float64 `json:"exp"`
	}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return true
	}
	if claims.Exp == 0 {
		return true
	}

	return claims.Exp <= float64(time.Now().Unix())+5
}

// authenticateWithAccessKey authenticates using access key with challenge-response
func (c *Client) authenticateWithAccessKey() error {
	path := "/v2/auth"
	method := "POST"

	authData := map[string]string{
		"id": c.accessKey.ID,
	}

	bodyBytes, err := json.Marshal(authData)
	if err != nil {
		return fmt.Errorf("failed to marshal auth request: %w", err)
	}

	reqURL := c.baseURL + path
	req, err := http.NewRequest(method, reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	// Calculate MD5 hash of body
	hash := md5.Sum(bodyBytes)
	contentMD5 := fmt.Sprintf("%x", hash)
	contentType := "application/json"

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-MD5", contentMD5)
	req.Header.Set("X-Lara-Date", c.httpDate())
	req.Header.Set("X-Lara-SDK-Name", c.sdkName)
	req.Header.Set("X-Lara-SDK-Version", c.sdkVersion)

	// Sign request with HMAC for authentication
	signature := c.sign(method, path, contentMD5, contentType, req.Header.Get("X-Lara-Date"))
	req.Header.Set("Authorization", fmt.Sprintf("Lara:%s", signature))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute auth request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("authentication failed: %s", string(respBody))
	}

	var authResp authResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return fmt.Errorf("failed to parse auth response: %w", err)
	}

	c.token = authResp.Token
	c.refreshToken = resp.Header.Get("x-lara-refresh-token")

	return nil
}

// refreshTokens refreshes the JWT token using the refresh token.
func (c *Client) refreshTokens() error {
	path := "/v2/auth/refresh"
	method := "POST"

	reqURL := c.baseURL + path
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.refreshToken)
	req.Header.Set("X-Lara-SDK-Name", c.sdkName)
	req.Header.Set("X-Lara-SDK-Version", c.sdkVersion)
	req.Header.Set("X-Lara-Date", c.httpDate())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute refresh request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("refresh failed: %s", string(respBody))
	}

	var authResp authResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return fmt.Errorf("failed to parse refresh response: %w", err)
	}

	c.token = authResp.Token
	if newRefreshToken := resp.Header.Get("x-lara-refresh-token"); newRefreshToken != "" {
		c.refreshToken = newRefreshToken
	}

	return nil
}

// refreshOrReauthenticate tries to refresh the token first, falls back to full authentication.
func (c *Client) refreshOrReauthenticate() error {
	if c.refreshToken != "" {
		if err := c.refreshTokens(); err != nil {
			c.refreshToken = ""
			if c.accessKey == nil {
				return err
			}
		} else {
			return nil
		}
	}

	if c.accessKey != nil {
		return c.authenticateWithAccessKey()
	}

	return fmt.Errorf("no authentication method available for token renewal")
}

func (c *Client) request(method, path string, params map[string]string, body interface{}, files map[string]io.Reader, headers map[string]string) ([]byte, error) {
	return c.doRequest(method, path, params, body, files, headers, 0)
}

func (c *Client) doRequest(method, path string, params map[string]string, body interface{}, files map[string]io.Reader, headers map[string]string, retryCount int) ([]byte, error) {
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
	req.Header.Set("X-Lara-Date", c.httpDate())
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

	// Ensure we have a valid, non-expired token before making the request
	if c.isTokenExpired() {
		c.token = ""
		if err := c.refreshOrReauthenticate(); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Use JWT Bearer token for authorization
	req.Header.Set("Authorization", "Bearer "+c.token)

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

	// Handle 401 with automatic token refresh and retry (once)
	if resp.StatusCode == 401 && retryCount < 1 {
		c.token = ""
		if err := c.refreshOrReauthenticate(); err != nil {
			return nil, fmt.Errorf("token refresh failed: %w", err)
		}
		return c.doRequest(method, path, params, body, files, headers, retryCount+1)
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

	mac := hmac.New(sha256.New, []byte(c.accessKey.Secret))
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
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &apiError); err == nil {
		if apiError.Error.Type != "" {
			return &LaraError{
				Status:  statusCode,
				Type:    apiError.Error.Type,
				Message: apiError.Error.Message,
			}
		}
		if apiError.Type != "" {
			return &LaraError{
				Status:  statusCode,
				Type:    apiError.Type,
				Message: apiError.Message,
			}
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
		if err := json.Unmarshal(resp.Data, &d); err == nil && len(d.Content) > 0 {
			content = d.Content
		} else {
			content = resp.Data
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
	return json.Unmarshal(respBytes, result)
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

// PostAndGetStream makes a POST request and processes the response as an NDJSON stream.
func (c *Client) PostAndGetStream(path string, body interface{}, headers map[string]string, callback func([]byte) error) error {
	return c.doPostAndGetStream(path, body, headers, callback, 0)
}

func (c *Client) doPostAndGetStream(path string, body interface{}, headers map[string]string, callback func([]byte) error, retryCount int) error {
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
	req.Header.Set("X-Lara-Date", c.httpDate())
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

	// Ensure we have a valid, non-expired token before making the request
	if c.isTokenExpired() {
		c.token = ""
		if err := c.refreshOrReauthenticate(); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Use JWT Bearer token for authorization
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return &LaraTimeoutError{Message: err.Error()}
		}
		return &LaraConnectionError{Message: err.Error()}
	}
	defer resp.Body.Close()

	// Handle 401 with automatic token refresh and retry (once)
	if resp.StatusCode == 401 && retryCount < 1 {
		c.token = ""
		if err := c.refreshOrReauthenticate(); err != nil {
			return fmt.Errorf("token refresh failed: %w", err)
		}
		return c.doPostAndGetStream(path, body, headers, callback, retryCount+1)
	}

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
