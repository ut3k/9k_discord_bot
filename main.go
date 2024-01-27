package main

import (
	"log"
	"os"
  "9k_discord_bot/bot"
  

	"github.com/joho/godotenv"
)

func main () {
  err := godotenv.Load()
  if err !=nil {
    log.Fatal("Error loading .env file")
  }
  bot.BotToken = os.Getenv("BOT_TOKEN")
  bot.Run()

}
