package main

import (
	"github.com/bwmarrin/discordgo"
)

// todo: try the strategy pattern
type Command interface {
	HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
  // All the slash commands that the bot will have
	commands = []*discordgo.ApplicationCommand{
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

func addCommandHandlers(session *discordgo.Session) {
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
    // get the command name
		commandData := interaction.ApplicationCommandData()
		switch commandData.Name {

		case "hello-world":
			helloWorldHandler(session, interaction)
		case "responses":
			responsesHandler(session, interaction)
		}
	})
}
