package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// todo: try the strategy pattern
type Command interface {
	HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name: "hello-world",
			// all commands must have a description
			Description: "A basic way to check that everything is working",
		},
		{
			Name:        "responses",
			Description: "A way to check the responses",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "resp-type",
					Description: "The type of response you want to see",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Channel message with source",
							Value: 4,
						},
						{
							Name:  "Deferred response with source",
							Value: 5,
						},
					},
					Required: true,
				},
			},
		},
	}
)

func AddCommandHandlers(session *discordgo.Session) {
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		GetHandler(session, interaction)
	})
}

func GetHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// get the command name
	commandData := interaction.ApplicationCommandData()
	switch commandData.Name {

	case "hello-world":
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

	case "responses":
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
}
