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
	// All the slash commands that the bot will have
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "hello-world",
			// all commands must have a description
			Description: "A basic way to check that everything is working",
		},

		{
			Name:        "ask",
			Description: "Ask Raino a question or tell him something",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "message",
					Description: "The message you want to send to Raino",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},

		{
			Name:        "convert",
			Description: "Converts a given image to a desired format",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "image-format",
					Description: "The format you want to convert the image to",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "png",
							Value: "png",
						},
						{
							Name:  "jpeg",
							Value: "jpeg",
						},
            {
              Name:  "webp",
              Value: "webp",
            },
					},
					Required: true,
				},
			},
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
		case "ask":
			// create a map of the options for easy access
			optionMap := make(optionMap)
			for _, option := range commandData.Options {
				optionMap[option.Name] = option
			}
			askHandler(session, interaction, optionMap)
    case "convert":
      convertHandler(session, interaction)
		}
	})
}

func removeUnusedCommands(session *discordgo.Session) {
	// https://github.com/bwmarrin/discordgo/issues/1518#issuecomment-2076083061
	// Get all the existing commands in the guild
	existingCommands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		fmt.Println("Error getting existing commands: ", err)
		return
	}

	// create a map of the command names
	commandNames := make(map[string]bool, len(commands))
	for _, command := range commands {
		commandNames[command.Name] = true
	}

	// iterate over the existing commands and remove the ones that are not in the list
	for _, command := range existingCommands {
		if _, found := commandNames[command.Name]; found {
			fmt.Println("Command found in list, keeping it", command.Name)
			continue
		}
		fmt.Println("Removing command: ", command.Name)
		// if the command can be found in the existing commands, continue
		err = session.ApplicationCommandDelete(session.State.User.ID, "", command.ID)
		if err != nil {
			fmt.Println("Error removing command: ", err)
			return
		}
	}
}
