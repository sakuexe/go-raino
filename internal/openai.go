package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func generateContent(content string) []byte {
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

func SendChat(apiKey string, chatContent string) (ChatResponse, error) {
  // openai API endpoint
  var url string = "https://api.openai.com/v1/chat/completions"

  jsonBody := generateContent("Explain why rocks are cool.")

  if len(jsonBody) == 0 {
    return ChatResponse{}, fmt.Errorf("Error generating content")
  }

  // generate a new post request
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

  if err != nil {
    return ChatResponse{}, err
  }

  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

  response, err := http.DefaultClient.Do(req)

  if err != nil {
    return ChatResponse{}, err
  }

  chatResponse, err := parseResponse(response)
  if err != nil {
    return ChatResponse{}, err
  }

  return chatResponse, nil
}
