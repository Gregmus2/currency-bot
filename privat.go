package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type privatCurrency struct {
	Currency     string `json:"ccy"`
	BaseCurrency string `json:"base_ccy"`
	Buy          string `json:"buy"`
	Sale         string `json:"sale"`
}

type Privat struct {
	Receiver
}

func (p *Privat) GetRates() error {
	res, err := http.Get("https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5")
	if err != nil {
		return errors.Wrap(err, "Privat API error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "Error on read response")
	}

	currencies := make([]privatCurrency, 0)
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		return errors.Wrap(err, "Privat: Json decode error")
	}

	for _, currency := range currencies {
		if currency.Currency == "USD" && currency.BaseCurrency == "UAH" {
			p.buy, _ = strconv.ParseFloat(currency.Buy, 64)
			p.sell, _ = strconv.ParseFloat(currency.Sale, 64)

			return nil
		}
	}

	return errors.New("Privat have no USD in response")
}

func (p *Privat) BankName() string {
	return "Privat  "
}
