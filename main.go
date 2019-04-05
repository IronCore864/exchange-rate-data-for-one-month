package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ironcore864/exchange-rate-data-for-one-month/model"
	"github.com/ironcore864/exchange-rate-data-for-one-month/redisclient"
)

func main() {
	netClient := &http.Client{
		Timeout: time.Second * 1,
	}
	// here the url should be /latest
	// but at the moment it does not work, so I changed to /2019-04-05
	url := "https://api.exchangeratesapi.io/2019-04-05"

	res, err := netClient.Get(url)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var rates model.Rates

	if err := json.Unmarshal(body, &rates); err != nil {
		panic(err.Error())
	}

	date := rates.Date

	for currency, rate := range rates.Rates {
		redisKey := fmt.Sprintf("%s-%s", currency, date)
		_, err := redisclient.Set(redisKey, rate, time.Second*86400*30)
		if err != nil {
			log.Fatal(err)
		}
	}
}
