package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func helloWorldHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(
		interaction.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello, World!",
			},
		})
	if err != nil {
		fmt.Println("Error responding to hello-world command: ", err)
	}
}

type optionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func askHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, options optionMap) {
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
	content := gpt(options["message"].StringValue())
	session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	return
}

func responsesHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	content := ""
	switch interaction.ApplicationCommandData().Options[0].IntValue() {
	case int64(discordgo.InteractionResponseChannelMessageWithSource):
		content =
			"You just responded to an interaction, sent a message and showed the original one. " +
				"Congratulations!"
		content +=
			"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
	default:
		err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseType(interaction.ApplicationCommandData().Options[0].IntValue()),
		})
		if err != nil {
			session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong",
			})
		}
		return
	}

	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(interaction.ApplicationCommandData().Options[0].IntValue()),
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
		return
	}
	time.AfterFunc(time.Second*5, func() {
		content := content + "\n\nWell, now you know how to create and edit responses. " +
			"But you still don't know how to delete them... so... wait 10 seconds and this " +
			"message will be deleted."
		_, err = session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		if err != nil {
			session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong",
			})
			return
		}
		time.Sleep(time.Second * 10)
		session.InteractionResponseDelete(interaction.Interaction)
	})
}
