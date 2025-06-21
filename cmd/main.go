package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/sakuexe/go-raino/internal/env"
	"github.com/sakuexe/go-raino/internal/handlers"
)

func createMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// make the message handler asynchronous
	go func() {
		// Ignore all messages created by the bot itself
		if message.Author.ID == session.State.User.ID {
			return
		}

		// if the message is hello reply with hello!
		askPattern, _ := regexp.Compile("([Hh]ey\\s?[Rr]aino)")
		if askPattern.MatchString(message.Content) {
			response := askMessageHandler(session, message)
			session.ChannelMessageSend(message.ChannelID, response)
		}

		rockPattern, _ := regexp.Compile("[Rr]ock")
		if rockPattern.MatchString(message.Content) {
			session.MessageReactionAdd(message.ChannelID, message.ID, "🪨")
		}
	}()
}

func main() {

	// invite raino to the server: https://discord.com/oauth2/authorize?client_id=1241964425317978193&permissions=40667002567744&scope=bot
	discord, err := discordgo.New("Bot " + env.GetDotenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Show that the app was started successfully
	discord.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		fmt.Printf("Raino is running as: %v#%v \n", session.State.User.Username, session.State.User.Discriminator)
	})

	// Open the connection
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	// close the discord session automatically once the program ends
	defer discord.Close()

	// in development, register the commands to a single guild in the
	// .env GUILD_ID variable. Otherwise register the commands globally.
	// This is so that the bot methods will be available instantly,
	// instead of waiting for the global commands to be registered
	// (about an hour).
	var guildID string = env.GetDotenv("GUILD_ID")

	_, err = discord.ApplicationCommandBulkOverwrite(discord.State.User.ID, guildID, handlers.Commands)
	if err != nil {
		fmt.Println("Error registering commands: ", err)
		return
	}

	// remove all the commands that do not exist anymore
	handlers.RemoveUnusedCommands(discord)

	// initialize the command handlers
	handlers.AddCommandHandlers(discord)

	// initialize the message handler
	discord.AddHandler(createMessage)

	// closing
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop // block until we receive a signal

	fmt.Println("Exiting.")

}
