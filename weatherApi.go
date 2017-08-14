package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// WeatherAPI db
type WeatherAPI struct {
	URL string
}

// dasdas
func NewWeatherAPI() *WeatherAPI {
	return &WeatherAPI{URL: "https://pogoda.ngs.ru"}
}

func (api *WeatherAPI) getCities(city string) ([]*City, error) {
	if len(city) == 0 {
		return nil, errors.New("пустой город")
	}
	var weatherCities WeatherCities
	err = api.fetchStruct("/api/v1/cities?q="+city, &weatherCities)
	if err != nil {
		return nil, err
	}

	if weatherCities.Errors.Message != "" {
		return nil, errors.New(weatherCities.Errors.Message)
	}
	return weatherCities.Cities, nil
}

func (api *WeatherAPI) getCity(city string) (*City, error) {
	cities, err := api.getCities(city)
	if err != nil {
		return nil, err
	}
	if len(cities) == 0 {
		return nil, nil
	}
	return cities[0], nil

}

func (api *WeatherAPI) getCurrentWeather(arg string) (*CurrentWeather, error) {
	var weatherResponse WeatherResponse
	err := api.fetchStruct("/api/v1/forecasts/current?city="+arg, &weatherResponse)
	if err != nil {
		return nil, err
	}

	if len(weatherResponse.Forecasts) == 0 {
		return nil, nil
	}

	return weatherResponse.Forecasts[0], nil
}

func (api *WeatherAPI) getForecast(arg string) (*WeatherResponseForecasts, error) {
	var weatherResponseForecasts WeatherResponseForecasts
	err := api.fetchStruct("/api/v1/forecasts/forecast?city="+arg, &weatherResponseForecasts)
	if err != nil {
		return nil, err
	}
	return &weatherResponseForecasts, nil
}

func (api *WeatherAPI) fetchStruct(url string, structToParse interface{}) error {
	resp, err := http.Get(api.URL + url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, structToParse); err != nil {
		return err
	}
	return nil
}
