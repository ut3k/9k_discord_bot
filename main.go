package main

import (
	"log"
	"os"
  "fmt"

	// "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main () {
  err := godotenv.Load()
  if err !=nil {
    log.Fatal("Error loading .env file")
  }
  BotToken := os.Getenv("BOT_TOKEN")

  // discord, err := discordgo.New("Bot " + BotToken)
  fmt.Println(BotToken)
}
