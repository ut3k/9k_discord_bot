package main

import (
	"log"
	"os"
  "fmt"
  "encoding/binary"
  "io"
  "os/signal"
  "strings"
  "syscall"
  "time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)
func init(){
  // load ENV
  err := godotenv.Load()
  if err !=nil {
    log.Fatal("Error loading .env file")
  }
}

var  BotToken string = os.Getenv("BOT_TOKEN")

func main () {
  
  // load sound effect
  err := loadSound()
  if err != nil{
    fmt.Println(" Error loading sound", err)
    return
  }


  // New discord session
  dg, err := discordgo.New("Bot " + BotToken)
  if err != nil {
    fmt.Println("Error creating Discord SESSION: ", err)
    return
  }

  // Register READY as a call back for the ready events
  dg.AddHandler(ready)

  // Register messageCreate as a callback for messageCreate events
  dg.AddHandler(messageCreate)
  
  // Register guildCreate as a callback for the GuildCreate events
  dg.AddHandler(guildCreate)

  // info about guilds and therir channels
  // messages and voice states

  dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

  // open websocet and beging listening
  err = dg.Open()
  if err !=nil {
    fmt.Println("Error opening Discord session: ", err)
  }


  // wait here until CTRL-C or other term signal is recived.
  fmt.Println("Airhorn is now running. Press CTRL-C to exit.")
  sc :=  make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
  <-sc

  // close session
  dg.Close()

}

func ready (s *discordgo.Session, event *discordgo.Ready) {

  // set the playing status.
  s.UpdateGameStatus(0, "!airhorn")

}
// pauza 
 // https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go 
