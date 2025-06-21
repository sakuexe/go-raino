package openai

import "github.com/sakuexe/go-raino/internal/env"

var (
	Model  = env.GetDotenv("GPT_MODEL")
	ApiUrl = "https://api.openai.com/v1/chat/completions"
	ApiKey = env.GetDotenv("OPENAI_API_KEY")
)

func init() {
	if Model == "" {
		Model = "gpt-3.5-turbo"
	}
}
