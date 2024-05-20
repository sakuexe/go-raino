package main

import (
	"fmt"

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
							Value: 1,
						},
						{
							Name:  "Deferred response with source",
							Value: 2,
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
  fmt.Printf("Handling interaction %s \n", commandData.Name)
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
    return
	}
}
