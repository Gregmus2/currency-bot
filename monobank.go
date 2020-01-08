package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

const USD int = 840

type monoCurrency struct {
	CurrencyA int     `json:"currencyCodeA"`
	CurrencyB int     `json:"currencyCodeB"`
	Date      int     `json:"date"`
	RateBuy   float64 `json:"rateBuy"`
	RateSell  float64 `json:"rateSell"`
}

type Monobank struct {
	Receiver
}

func (m *Monobank) GetRates() error {
	res, err := http.Get("https://api.monobank.ua/bank/currency")
	if err != nil {
		return errors.Wrap(err, "Monobank API error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "Error on read response")
	}

	currencies := make([]monoCurrency, 0)
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		return errors.Wrap(err, "Json decode error")
	}

	for _, currency := range currencies {
		if currency.CurrencyA == USD {
			m.buy = currency.RateBuy
			m.sell = currency.RateSell

			return nil
		}
	}

	return errors.New("monobank have no USD in response")
}

func (m *Monobank) BankName() string {
	return "Monobank"
}
