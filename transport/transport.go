package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Transport struct {
	baseURL     string
	apiKey      string
	accessToken string
	timeout     time.Duration
	debug       bool
	userAgent   string
	client      *http.Client
	mu          sync.RWMutex
}

type RequestOptions struct {
	JSON          interface{}
	Params        map[string]string
	Headers       map[string]string
	SkipAuth      bool
	TokenOverride string
}

type RequestInterceptor func(method, path, url string, headers map[string]string, body []byte)
type ResponseInterceptor func(method, path string, status int, ok bool)

func NewTransport(baseURL, apiKey, accessToken string, timeout time.Duration, debug bool) *Transport {
	return NewTransportWithClient(baseURL, apiKey, accessToken, timeout, debug, "", nil)
}

// NewTransportWithClient creates a transport with an optional custom HTTP client.
// The caller retains ownership of the supplied client and its transport.
func NewTransportWithClient(baseURL, apiKey, accessToken string, timeout time.Duration, debug bool, userAgent string, client *http.Client) *Transport {
	base := strings.TrimRight(baseURL, "/")
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "https://" + base
	}
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}
	if userAgent == "" {
		userAgent = "blocklog-go/1.0.0"
	}
	return &Transport{
		baseURL:     base,
		apiKey:      apiKey,
		accessToken: accessToken,
		timeout:     timeout,
		debug:       debug,
		userAgent:   userAgent,
		client:      client,
	}
}

func (t *Transport) SetAccessToken(token string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.accessToken = token
}

func (t *Transport) Request(ctx context.Context, method, path string, opts *RequestOptions) ([]byte, error) {
	t.mu.RLock()
	baseURL, apiKey, accessToken, userAgent, client := t.baseURL, t.apiKey, t.accessToken, t.userAgent, t.client
	t.mu.RUnlock()
	cleanPath := path
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}
	if strings.HasSuffix(baseURL, "/api/v1") && strings.HasPrefix(cleanPath, "/api/v1/") {
		cleanPath = cleanPath[7:]
	} else if strings.HasSuffix(baseURL, "/api/v1") && cleanPath == "/api/v1" {
		cleanPath = ""
	} else if !strings.HasSuffix(baseURL, "/api/v1") && !strings.HasPrefix(cleanPath, "/api/v1/") && cleanPath != "/api/v1" {
		cleanPath = "/api/v1" + cleanPath
	}

	fullURL := baseURL + cleanPath

	var body []byte
	var err error
	if opts != nil && opts.JSON != nil {
		body, err = json.Marshal(opts.JSON)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	if opts != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	token := accessToken
	if opts != nil && opts.TokenOverride != "" {
		token = opts.TokenOverride
	}
	if opts == nil || !opts.SkipAuth {
		if token == "" && apiKey == "" {
			return nil, ErrMissingAuth
		}
		authToken := token
		if authToken == "" {
			authToken = apiKey
		}
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	if opts != nil && opts.Params != nil {
		q := req.URL.Query()
		for k, v := range opts.Params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if t.debug {
		fmt.Printf("[blocklog] %s %s\n", method, req.URL.String())
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if t.debug {
		fmt.Printf("[blocklog] Response status: %d\n", resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		msg := string(respBody)
		if msg == "" {
			msg = fmt.Sprintf("HTTP Error %d", resp.StatusCode)
		}
		return nil, mapHTTPError(resp.StatusCode, msg, respBody)
	}

	if resp.StatusCode == 204 {
		return nil, nil
	}

	return respBody, nil
}

func (t *Transport) Get(ctx context.Context, path string, opts *RequestOptions) ([]byte, error) {
	return t.Request(ctx, "GET", path, opts)
}

func (t *Transport) Post(ctx context.Context, path string, opts *RequestOptions) ([]byte, error) {
	return t.Request(ctx, "POST", path, opts)
}

func (t *Transport) Put(ctx context.Context, path string, opts *RequestOptions) ([]byte, error) {
	return t.Request(ctx, "PUT", path, opts)
}

func (t *Transport) Patch(ctx context.Context, path string, opts *RequestOptions) ([]byte, error) {
	return t.Request(ctx, "PATCH", path, opts)
}

func (t *Transport) Delete(ctx context.Context, path string, opts *RequestOptions) ([]byte, error) {
	return t.Request(ctx, "DELETE", path, opts)
}

func BuildURL(baseURL, path string, params map[string]string) string {
	cleanPath := path
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}
	if strings.HasSuffix(baseURL, "/api/v1") && strings.HasPrefix(cleanPath, "/api/v1/") {
		cleanPath = cleanPath[7:]
	} else if strings.HasSuffix(baseURL, "/api/v1") && cleanPath == "/api/v1" {
		cleanPath = ""
	} else if !strings.HasSuffix(baseURL, "/api/v1") && !strings.HasPrefix(cleanPath, "/api/v1/") && cleanPath != "/api/v1" {
		cleanPath = "/api/v1" + cleanPath
	}
	fullURL := baseURL + cleanPath
	if len(params) > 0 {
		u, _ := url.Parse(fullURL)
		q := u.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}
	return fullURL
}
