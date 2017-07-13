package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"errors"
	"fmt"
)

func getCities(arg string) ([]City, error) {
	log.Printf("call method getCities:" + arg)
	body, err := getDataByUrl("http://pogoda.ngs.ru/api/v1/cities?q=" + arg)
	if err != nil {
		return []City{}, err
	}

	var city = new(WeatherCitys)
	if err = json.Unmarshal(body, &city); err != nil {
		return nil, err
	}

	if city.Errors.Message != "" {
		return nil, errors.New(city.Errors.Message)
	}
	return city.Cities, nil
}

func getCity(arg string) (City, error) {
	cities, err := getCities(arg)
	if err != nil {
		return City{}, err
	}
	return cities[0], nil

}

func getCurrentWeather(arg string) (CurrentWeather, error) {
	log.Printf("call method getCurrentWeather:" + arg)
	//TODO сделать кэш на 10 минут
	body, err := getDataByUrl("http://pogoda.ngs.ru/api/v1/forecasts/current?city=" + arg)
	if err != nil {
		return CurrentWeather{}, err
	}
	var weather = new(WeatherResponce)
	if err = json.Unmarshal(body, &weather); err != nil {
		return CurrentWeather{}, err
	}
	currentWeather := weather.Forecasts
	return currentWeather[0], nil
}

func getDataByUrl(url string) ([]byte, error)  {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func formatCurrentWeather(weather CurrentWeather) string  {
	messageText := fmt.Sprintf("%g °C, Ветер: %g м/с, %s %s %s",
		weather.Temperature,
		weather.Wind.Speed,
		weather.Wind.Direction.Title,
		weather.Cloud.Title,
		weather.Precipitation.Title,
	)
	return messageText
}