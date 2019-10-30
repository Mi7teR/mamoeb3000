package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	ga "github.com/jpillora/go-ogle-analytics"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Templates struct {
	Curses `json:"curses"`
}
type Curses []string

var data []byte
var templates Templates
var cursesSize int
var gaKey, hostName, locationUrl string

func init() {
	fileData, err := ioutil.ReadFile("templates.json")
	if err != nil {
		log.Fatalln(err)
	}
	data = fileData
	if err := json.Unmarshal(data, &templates); err != nil {
		panic(err)
	}
	cursesSize = len(templates.Curses)
}

func sendGaEvent(clientID, messageSource, messageType string) {
	if len(gaKey) == 0 {
		return
	}
	client, err := ga.NewClient(gaKey)
	if err != nil {
		log.Println(err)
	}
	client.DocumentHostName(hostName)
	client.DocumentLocationURL(fmt.Sprintf("%s/%s/%s", locationUrl, messageSource, messageType))
	client.UserID(clientID)
	err = client.Send(ga.NewPageview())
	if err != nil {
		log.Println(err)
	}
	return
}

func main() {
	var token string
	flag.StringVar(&token, "token", "empty", "telegram bot token")
	flag.StringVar(&hostName, "hostName", "bot.mi7ter.xyz", "ga hostname")
	flag.StringVar(&locationUrl, "locationUrl", "https://bot.mi7ter.xyz", "ga location url")
	flag.StringVar(&gaKey, "gaKey", "", "google analytics key")
	flag.Parse()
	if token == "empty" {
		panic(getRandomCurse())
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			var messageText string
			messageText = getRandomCurse()
			message := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			_, err = bot.Send(message)
			if err != nil {
				sendGaEvent(string(update.Message.From.ID), "direct", "error")
			} else {
				sendGaEvent(string(update.Message.From.ID), "direct", "curse")
			}
			continue
		}

		if update.InlineQuery == nil {
			continue
		}

		curse := getRandomCurse()
		var articles []interface{}
		article := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, "Выебать мамку", curse)
		article.Description = "И не только мамку"
		articles = append(articles, article)

		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			CacheTime:     0,
			Results:       articles,
		}

		if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
			sendGaEvent(string(update.InlineQuery.From.ID), "inline", "error")
		} else {
			sendGaEvent(string(update.InlineQuery.From.ID), "inline", "curse")
		}
	}
}

func getRandomCurse() string {
	curse := templates.Curses[rand.Intn(cursesSize)]
	return curse
}
