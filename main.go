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

  // błąd po stronie discorda na boty zapraszane samemu na serwer
  // zawratość pola "Content" jest pusta, wszysto po niżej nie zadziała

  // obejscie problemu  z githuba'a też nie działa

  if m.Content == "" {
		chanMsgs, err := s.ChannelMessages(m.ChannelID, 1, "", "", m.ID)
		if err != nil {
			// log.Errorf("unable to get messages: %s", err)
			return
		}
		m.Content = chanMsgs[0].Content
		m.Attachments = chanMsgs[0].Attachments
    var tekst string
    tekst = chanMsgs[0].Content
    fmt.Println(tekst)
	}
  // omijanie własnych rekacji bota na siebie samego
  if m.Content == "aa" {
    s.ChannelMessageSend(m.ChannelID, "Pong !")
    fmt.Println("wykryto @everyone")
    fmt.Println(m.Content)
  }
  if m.Content == "pong" {
    s.ChannelMessageSend(m.ChannelID, "ping !!")
  }



}

