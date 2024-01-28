package main

import (
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"

  "github.com/bwmarrin/discordgo"
  "github.com/joho/godotenv"
)



func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  var BotToken string

  BotToken = os.Getenv("BOT_TOKEN")


  dg, err := discordgo.New("Bot " + BotToken)
  if err != nil {
    fmt.Println("error creating Discord session")
    return
  }

  // rejestrowanie jakiejś funkcji, do której będzie odwołanie
  dg.AddHandler(messageCreate)

  // to nie wiem co robi.
  dg.Identify.Intents = discordgo.IntentGuildMessages

  // otwarcie sesji i wywołanie ewentualnego błędu
  err = dg.Open()
  if err != nil {
    fmt.Println("nie udało się otowrzyć połączenia, bo:", err)
    return
  }

  // wywołanie przerwania, wyłączenia z terminalu bota
  fmt.Println("Bot is running now... press CTR+C to stop")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
  <-sc

  dg.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate){

  // omijanie własnych rekacji bota na siebie samego
  if m.Content == "ping" {
    s.ChannelMessageSend(m.ChannelID, "Pong !")
  }
  if m.Content == "pong" {
    s.ChannelMessageSend(m.ChannelID, "ping !!")
  }

}

