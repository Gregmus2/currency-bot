package main

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ideaCurrency struct {
	Currencies []struct {
		ByRate   string `json:"byRate"`
		SellRate string `json:"sellRate"`
		Currency string `json:"currency"`
	} `json:"currencies"`
}

type Idea struct {
	Receiver
}

func (p *Idea) GetRates() error {
	jsonStr := []byte(`{"date":"` + time.Now().Format("2006-01-02T15:04:05-0700") + `"}`)
	url := "https://admin.ideabank.ua/ru/jsonapi/wise/proxyrequest?method=GetExchangeRate"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return errors.Wrap(err, "Idea New request error")
	}
	req.Header.Set("authorization", "Basic aWRlYTphV1JsWVdKaGJtcz0=")
	req.Header.Set("Content-Type", "application/json")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return errors.Wrap(err, "Idea API error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "Error on read response")
	}

	currencies := ideaCurrency{}
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		return errors.Wrap(err, "Idea: Json decode error"+string(body))
	}

	for _, currency := range currencies.Currencies {
		if currency.Currency == "USD" {
			currency.ByRate = strings.ReplaceAll(currency.ByRate, ",", ".")
			currency.SellRate = strings.ReplaceAll(currency.SellRate, ",", ".")
			p.buy, _ = strconv.ParseFloat(currency.ByRate, 64)
			p.sell, _ = strconv.ParseFloat(currency.SellRate, 64)

			return nil
		}
	}

	return errors.New("Idea have no USD in response")
}

func (p *Idea) BankName() string {
	return "Idea    "
}
