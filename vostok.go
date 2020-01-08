package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type vostokCurrency struct {
	CurrencyRates []struct {
		RateType      string `json:"rateType"`
		ExchangeRates []struct {
			CurrencyShortNameFrom string  `json:"currencyShortNameFrom"`
			CurrencyShortNameTo   string  `json:"currencyShortNameTo"`
			RateBuy               float64 `json:"rateBuy"`
			RateSell              float64 `json:"rateSell"`
		} `json:"exchangeRates"`
	} `json:"currencyRates"`
}

type Vostok struct {
	Receiver
}

func (p *Vostok) GetRates() error {
	res, err := http.Get("https://bankvostok.com.ua/api/getCurrencyRates?onDate=" + time.Now().Format("2006-01-02"))
	if err != nil {
		return errors.Wrap(err, "Vostok API error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "Error on read response")
	}

	currencies := vostokCurrency{}
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		return errors.Wrap(err, "Vostok: Json decode error"+string(body))
	}

	for _, currency := range currencies.CurrencyRates {
		if currency.RateType != "cashless" {
			continue
		}

		for _, cash := range currency.ExchangeRates {
			if cash.CurrencyShortNameFrom == "USD" && cash.CurrencyShortNameTo == "UAH" {
				p.buy = cash.RateBuy
				p.sell = cash.RateSell

				return nil
			}
		}
	}

	return errors.New("Vostok have no USD in response" + string(body))
}

func (p *Vostok) BankName() string {
	return "Vostok  "
}
