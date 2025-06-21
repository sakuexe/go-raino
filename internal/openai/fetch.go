package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func makeRequest(jsonBody []byte) (*ChatResponse, error) {
	if len(jsonBody) == 0 {
		slog.Error("json that was passed is empty")
		return nil, fmt.Errorf("I couldn't read your message...")
	}

	// generate a new post request
	req, err := http.NewRequest(http.MethodPost, ApiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		slog.Error("Error creating request to OpenAI API", "error", err.Error(), "url", ApiUrl, "body", jsonBody)
		return nil, fmt.Errorf("I couldn't find my thinking rock...")
	}

	// the required headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ApiKey))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Error sending request to OpenAI API", "error", err.Error(), "url")
		return nil, fmt.Errorf("I could not get a connection on my stone tablet... Sorry about that.")
	}

	chatResponse, err := parseResponse(response)
	if err != nil {
		fmt.Println("Error parsing response from OpenAI API:", err)
		return nil, fmt.Errorf("Something went wrong while formatting my response...")
	}

	return &chatResponse, nil
}

func parseResponse(response *http.Response) (ChatResponse, error) {
	// read the response body
	body, err := io.ReadAll(response.Body)

	if err != nil {
		slog.Error("error while parsing response", err.Error())
		return ChatResponse{}, err
	}

	responseMessage := ChatResponse{}
	err = json.Unmarshal(body, &responseMessage)

	if err != nil {
		slog.Error("Error converting json response body", "error", err.Error())
		return ChatResponse{}, err
	}

	return responseMessage, nil
}
