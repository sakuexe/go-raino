package openai

import (
	"encoding/json"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/sakuexe/go-raino/internal/raino"
)

// gets a ChatRequest and converts it to a json byte array
func createRequestBody(request *ChatRequest) []byte {
	jsonBody, err := json.Marshal(request)

	if err != nil {
		slog.Error("Error converting request to json", "error", err.Error())
		return []byte{}
	}

	return jsonBody
}

func AskQuestion(message string) (*ChatResponse, error) {
	prompt := []map[string]string{
		{
			"role":    "developer",
			"content": raino.SystemPrompt,
		},
		{
			"role":    "user",
			"content": message,
		},
	}

	jsonBody := createRequestBody(&ChatRequest{Model: Model, Messages: prompt})
	return makeRequest(jsonBody)
}

func AnswerChat(chat []*discordgo.Message, sessionID string) (*ChatResponse, error) {
	messages := []map[string]string{
		{
			"role":    "developer",
			"content": raino.SystemPrompt,
		},
	}

	for _, message := range chat {
		newMessage := map[string]string{
			"role":    "assistant",
			"content": message.Content,
		}

		// if the message is not from raino itself, add user information
		if message.Author.ID != sessionID {
			newMessage["role"] = "user"
			newMessage["username"] = message.Author.Username
			newMessage["user_id"] = message.Author.ID
		}

		messages = append(messages, newMessage)
	}

	jsonBody := createRequestBody(&ChatRequest{Model: Model, Messages: messages})
	return makeRequest(jsonBody)
}
