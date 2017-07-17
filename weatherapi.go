package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const POGODA_URL = "http://pogoda.ngs.ru"

func getCities(arg string) ([]City, error) {
	body, err := getDataByUrl("/api/v1/cities?q=" + arg)
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

	cachedCurrentWeather, _ := readCache("current_weather" + arg)
	if cachedCurrentWeather != (CachedItem{}) {
		currentWeather := new(CurrentWeather)
		json.Unmarshal([]byte(cachedCurrentWeather.cache_value), &currentWeather)
		return *currentWeather, nil
	}

	body, err := getDataByUrl("/api/v1/forecasts/current?city=" + arg)
	if err != nil {
		return CurrentWeather{}, err
	}
	var weather = new(WeatherResponce)
	if err = json.Unmarshal(body, &weather); err != nil {
		return CurrentWeather{}, err
	}
	currentWeather := weather.Forecasts[0]

	cacheValue, _ := json.Marshal(currentWeather)
	saveCache("current_weather"+arg, string(cacheValue), time.Now().Unix()+60*10)
	return currentWeather, nil
}

func getDataByUrl(url string) ([]byte, error) {
	resp, err := http.Get(POGODA_URL + url)
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

func formatCurrentWeather(weather CurrentWeather) string {
	messageText := fmt.Sprintf("%g °C, Ветер: %g м/с, %s %s %s",
		weather.Temperature,
		weather.Wind.Speed,
		weather.Wind.Direction.Title,
		weather.Cloud.Title,
		weather.Precipitation.Title,
	)
	return messageText
}
