package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// WeatherService WeatherService
type WeatherService struct {
	cache         CacheService
	weatherSource WeatherSource
}

// NewWeatherService WeatherService
func NewWeatherService(cache CacheService, weatherSource WeatherSource) *WeatherService {
	return &WeatherService{
		cache:         cache,
		weatherSource: weatherSource,
	}
}

func (wc *WeatherService) getCities(arg string) ([]*City, error) {
	return wc.weatherSource.getCities(arg)
}

func (wc *WeatherService) getCity(arg string) (*City, error) {
	return wc.weatherSource.getCity(arg)
}

func (wc *WeatherService) getCurrentWeather(arg string) (*CurrentWeather, error) {
	if cachedCurrentWeather, _ := wc.cache.read("current_weather" + arg); cachedCurrentWeather != nil {
		currentWeather := &CurrentWeather{}
		json.Unmarshal([]byte(cachedCurrentWeather.CacheValue), &currentWeather)
		return currentWeather, nil
	}

	if currentWeather, _ := wc.weatherSource.getCurrentWeather(arg); currentWeather != nil {
		cacheValue, _ := json.Marshal(currentWeather)
		wc.cache.write("current_weather"+arg, string(cacheValue), time.Now().Unix()+60*10)
		return currentWeather, nil
	}
	return nil, nil
}

func (wc *WeatherService) formatCurrentWeather(weather *CurrentWeather) string {
	messageText := fmt.Sprintf("%g °C, Ветер: %g м/с, %s %s %s",
		weather.Temperature,
		weather.Wind.Speed,
		weather.Wind.Direction.Title,
		weather.Cloud.Title,
		weather.Precipitation.Title,
	)
	return messageText
}
