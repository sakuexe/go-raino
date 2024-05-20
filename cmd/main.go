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
	discord, err := discordgo.New("Bot " + GetDotenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

  // Open the connection
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	// close the discord session automatically once the program ends
	defer discord.Close()

  // Show that the app was started successfully
  discord.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
    fmt.Printf("Raino is running as: %v#%v \n", session.State.User.Username, session.State.User.Discriminator)
  })

	discord.AddHandler(createMessage)
	discord.Identify.Intents = discordgo.IntentsGuildMessages

  registeredCommands := []*discordgo.ApplicationCommand{}
  for _, command := range Commands {
    fmt.Printf("Registering command: %v \n", command.Name)
    cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
    if err != nil {
      fmt.Printf("Error creating command: %v", command.Name)
      panic(err)
    }
    registeredCommands = append(registeredCommands, cmd)
  }

  // initialize the command handlers
  AddCommandHandlers(discord)

	// closing
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop // block until we receive a signal

	fmt.Println("Exiting.")

  if len(registeredCommands) == 0 {
    return
  }

  // remove all the registered commands
  for _, command := range registeredCommands {
    fmt.Printf("Removing command: %v \n", command.Name)
    err := discord.ApplicationCommandDelete(discord.State.User.ID, "", command.ID)
    if err != nil {
      fmt.Printf("Couldn't delete command: %v", command.Name)
      panic(err)
    }
  }
}
