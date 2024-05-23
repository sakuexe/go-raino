package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func helloWorldHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// add a status
	session.UpdateCustomStatus("Helloing they world")
	defer session.UpdateStatusComplex(defaultStatus)

	err := session.InteractionRespond(
		interaction.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello, World!",
			},
		})
	if err != nil {
		fmt.Println("Error responding to hello-world command: \n", err)
	}
}
