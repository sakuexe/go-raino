package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	openai "github.com/sakuexe/go-raino/internal"
)

func gpt() {
	var apiKey string = GetDotenv("OPENAI_API_KEY")
	var content string = "Raino, why are rocks so cool?"

	chat, err := openai.SendChat(apiKey, content)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(chat.Choices[0].Message.Content)
}

func createMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if message.Author.ID == session.State.User.ID {
		return
	}
	// if the message is ping reply with pong!
	if message.Content == "ping" {
		session.ChannelMessageSend(message.ChannelID, "pong")
	}
	// if the message is hello reply with hello!
	if message.Content == "hello" {
		session.ChannelMessageSend(message.ChannelID, "Hello!")
	}
}

func main() {
  // invite raino to the server: https://discord.com/oauth2/authorize?client_id=1241964425317978193&permissions=40667002567744&scope=bot
	discord, err := discordgo.New("Bot " + GetDotenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

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

  // register the commands
  _, err = discord.ApplicationCommandBulkOverwrite(discord.State.User.ID, "", commands)
  if err != nil {
    fmt.Println("Error registering commands: ", err)
    return
  }

	// initialize the command handlers
	addCommandHandlers(discord)

  // initialize the message handler
	discord.AddHandler(createMessage)
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// closing
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop // block until we receive a signal

	fmt.Println("Exiting.")
}
