package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
)

func getCities(arg string) (*WeatherCitys, string) {
	log.Printf("call method getCities:" + arg)
	resp, err := http.Get("http://pogoda.ngs.ru/api/v1/cities?q=" + arg)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var s = new(WeatherCitys)
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Printf("whoops:", err)
	}

	if s.Errors.Message != "" {
		return s, s.Errors.Message
	}

	return s, ""
}

func getCity(arg string) (string, string) {
	s, err := getCities(arg)
	if err != "" {
		return "", err
	}
	return s.Cities[0].Alias, ""

}

func getStations(arg string) string {
	log.Printf("call method getStations:" + arg)
	resp, err := http.Get("http://pogoda.ngs.ru/api/v1/forecasts/current?city=" + arg)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var s = new(WeatherResponce)
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Printf("whoops:", err)
	}

	messageText := strconv.FormatFloat(float64(s.Forecasts[0].Temperature), 'f', 1, 32) + "°C, " + "Ветер: " + strconv.FormatFloat(float64(s.Forecasts[0].Wind.Speed), 'f', 1, 32) + "м/с, " + s.Forecasts[0].Wind.Direction.Title + ", " + s.Forecasts[0].Cloud.Title + " " + s.Forecasts[0].Precipitation.Title
	messageText = arg + " " + messageText
	return messageText
}
