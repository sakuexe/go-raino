package handlers

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	imageconversion "github.com/sakuexe/go-raino/internal/image-conversion"
)

var messageHistoryLimit int = 100

func findAttachments(session *discordgo.Session,
	interaction *discordgo.Interaction) ([]*discordgo.MessageAttachment, error) {

	// get the user that sent the command
	user := interaction.Member.User

	// get all the messages in the channel (up to the set limit)
	messages, _ := session.ChannelMessages(interaction.ChannelID, messageHistoryLimit, "", "", "")

	// look for the newest message from the user that has an attachment
	var message *discordgo.Message = nil
	// loop in reverse order (newest to oldest)
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Author.ID == user.ID && len(msg.Attachments) > 0 {
			message = msg
		}
	}

	if message == nil {
		rainoMessage := fmt.Sprintf("I looked at the last %d messages and didn't find yours with attachments, %s",
			messageHistoryLimit, user.Username)
		return nil, fmt.Errorf(rainoMessage)
	}

	// get the image from the message
	return message.Attachments, nil
}

func convertAttachmentsParallel(attachments []*discordgo.MessageAttachment, format string) ([]*discordgo.File, []error) {
	files := []*discordgo.File{}
	waitGroup := sync.WaitGroup{}
  // create a channel to send errors to
	errChannel := make(chan error, len(attachments))

	for index, attachment := range attachments {
		waitGroup.Add(1)
		go func(index int, attachment *discordgo.MessageAttachment) {
			defer waitGroup.Done()
			imageResponse, err := imageconversion.ConvertImage(format, attachment.URL)

			// if there was an error, send it to the error channel and return
			if err != nil {
				errorMessage := fmt.Errorf("**Error with image %d**: %s", index+1, err.Error())
				errChannel <- errorMessage
				return
			}

			reader := bytes.NewReader(imageResponse.Buffer.Bytes())
			files = append(files, &discordgo.File{
				Name:        imageResponse.Filename,
				ContentType: imageResponse.ContentType,
				Reader:      reader,
			})
		}(index, attachment)
	}

  // wait for all the goroutines to finish
	waitGroup.Wait()
	close(errChannel)

  // get all the errors from the channel
  errors := []error{}
  for err := range errChannel {
    errors = append(errors, err)
  }
	return files, errors
}

func convertHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var format string = interaction.ApplicationCommandData().Options[0].StringValue()

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

	attachments, err := findAttachments(session, interaction.Interaction)
	if err != nil {
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: err.Error(),
		})
		return
	}

	files, errors := convertAttachmentsParallel(attachments, format)

	// send a follow up message
	user := interaction.Member.User
	var content string = fmt.Sprintf("Here you go, <@%s>:", user.ID)

	if len(errors) > 0 {
		for _, err := range errors {
			content += fmt.Sprintf("\n%s", err.Error())
		}
	}

	session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
		Content: content,
		Files:   files,
	})
}
