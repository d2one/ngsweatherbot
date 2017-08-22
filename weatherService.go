package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bot-api/telegram"
	"log"
)

// WeatherService WeatherService
type WeatherService struct {
	cache CacheService
	api   WeatherSource
}

// NewWeatherService WeatherService
func NewWeatherService(cache CacheService) *WeatherService {
	return &WeatherService{
		cache: cache,
		api:   NewWeatherAPI(),
	}
}

func (service *WeatherService) getCities(arg string) ([]*City, error) {
	return service.api.getCities(arg)
}

func (service *WeatherService) getCity(arg string) (*City, error) {
	return service.api.getCity(arg)
}

func (service *WeatherService) getCurrentWeather(cityTitle string) (*CurrentWeather, error) {
	log.Println(cityTitle)
	if cachedCurrentWeather, _ := service.cache.read("current_weather" + cityTitle); cachedCurrentWeather != nil {
		var currentWeather CurrentWeather
		json.Unmarshal([]byte(cachedCurrentWeather.CacheValue), &currentWeather)
		return &currentWeather, nil
	}

	if currentWeather, _ := service.api.getCurrentWeather(cityTitle); currentWeather != nil {
		cacheValue, _ := json.Marshal(currentWeather)
		service.cache.write("current_weather"+cityTitle, string(cacheValue), time.Now().Unix()+60*10)
		log.Println(currentWeather)
		return currentWeather, nil
	}

	return nil, nil
}

func (service *WeatherService) getForecast(arg string) (*WeatherResponseForecasts, error) {
	if cachedForecast, _ := service.cache.read("forecast_weather" + arg); cachedForecast != nil {
		var forecast WeatherResponseForecasts
		json.Unmarshal([]byte(cachedForecast.CacheValue), &forecast)
		return &forecast, nil
	}

	if forecast, _ := service.api.getForecast(arg); forecast != nil {
		cacheValue, _ := json.Marshal(forecast)
		service.cache.write("forecast_weather"+arg, string(cacheValue), time.Now().Unix()+60*10)
		return forecast, nil
	}
	return nil, nil
}

func (service *WeatherService) formatCurrentWeather(weather *CurrentWeather) (string, telegram.InlineKeyboardMarkup) {
	var inlineKeyboard = [][]telegram.InlineKeyboardButton{}
	messageText := fmt.Sprintf("%g °C, %s %g м/с, %s %s\n [Источник](https://pogoda.ngs.ru/%s)",
		weather.Temperature,
		weatherIcons.getWind(weather.Wind.Direction.Value),
		weather.Wind.Speed,
		weatherIcons.getClouds(weather.Cloud.Value),
		weatherIcons.getPrecipitations(weather.Precipitation.Value),
		weather.Links.City,
	)

	inlineKeyboard = [][]telegram.InlineKeyboardButton{
		{
			{
				Text:         "Подробно",
				CallbackData: "/weather_forecast_current:full",
			},
		},
	}

	replyMarkup := telegram.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
	return messageText, replyMarkup
}

func (service *WeatherService) formatFullCurrentWeather(weather *CurrentWeather) (string, telegram.InlineKeyboardMarkup) {
	var inlineKeyboard = [][]telegram.InlineKeyboardButton{}
	messageText := fmt.Sprintf("*Температура:* %g°C, ощущается как  %g°C. \n"+
		"*Ветер:* %s %g м/с. \n%s, %s. \n%s, %s - %s\n%s\n"+
		"[Источник](https://pogoda.ngs.ru/%s)",
		weather.Temperature,
		weather.FeelLikeTemperature,
		weather.Wind.Direction.Title,
		weather.Wind.Speed,
		weather.Cloud.Title,
		weather.Precipitation.Title,
		weatherIcons.getAstronomy("sunrise")+" "+weather.Astronomy.Sunrise,
		weatherIcons.getAstronomy("sunset")+" "+weather.Astronomy.Sunset,
		weather.Astronomy.LengthDayHuman,
		weather.MagneticStatus,
		weather.Links.City,
	)

	inlineKeyboard = [][]telegram.InlineKeyboardButton{
		{
			{
				Text:         "Кратко",
				CallbackData: "/weather_forecast_current:short",
			},
		},
	}

	replyMarkup := telegram.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
	return messageText, replyMarkup
}

func (service *WeatherService) formatForecastWeather(weather *WeatherResponseForecasts) string {
	messageText := ""
	t := time.Now()
	h := t.Hour()
	counter := 0
	cityAlias := ""
	for indexForecast, forecast := range weather.Forecasts {
		cityAlias = forecast.Links.City
		if indexForecast > 3 {
			break
		}
		if indexForecast > 0 {
			t, _ := time.Parse("2006-01-02T15:04:05-0700", forecast.Date)
			messageText = messageText + "\n*" + t.Format("2 Jan") + "*\n"
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
	messageText = messageText + fmt.Sprintf("\n[Источник](https://pogoda.ngs.ru/%s)", cityAlias)
	return messageText
}

func formatForecast(hourForecast *HourForecast) string {
	return fmt.Sprintf("%v °C, Ветер: %v м/с, %s %s %s \n",
		hourForecast.Temperature.Avg,
		hourForecast.Wind.Speed.Avg,
		weatherIcons.getWind(hourForecast.Wind.Direction.Value),
		weatherIcons.getClouds(hourForecast.Cloud.Value),
		weatherIcons.getPrecipitations(hourForecast.Precipitation.Value),
	)
}
