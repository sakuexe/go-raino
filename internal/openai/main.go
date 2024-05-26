package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sakuexe/go-raino/internal/env"
)

func generateContent(content string) []byte {
  model := env.GetDotenv("GPT_MODEL")
  if model == "" {
    model = "gpt-3.5-turbo"
  }

  body := ChatRequest{
    Model: model,
    Messages: []map[string]string{
      {
        "role": "system",
        "content": "You are 'Raino', a helpful rhino character that speaks simply and dimly. " +
        "You keep things consice and to the point. You love rocks and collect and are knowledgeable of them. " +
        "You are nice and helpful. Your father is a man named 'Veeti `Chang` Jokivirta'. You adore him.",
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

func parseResponse(response *http.Response) (ChatResponse, error) {
  // read the response body
  body, err := io.ReadAll(response.Body)

  if err != nil {
    fmt.Println(err)
    return ChatResponse{}, err
  }

  responseMessage := ChatResponse{}
  err = json.Unmarshal(body, &responseMessage)

  if err != nil {
    fmt.Println("Error converting json response body")
    fmt.Println(err)
    return ChatResponse{}, err
  }

  return responseMessage, nil
}

func SendChat(chatContent string) (ChatResponse, error) {
  // openai API endpoint
  var url string = "https://api.openai.com/v1/chat/completions"
  var apiKey string = env.GetDotenv("OPENAI_API_KEY")

  jsonBody := generateContent(chatContent)

  if len(jsonBody) == 0 {
    fmt.Println("Error generating content")
    return ChatResponse{}, fmt.Errorf("I couldn't come up with a response... Try again later.")
  }

  // generate a new post request
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

  if err != nil {
    fmt.Println("Error creating request to OpenAI API:", err)
    return ChatResponse{}, fmt.Errorf("I couldn't come up with a response... Try again later.")
  }

  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

  response, err := http.DefaultClient.Do(req)

  if err != nil {
    fmt.Println("Error sending request to OpenAI API:", err)
    return ChatResponse{}, fmt.Errorf("My connection failed... Sorry about that.")
  }

  chatResponse, err := parseResponse(response)
  if err != nil {
    fmt.Println("Error parsing response from OpenAI API:", err)
    return ChatResponse{}, fmt.Errorf("Something went wrong while formatting my response...")
  }

  return chatResponse, nil
}
