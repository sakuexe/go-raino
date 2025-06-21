package main

import (
	"log/slog"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/sakuexe/go-raino/internal/openai"
	"github.com/sakuexe/go-raino/internal/raino"
)

func askMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) string {
	// get the last N messages
	messages, err := session.ChannelMessages(message.ChannelID, raino.MaxHistorySize, "", "", "")
	if err != nil {
		slog.Error("Error while fetching chat history", "error", err.Error(), "history_size", raino.MaxHistorySize)
		session.ChannelMessageSend(message.ChannelID, "I'm sorry, I couldn't get the context of your message.")
	}

	// make sure that the messages are sorted from latest to newest
	// we want it this way for the api calls
	slices.Reverse(messages)

	response, err := openai.AnswerChat(messages, session.State.SessionID)
	if err != nil {
		slog.Error("error getting a response from openai", "error", err.Error())
		session.ChannelMessageSend(message.ChannelID, err.Error())
	}

	if len(response.Choices) == 0 {
		slog.Error("response.Choices is empty", "response", response)
		return "Could not find message content in response..."
	}

	return response.Choices[0].Message.Content
}
