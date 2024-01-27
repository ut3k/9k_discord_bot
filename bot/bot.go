package bot

import (
  "fmt"
  "log"
  "os"
  "os/signal"
  "strings"

  "github.com/bwmarrin/discordgo"
)

var BotToken string

func checNilErr (e error){
  if e != nil {
    log.Fatal("Error message")
  }
}

func Run() {

  // create session
  discord, err := discordgo.New("Bot " +BotToken)
  checNilErr(err)

  //add Event handler
  discord.AddHandler(newMessage)

  // open session 
  discord.Open()
  defer discord.Close() // close session

  // keep bot running nitil no interrupt
  fmt.Println("Bot is running...")
  c := make (chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  <-c

}

func newMessage (discord *discordgo.Session, message *discordgo.MessageCreate){
  // prevent bot responding to self
  // achived by looking into  message authot ID
  // if msg author ID same as bot then just return

  if message.Author.ID == discord.State.User.ID {
    return
  }

  switch {
  case strings.Contains(message.Content,  "!help"):
    discord.ChannelMessageSend(message.ChannelID, "Czego chcesz?")
  fmt.Println("Wykryto wiadmość !help")
  case strings.Contains(message.Content, "!bye"):
    discord.ChannelMessageSend(message.ChannelID, "nara!")
  case strings.Contains(message.Content, "tescik"):
  discord.ChannelMessageSend(message.ChannelID, "bot odpowiada")
  fmt.Println("wykryto wiadomość TESCIK")
  }

}
