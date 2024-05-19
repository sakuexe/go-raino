package openai

type ChatRequest struct {
  Model string `json:"model"`
  Messages []map[string]string `json:"messages"`
}

type Message struct {
  Role string `json:"role"`
  Content string `json:"content"`
}

type Choice struct {
  Index int `json:"index"`
  Message Message `json:"message"`
  LogProps map[string]string `json:"log_props"`
  FinishReason string `json:"finish_reason"`
}

type Usage struct {
  PromptTokens int `json:"prompt_tokens"`
  CompletionTokens int `json:"completion_tokens"`
  TotalTokens int `json:"total_tokens"`
}

type ChatResponse struct {
  Id string `json:"id"`
  Object string `json:"object"`
  Created int `json:"created"`
  Model string `json:"model"`
  Choices []Choice `json:"choices"`
  Usage Usage `json:"usage"`
}

