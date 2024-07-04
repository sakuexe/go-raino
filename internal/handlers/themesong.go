package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	audioBuffer = make([][]byte, 0)
)

func themeSongHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// send a message about the process being in progress
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(5),
	})
	if err != nil {
		err := session.InteractionRespond(
			interaction.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "I'm sorry, I couldn't defer my response to `/themesong`...",
				},
			})
		if err != nil {
			fmt.Println("Error initializing response to themeSongHandler: \n", err)
		}
	}

	// get the guild id for the channel the command was sent in
	guild, err := session.State.Guild(interaction.GuildID)
	if err != nil {
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Couldn't get the guild from the message... Sorry about that",
		})
		return
	}

	// get the voice state of the user who sent the command
	user := interaction.Member.User
	voiceState, err := session.State.VoiceState(guild.ID, user.ID)
	if err != nil {
		fmt.Println(err)
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("I didn't find you in a voice channel <@%s>... Are you even in one?", user.ID),
		})
		return
	}

	// load the stupid fucking dca file that is impossible to convert to
  // file := "./media/themesong.mp3"

	// play the sound in the voice channel the user is in
	vc, err := session.ChannelVoiceJoin(guild.ID, voiceState.ChannelID, false, true)
	if err != nil {
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Couldn't join the voice channel... how rude.",
		})
		return
	}

	session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
		Content: "Enjoy!",
	})

	// wait for a specified time before playing the sound
	time.Sleep(500 * time.Millisecond)
	vc.Speaking(true)

	time.Sleep(500 * time.Millisecond)
	vc.Speaking(false)
	vc.Disconnect()
}
