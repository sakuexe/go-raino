package env

import (
	"os"

	"github.com/joho/godotenv"
)

func GetDotenv(variable string) string {
  err := godotenv.Load()

  if err != nil {
    panic("Error loading .env file at project root")
  }

  token := os.Getenv(variable)
  return token
}
