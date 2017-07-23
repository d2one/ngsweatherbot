package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// weather.WeatherAPI db
type WeatherAPI struct {
	URL string
}

// dasdas
func NewWeatherAPI() *WeatherAPI {
	return &WeatherAPI{URL: "https://pogoda.ngs.ru"}
}

func (Weatherapi *WeatherAPI) getCities(arg string) ([]*City, error) {
	if len(arg) == 0 {
		return nil, errors.New("пустой город")
	}

	log.Println("city " + arg)
	body, err := Weatherapi.getDataByURL("/api/v1/cities?q=" + arg)
	if err != nil {
		return nil, err
	}

	var city = &WeatherCities{}
	if err = json.Unmarshal(body, &city); err != nil {
		return nil, err
	}

	if city.Errors.Message != "" {
		return nil, errors.New(city.Errors.Message)
	}
	return city.Cities, nil
}

func (Weatherapi *WeatherAPI) getCity(arg string) (*City, error) {
	cities, err := Weatherapi.getCities(arg)
	if err != nil {
		return nil, err
	}
	return cities[0], nil

}

func (Weatherapi *WeatherAPI) getCurrentWeather(arg string) (*CurrentWeather, error) {
	body, err := Weatherapi.getDataByURL("/api/v1/forecasts/current?city=" + arg)
	if err != nil {
		return nil, err
	}
	var weather = &WeatherResponce{}
	if err = json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}
	currentWeather := weather.Forecasts[0]
	return currentWeather, nil
}

func (Weatherapi *WeatherAPI) getForecast(arg string) (*WeatherResponceForecasts, error) {
	body, err := Weatherapi.getDataByURL("/api/v1/forecasts/forecast?city=" + arg)
	if err != nil {
		return nil, err
	}
	var weather = &WeatherResponceForecasts{}
	if err = json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}
	return weather, nil
}

func (Weatherapi *WeatherAPI) getDataByURL(url string) ([]byte, error) {
	resp, err := http.Get(Weatherapi.URL + url)
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
