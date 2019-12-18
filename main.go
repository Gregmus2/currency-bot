package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type MonoCurrency struct {
	CurrencyA int `json:"currencyCodeA"`
	CurrencyB int `json:"currencyCodeB"`
	Date      int `json:"date"`
	RateBuy   float32 `json:"rateBuy"`
	RateSell  float32 `json:"rateSell"`
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

	res, err := http.Get("https://api.monobank.ua/bank/currency")
	if err != nil {
		bot.errorReport("Something wrong with fetching monobank api" + err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		bot.errorReport("Error on read response" + err.Error())
	}

	currencies := make([]MonoCurrency, 0)
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		bot.errorReport("Error json decode" + err.Error())
	}

	for _, currency := range currencies {
		if currency.CurrencyA == cfg.Currency {
			var msg tgbotapi.Chattable
			if currency.RateSell >= 24 {
				msg = bot.factory(fmt.Sprintf("FUCK, your deposite in ass %f", currency.RateSell))
			} else {
				msg = bot.factory(fmt.Sprintf("%f", currency.RateSell))
			}
			_, err = bot.api.Send(msg)
			if err != nil {
				bot.errorReport("Error send msg" + err.Error())
			}

			return
		}
	}
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

func (b *Bot) errorReport(msg string) {
	_, err := b.api.Send(b.factory(msg))
	if err != nil {
		log.Fatal(err)
	}
}
