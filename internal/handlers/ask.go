package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sakuexe/go-raino/internal/openai"
)

type optionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func getGptAnswer(message string) string {
	var content string = message

	chat, err := openai.SendChat(content)

	if err != nil {
		fmt.Println(err)
		return "An error happened while trying to come up with a response..."
	}

	return chat.Choices[0].Message.Content
}

func askHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, options optionMap) {
	status := fmt.Sprintf("Responding to %s", interaction.Interaction.Member.User.Username)
	session.UpdateCustomStatus(status)
	defer session.UpdateStatusComplex(defaultStatus)

	// send a message about the question being in process, so a followup will come soon
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		// Type 5 is a deferred response with source
		Type: discordgo.InteractionResponseType(5),
	})
	if err != nil {
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Couldn't defer the response... Sorry about that",
		})
		return
	}

	// respond with a followup message including the generated answer
	content := getGptAnswer(options["message"].StringValue())
	session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	return
}

