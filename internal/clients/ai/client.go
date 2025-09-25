package ai

import "context"

type Message struct {
    Role    string // "system" | "user" | "assistant"
    Content string
}

type ChatOptions struct {
    Model string
}

type ChatModel interface {
    Name() string
    Chat(ctx context.Context, msgs []Message, opts ChatOptions) (string, error)
}
