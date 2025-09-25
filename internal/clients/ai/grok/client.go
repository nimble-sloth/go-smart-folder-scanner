package grok

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "os"
    "time"

    "github.com/nimble-sloth/go-smart-folder-scanner/internal/clients/ai"
)

type Client struct {
    http  *http.Client
    key   string
    url   string
    model string
}

func New(model, endpoint string, timeout time.Duration) *Client {
    url := endpoint
    if url == "" { url = "https://api.x.ai/v1/chat/completions" }
    return &Client{
        http:  ai.DefaultHTTPClient(timeout),
        key:   os.Getenv("GROK_API_KEY"),
        url:   url,
        model: model,
    }
}

func (c *Client) Name() string { return "grok" }

func (c *Client) Chat(ctx context.Context, msgs []ai.Message, opts ai.ChatOptions) (string, error) {
    if c.key == "" { return "", errors.New("GROK_API_KEY not set") }
    payload := struct {
        Model    string       `json:"model"`
        Messages []ai.Message `json:"messages"`
    }{ Model: choose(opts.Model, c.model), Messages: msgs }

    body, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
    req.Header.Set("Authorization", "Bearer "+c.key)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.http.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode/100 != 2 { return "", errors.New("grok: non-2xx response") }

    var out struct {
        Choices []struct{ Message ai.Message `json:"message"` } `json:"choices"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    if len(out.Choices) == 0 { return "", errors.New("grok: empty choices") }
    return out.Choices[0].Message.Content, nil
}

func choose(a, b string) string { if a != "" { return a }; return b }
