package main

import (
	"fmt"

	openai "github.com/sakuexe/go-raino/internal"
)

func main() {
  fmt.Println("Hello, World!")
  var apiKey string = GetDotenv("OPENAI_API_KEY")
  var content string = "Raino, why are rocks so cool?"

  chat, err := openai.SendChat(apiKey, content)

  if err != nil {
    fmt.Println(err)
    return
  }

  fmt.Println(chat.Choices[0].Message.Content)
}
