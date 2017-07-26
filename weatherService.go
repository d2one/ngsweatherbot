package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func (wc *WeatherService) getForecast(arg string) (*WeatherResponceForecasts, error) {
	if cachedForecast, _ := wc.cache.read("forecast_weather" + arg); cachedForecast != nil {
		forecast := &WeatherResponceForecasts{}
		json.Unmarshal([]byte(cachedForecast.CacheValue), &forecast)
		return forecast, nil
	}

	if forecast, _ := wc.weatherSource.getForecast(arg); forecast != nil {
		cacheValue, _ := json.Marshal(forecast)
		wc.cache.write("forecast_weather"+arg, string(cacheValue), time.Now().Unix()+60*10)
		return forecast, nil
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

func (wc *WeatherService) formatForecasttWeather(weather *WeatherResponceForecasts) string {
	messageText := ""
	t := time.Now()
	h := t.Hour()
	counter := 0
	for indexForecast, forecast := range weather.Forecasts {

		if indexForecast > 1 || counter > 1 {
			break
		}

		for _, hourForecast := range forecast.Hours {
			if hourForecast.Hour < h && indexForecast < 1 {
				continue
			}
			counter++
			switch hourForecast.Hour {
			case 0:
				messageText = messageText + "*Ночью:* " + formatForecast(hourForecast)
				continue
			case 6:
				if indexForecast > 0 {
					messageText = messageText + "Завтра: \n"
				}
				messageText = messageText + "*Утром:* " + formatForecast(hourForecast)
				continue
			case 12:
				messageText = messageText + "*Днем:* " + formatForecast(hourForecast)
				continue
			case 18:
				messageText = messageText + "*Вечером:* " + formatForecast(hourForecast)
				continue
			}
		}

	}
	return messageText
}

func formatForecast(hourForecast *HourForecat) string {
	log.Println(hourForecast.Temperature)
	return fmt.Sprintf("%v °C, Ветер: %v м/с, %s %s %s \n",
		hourForecast.Temperature.Avg,
		hourForecast.Wind.Speed.Avg,
		hourForecast.Wind.Direction.Title,
		hourForecast.Cloud.Title,
		hourForecast.Precipitation.Title,
	)
}
