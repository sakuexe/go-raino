# Go Raino! - The Helpful Discord Bot

Raino is a Discord bot that mainly provides some useful functionalities
that I was needing for my own server. These functionalities include:

- Image Conversion: Convert images to different formats.
- Chatbot: Talk to Raino and get responses. He's a nice guy.

This project is still very much a work in progress, so expect some bugs
and missing features. If you find any, feel free to open an issue or
a pull request.

From the media directory, you can find some images to use for the bot's
image, so you can get the authentic Raino experience!

## Installation

To get Raino up and running you need the following:

1. **Install Dependencies**

Raino uses the DiscordGo library for interfacing with Discord's API.
It also needs the necessary libraries for image conversion to work.

```bash
# Go packages
go mod download && go mod verify
# Webp library for image conversion
sudo apt install libwebp-dev -y
```

2. **Discord Bot Token** 

If you haven't hosted your own discord bot before, you need to set up
a discord application from the [Discord Developer Portal](https://discord.com/developers/applications).

Once you have it created, grab the token from the `Bot` section. 
Click the `Reset Token` for the code to show up. Then Save it somewhere safe.
It will be used in a later step.

3. **OAuth2 Permissions** 

Next step is to invite the bot to your server.
Go to the `OAuth2` section of your application and select the `bot` scope.
Then select the permissions you want the bot to have.

Once finished, copy the generated URL and paste it in your browser. Follow
the instructions to add the bot to a server of your choosing.

4. **Optional: OpenAI API Key**

If you want to use the chatbot functionality, you need to get an API key from
OpenAI. You can get one by signing up at their [website](https://platform.openai.com/).

Generate a token and add some balance to the account if needed. Save the token,
because it will get put into the bot's configuration file.

You can also choose the exact model you want to use for the chatbot functionality.
Just fill the `GPT_MODEL` field in the `.env` file with the model name. The default
model used, if no model is specified, is `gpt-3.5-turbo`.

## Run Go Raino!

To run the bot, you need to have Go installed on your machine. If you don't have it,
you can download it from the [official website](https://golang.org/). Or through
the package manager of your system. (e.g. `zypper in go go-doc`)

1. Clone the repository to your machine. And navigate to the directory.

```bash
git clone https://github.com/sakuexe/go-raino.git
cd go-raino
```

2. Add a `.env` file for configuring the bot. You can copy the `env.example` file
and fill in the necessary values.

If you are developing this, remember to add the `GUILD_ID` for faster iteration.

```bash
mv env.example .env
vi .env # remember to fill in the values
```

3. Build the bot and run it.

```bash
go build cmd/
./cmd
```

### Run with Docker

1. Build the Docker image.

```bash
docker build -t go-raino .
```

2. Run the Docker container.

```bash
docker compose up --build -d
```

Because the bot was authenticated with OAuth2, it shouldn't need anything else.

## Help

Here are some helpful guides and resources that I used to build the bot. Maybe
you can find them useful too!

- [DiscordGo Documentation](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [DiscordGo Examples](https://github.com/bwmarrin/discordgo/tree/master/examples) 
- [OpenAI API Documentation](https://beta.openai.com/docs/)
- [Writing Modern Discord Bots on Go](url) by: [Fedor Lapshin](https://medium.com/@lapfed255) 
- [Create a Bot for Discord in Go](https://dev.to/aurelievache/learning-go-by-examples-part-4-create-a-bot-for-discord-in-go-43cf) by [Aurélie Vache](https://dev.to/aurelievache)
