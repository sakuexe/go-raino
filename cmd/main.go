package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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

func getOpenAiToken() string {
  err := godotenv.Load()

  if err != nil {
    fmt.Println("Error loading .env file at project root")
    return ""
  }

  token := os.Getenv("OPENAI_API_KEY")
  if token == "" {
    fmt.Println("Error: OPENAI_API_KEY not found in .env file")
    return ""
  }

  return token
}

func generateRequest(content string) []byte {
  body := ChatRequest{
    Model: "gpt-3.5-turbo",
    Messages: []map[string]string{
      {
        "role": "system",
        "content": "You are 'Raino', a helpful rhino character that speaks simply and dimly.",
      },
      {
        "role": "user",
        "content": content,
      },
    },
  }

  // convert into a json string
  jsonBody, err := json.Marshal(body)
  if err != nil {
    fmt.Println("Error converting body to json")
    fmt.Println(err)
    return []byte{}
  }

  return jsonBody
}

func main() {
  fmt.Println("Hello, World!")
  var url string = "https://api.openai.com/v1/chat/completions"
  var token string = getOpenAiToken()
  if token == "" {
    return
  }

  jsonBody := generateRequest("Explain why rocks are cool.")
  if len(jsonBody) == 0 {
    return
  }

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

  if err != nil {
    fmt.Printf("Error with get request to %s \n", url)
    fmt.Println(err)
  }

  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

  response, err := http.DefaultClient.Do(req)

  if err != nil {
    fmt.Printf("Error with response from %s \n", url)
    fmt.Println(err)
  }

  // read the response body
  body, err := io.ReadAll(response.Body)

  if err != nil {
    fmt.Printf("Error reading the response body from %s \n", url)
    fmt.Println(err)
  }

  responseMessage := ChatResponse{}
  err = json.Unmarshal(body, &responseMessage)

  if err != nil {
    fmt.Println("Error converting json response body")
    fmt.Println(err)
  }

  fmt.Println(responseMessage.Choices[0].Message.Content)
}
