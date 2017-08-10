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

func (weatherApi *WeatherAPI) getCities(arg string) ([]*City, error) {
	if len(arg) == 0 {
		return nil, errors.New("пустой город")
	}

	var city = &WeatherCities{}
	err = weatherApi.getDataByURL("/api/v1/cities?q="+arg, city)
	if err != nil {
		return nil, err
	}

	if city.Errors.Message != "" {
		return nil, errors.New(city.Errors.Message)
	}
	return city.Cities, nil
}

func (weatherApi *WeatherAPI) getCity(arg string) (*City, error) {
	cities, err := weatherApi.getCities(arg)
	if err != nil {
		return nil, err
	}
	return cities[0], nil

}

func (weatherApi *WeatherAPI) getCurrentWeather(arg string) (*CurrentWeather, error) {
	var weather = &WeatherResponse{}
	err := weatherApi.getDataByURL("/api/v1/forecasts/current?city="+arg, weather)
	if err != nil {
		return nil, err
	}

	return weather.Forecasts[0], nil
}

func (weatherApi *WeatherAPI) getForecast(arg string) (*WeatherResponseForecasts, error) {

	var weather = &WeatherResponseForecasts{}
	err := weatherApi.getDataByURL("/api/v1/forecasts/forecast?city="+arg, weather)
	if err != nil {
		return nil, err
	}
	return weather, nil
}

func (weatherApi *WeatherAPI) getDataByURL(url string, structToParse interface{}) error {
	resp, err := http.Get(weatherApi.URL + url)
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
