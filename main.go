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
var buffer = make ([][]byte, 0)

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

func messageCreate (s *discordgo.Session, m *discordgo.MessageCreate) {

  // ignore bot sefl msg
  if m.Author.ID == s.State.User.ID {
    return
  }
  
 // check if msg is "!airHorn"
  if strings.HasPrefix(m.Content, "!airhorn"){

    // find chanel msg came from
    c, err := s.State.Channel(m.ChannelID)
    if err != nil {
      return
    }

    // find guild channel
    g, err := s.State.Guild(c.GuildID)
    if err != nil {
      return
    }

 // look for the message sender in that current guild current voice states
    for _, vs := range g.VoiceStates {

      if vs.UserID == m.Author.ID {
        err = playSound(s, g.ID, vs.ChannelID)

        if err != nil{
          fmt.Println("error playing sound: ", err)
        }
        return




      }
    }

  }

}

func guildCreate (s *discordgo.Session, event *discordgo.GuildCreate) {

  if event.Guild.Unavailable {
    return
  }

  for _, channel := range event.Guild.Channels {

    if channel.ID == event.Guild.ID{
      _, _ = s.ChannelMessageSend(channel.ID, "Airhorn is ready! type !airhorn while in a voice channel to play sound")
      return
    }

  }

}


func loadSound() error{

  file, err := os.Open("airhorn.dca")
  if err != nil {
    fmt.Println("error opening DCA file: ", err)
    return err

  }

  var opuslen int16

  for {

    // read opus frame lenght from dca file
    err = binary.Read(file, binary.LittleEndian, &opuslen)


    // if this is the end of file, just return.
    if err == io.EOF || err == io.ErrUnexpectedEOF {
      err := file.Close()

      if err != nil {
        return err
      }

      return nil

    }

    if err != nil {
      fmt.Println("erro reading from DCA file: ", err)
      return err

    }

    // read encoded pcm from DCA file
    InBuf := make([]byte, opuslen)
    err = binary.Read(file, binary.LittleEndian, &InBuf)
    
    // there should not be any of file errors
    if err != nil {
      fmt.Println("Error reading form DCA file: ", err)
      return err
    }

    buffer = append(buffer, InBuf)
  }


}

func playSound(s *discordgo.Session, GuildID, ChannelID string)(err error) {

  // join the provided voice channel
  vc, err := s.ChannelVoiceJoin(GuildID, ChannelID, false, true)
  if err != nil {
    return err
  }

  // sleep for specific amount of time before speaking
  time.Sleep(250 * time.Millisecond)

  // speak
  vc.Speaking(true)

  // send buffer data
  for _, buff := range buffer{
    vc.OpusSend <- buff
  }

  // stop speaking
  vc.Speaking(false)

  // sleep for specific amount of time before ending
  time.Sleep(250 * time.Millisecond)

  // disconnect from voice
  vc.Disconnect()

  return nil


}
