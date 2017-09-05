package main

import (
  "os"
  "log"
  "fmt"
  "encoding/json"
  "github.com/jasonlvhit/gocron"
  "gopkg.in/telegram-bot-api.v4"
  h "github.com/qube81/hackernews-api-go"
)

type Configuration struct {
  BotToken string
}

type Stats struct {
  prevStoryId int
}

func (self Stats) getPrevTopStoryId() int {
  return self.prevStoryId
}

func (self *Stats) updatePrevTopStoryId(id int) {
  self.prevStoryId = id
}

func getTopStory() h.Story {
  topStories, _ := h.GetStories("top")
  topStory, _ := h.GetItem(topStories[0])
  return topStory
}

func sendTopStoryToChannel(bot *tgbotapi.BotAPI, stats *Stats) {
  topStory := getTopStory()
  prevStoryId := stats.getPrevTopStoryId()

  if topStory.ID != prevStoryId {
    log.Printf("Sending top story to channel...")

    stats.updatePrevTopStoryId(topStory.ID)
    msg := tgbotapi.NewMessageToChannel("@top_hacker_news", topStory.URL)
    bot.Send(msg)
  } else {
    log.Printf("Same top story ID: %d, no message sent.", prevStoryId)
  }
}

func main() {
  file, _ := os.Open("config.json")
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  err := decoder.Decode(&configuration)

  if err != nil {
    fmt.Println("error:", err)
  }

  bot, err := tgbotapi.NewBotAPI(configuration.BotToken)

  if err != nil {
    log.Panic(err)
  }

  bot.Debug = true

  log.Printf("Authorized on account %s", bot.Self.UserName)

  topStory := getTopStory()
  stats := &Stats{topStory.ID}
  gocron.Every(5).Minutes().Do(sendTopStoryToChannel, bot, stats)
  <- gocron.Start()
}