package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Pivdeniy byte

func (p Pivdeniy) GetRates() (float64, float64, error) {
	now := time.Now()
	unixTime := strconv.Itoa(int(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()))
	res, err := http.Get("https://bank.com.ua/api/ru/v1/rest-ui/find-branch-course?date=" + unixTime)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Pivdeniy API error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Error on read response")
	}

	currencies := make([][]interface{}, 6)
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Pivdeniy: Json decode error"+string(body))
	}

	for _, currency := range currencies {
		title, ok := currency[1].(string)
		if !ok {
			return 0, 0, errors.Wrap(err, "Wrong response struct"+string(body))
		}

		if title != "USD" {
			continue
		}

		buy, ok := currency[2].(float64)
		if !ok {
			return 0, 0, errors.Wrap(err, "Wrong response struct: "+string(body))
		}

		sell, ok := currency[3].(float64)
		if !ok {
			return 0, 0, errors.Wrap(err, "Wrong response struct: "+string(body))
		}

		return buy, sell, nil
	}

	return 0, 0, errors.New("pivdeniy have no USD in response")
}

func (p Pivdeniy) BankName() string {
	return "Pivdeniy"
}
