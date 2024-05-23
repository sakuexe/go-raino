package handlers

import (
	"bytes"
	"fmt"

	"github.com/bwmarrin/discordgo"
	imageconversion "github.com/sakuexe/go-raino/internal/image-conversion"
)

func convertHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// customize some variables
	var format string = interaction.ApplicationCommandData().Options[0].StringValue()
	var messageHistoryLimit int = 100

	// add a status for the process
	status := fmt.Sprintf("Converting an image to %s", format)
	session.UpdateCustomStatus(status)
	defer session.UpdateStatusComplex(defaultStatus)

	// send a message about the process being in progress
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(5),
	})
	if err != nil {
		fmt.Println("No status could be created/sent for the `convertHandler` command")
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Couldn't defer the response... Sorry about that",
		})
	}

	// get the user that sent the command
	user, err := session.User(interaction.Interaction.Member.User.ID)
	if err != nil {
		fmt.Println("Couldn't get the user that sent the command")
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Couldn't get the user that sent the command",
		})
		return
	}

	// get user's messages in the channel
	messages, _ := session.ChannelMessages(interaction.Interaction.ChannelID, messageHistoryLimit, "", "", "")
	// look for the message that is from the user and has the attachment
	var message *discordgo.Message = nil
	// loop in reverse order (newest to oldest)
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Author.ID == user.ID && len(msg.Attachments) > 0 {
			message = msg
		}
	}
	if message == nil {
		fmt.Printf("Couldn't a message with attachments from the user %s. Looked at the last 100 messages.\n", user.Username)
		rainoMessage := fmt.Sprintf("I looked at the last 100 messages and didn't find yours with attachments, %s", user.Username)
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: rainoMessage,
		})
		return
	}

	// get the image from the message
	attachmentUrl := message.Attachments[0].URL

	imageResponse, err := imageconversion.ConvertImage(format, attachmentUrl)
	if err != nil {
		fmt.Println("Error converting the image: ", err)
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: err.Error(),
		})
		return
	}

	// read the image file
	reader := bytes.NewReader(imageResponse.Buffer.Bytes())
	// fmt.Println("Filename: ", imageResponse.Filename)
	// fmt.Println("Content Type: ", imageResponse.ContentType)
	// fmt.Println("Filepath: ", imageResponse.Filepath)

	// send a follow up message
	var content string = fmt.Sprintf("Here you go, <@%s>:", user.ID)
	session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
		Content: content,
		Files: []*discordgo.File{
			{
				Name:        imageResponse.Filename,
				ContentType: imageResponse.ContentType,
				Reader:      reader,
			},
		},
	})
}
