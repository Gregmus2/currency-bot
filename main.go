package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strconv"
)

type RatesReceiver interface {
	GetRates() (float64, float64, error)
	BankName() string
}

type Bot struct {
	api     *tgbotapi.BotAPI
	factory func(msg string) tgbotapi.MessageConfig
}

type Config struct {
	Token     string `json:"token"`
	ChannelId string `json:"channel_id"`
	Currency  int    `json:"currency"`
}

func main() {
	cfg := loadConfiguration()
	bot := NewBot(cfg)

	banks := []RatesReceiver{
		new(Pivdeniy),
		new(Monobank),
		new(Privat),
	}

	text := "<pre>"
	for _, bank := range banks {
		buy, sell, err := bank.GetRates()
		bot.errorReport(err)
		if err == nil {
			text += fmt.Sprintf("%s:\t%.2f\t%.2f\n", bank.BankName(), buy, sell)
		}
	}
	text += "</pre>"

	msg := bot.factory(text)
	msg.ParseMode = "html"
	_, err := bot.api.Send(msg)
	bot.errorReport(err)
}

func NewBot(cfg *Config) *Bot {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	api.Debug = false

	var factory func(msg string) tgbotapi.MessageConfig
	chatID, err := strconv.Atoi(cfg.ChannelId)
	if err == nil && chatID != 0 {
		factory = func(msg string) tgbotapi.MessageConfig { return tgbotapi.NewMessage(int64(chatID), msg) }
	} else {
		factory = func(msg string) tgbotapi.MessageConfig { return tgbotapi.NewMessageToChannel(cfg.ChannelId, msg) }
	}

	bot := &Bot{
		api:     api,
		factory: factory,
	}

	return bot
}

func loadConfiguration() *Config {
	config := new(Config)
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		log.Fatalf("Error on loading configuration file: %s", err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(config)

	return config
}

func (b *Bot) errorReport(err error) {
	if err == nil {
		return
	}

	_, err = b.api.Send(b.factory(err.Error()))
	if err != nil {
		log.Fatal(err)
	}
}
