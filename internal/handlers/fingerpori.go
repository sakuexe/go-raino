package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

func fingerporiHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	randomComic := false
	options := interaction.ApplicationCommandData().Options
	if len(options) > 0 {
		randomComic = options[0].BoolValue()
	}

	// add a status for the process
	status := "Fetching latest fingerpori"
	if randomComic {
		status = "Fetching random fingerpori"
	}
	session.UpdateCustomStatus(status)
	defer session.UpdateStatusComplex(defaultStatus)

	// send a message about the process being in progress
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(5),
	})
	if err != nil {
		slog.Error("No initial interaction status could be created/sent", "error", err.Error(), "handler", fingerporiHandler)
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Couldn't defer the response... Sorry about that",
		})
	}

	fingerporiUrl := getFingerporiUrl(randomComic)
	if fingerporiUrl == "" {
		session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Could not fetch the fingerpori comic... sorry!",
		})
		return
	}

	// send a follow up message
	user := interaction.Member.User
	content := fmt.Sprintf("Here you go <@%s>, your requested [fingerpori](%s)", user.ID, fingerporiUrl)

	session.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
}

func getFingerporiUrl(randomize bool) string {
	url := "https://www.kaleva.fi/sarjakuvat/fingerpori"
	if randomize {
		url += fmt.Sprintf("/%d", getRandomID())
		slog.Info("generated url", "url", url)
	}

	res, err := http.Get(url)
	if err != nil {
		slog.Error("could not fetch website from url", "error", err.Error(), "url", url)
		return ""
	}
	defer res.Body.Close()


	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("could not read response body", "error", err.Error())
	}

	// funky! match the first cartoon strip image and get the img source
	// https://regex101.com/r/5mYD9a/1
	re := regexp.MustCompile(`class="cartoon-strip__image.*?>[.\s]*(<.*>)?[.\s]*<img src="([^"]+?)\?`)

	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 1 {
		slog.Error("not enough matches were found", "len", len(matches), "regex", re.String())
		return ""
	}

	slog.Info("found image link", "url", matches[len(matches) - 1])
	return matches[len(matches) - 1]
}

func getRandomID() int {
	firstID := 11610182
	// get days between the first fingerpori with similiar URL and current date
	daysSinceFirst := time.Since(time.Date(2025, time.May, 13, 0, 0, 0, 0, time.Local)).Hours() / 24
	slog.Debug("days since 13.5.2025", "value", daysSinceFirst)
	lastID := firstID + int(daysSinceFirst)

	return rand.Intn(lastID - firstID) + firstID
}
