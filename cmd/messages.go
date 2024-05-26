package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sakuexe/go-raino/internal/openai"
)

func askMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) string {
  messageHistorySize := 6
	// get the last 10 messages
	messages, err := session.ChannelMessages(message.ChannelID, messageHistorySize, "", "", "")
	if err != nil {
		fmt.Println("Error getting messages: ", err)
		session.ChannelMessageSend(message.ChannelID, "I'm sorry, I couldn't get the context of your message.")
	}

	// make a string of the last 10 messages with usernames
	context := fmt.Sprintf("here are the last %d messages of the discord chat. " +
		"Use them as the context for the latest message that you are replying to. " +
		"\n\n", messageHistorySize)

	// iterate through the messages from oldest to newest
	for index := len(messages) - 1; index >= 0; index-- {
		msg := messages[index]
		context += fmt.Sprintf("%v: %v\n", msg.Author.Username, msg.Content)
	}

	response, err := openai.SendChat(context)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, err.Error())
	}

  // if the message starts with Raino: then remove it
  if strings.HasPrefix(response.Choices[0].Message.Content, "Raino:") {
    return strings.TrimPrefix(response.Choices[0].Message.Content, "Raino:")
  }

  return response.Choices[0].Message.Content
}
