package env

import (
	"os"

	"github.com/joho/godotenv"
)

func GetDotenv(variable string) string {
  // it is okay if the .env file is not found
  // we will check the environment variables and
  // return an empty string if not found, it is up to the
  // caller to handle the empty string
  godotenv.Load()

  token := os.Getenv(variable)
  return token
}
