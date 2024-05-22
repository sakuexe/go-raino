package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	imageconversion "github.com/sakuexe/go-raino/internal/image-conversion"
)

func helloWorldHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
  // add a status
  session.UpdateCustomStatus("Helloing they world")
  defer session.UpdateCustomStatus("")

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
  status := fmt.Sprintf("Responding to %s", interaction.Interaction.Member.User.Username)
  session.UpdateCustomStatus(status)
  defer session.UpdateCustomStatus("")

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

func convertHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var format string = interaction.ApplicationCommandData().Options[0].StringValue()
	var messageHistoryLimit int = 100

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
