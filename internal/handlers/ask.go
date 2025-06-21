package handlers

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/sakuexe/go-raino/internal/openai"
)

type optionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

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
		slog.Error("error while sending message", "error", err.Error())
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Couldn't defer the response... Sorry about that",
		})
		return
	}

	// respond with a followup message including the generated answer
	content, err := openai.AskQuestion(options["message"].StringValue())
	if err != nil {
		slog.Error("something went wrong during /ask response generation", "error", err.Error())
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "something went wrong while reading my stone tablet...",
		})
		return
	}

	if len(content.Choices) == 0 {
		slog.Error("Response's Choices is empty (maybe insufficent tokens?)")
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Sorry friend! My stone tablet ran out of data...",
		})
		return
	}

	session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
		Content: content.Choices[0].Message.Content,
	})
	return
}
