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

// WeatherAPI db
type WeatherAPI struct {
	URL string
}

func (api *WeatherAPI) init() {
	api.URL = "https://pogoda.ngs.ru"
}

func (api *WeatherAPI) getCities(arg string) ([]*City, error) {
	body, err := api.getDataByURL("/api/v1/cities?q=" + arg)
	if err != nil {
		return nil, err
	}

	var city = new(WeatherCities)
	if err = json.Unmarshal(body, &city); err != nil {
		return nil, err
	}

	if city.Errors.Message != "" {
		return nil, errors.New(city.Errors.Message)
	}
	return city.Cities, nil
}

func (api *WeatherAPI) getCity(arg string) (*City, error) {
	cities, err := api.getCities(arg)
	if err != nil {
		return nil, err
	}
	return cities[0], nil

}

func (api *WeatherAPI) getCurrentWeather(arg string) (*CurrentWeather, error) {
	log.Println("getCurrentWeather" + arg)
	cachedCurrentWeather, _ := readCache("current_weather" + arg)
	log.Println("CACHE READED getCurrentWeather" + arg)
	if cachedCurrentWeather != nil {
		currentWeather := new(CurrentWeather)
		json.Unmarshal([]byte(cachedCurrentWeather.CacheValue), &currentWeather)
		return currentWeather, nil
	}
	log.Println("DATAURL" + arg)
	body, err := api.getDataByURL("/api/v1/forecasts/current?city=" + arg)
	if err != nil {
		return nil, err
	}
	var weather = new(WeatherResponce)
	if err = json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}
	currentWeather := weather.Forecasts[0]
	cacheValue, _ := json.Marshal(currentWeather)
	log.Println("SAVE CACHE" + arg)
	saveCache("current_weather"+arg, string(cacheValue), time.Now().Unix()+60*10)
	return currentWeather, nil
}

func (api *WeatherAPI) getDataByURL(url string) ([]byte, error) {
	resp, err := http.Get(api.URL + url)
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

func (api *WeatherAPI) formatCurrentWeather(weather *CurrentWeather) string {
	messageText := fmt.Sprintf("%g °C, Ветер: %g м/с, %s %s %s",
		weather.Temperature,
		weather.Wind.Speed,
		weather.Wind.Direction.Title,
		weather.Cloud.Title,
		weather.Precipitation.Title,
	)
	return messageText
}
