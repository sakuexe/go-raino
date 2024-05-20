package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetDotenv(variable string) string {
  err := godotenv.Load()

  if err != nil {
    panic("Error loading .env file at project root")
  }

  token := os.Getenv(variable)
  if token == "" {
    errorMessage := fmt.Sprintf("Error: %s not found in .env file", variable)
    panic(errorMessage)
  }

  return token
}
