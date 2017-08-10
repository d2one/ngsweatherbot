package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bot-api/telegram"
)

// WeatherService WeatherService
type WeatherService struct {
	cache         CacheService
	weatherSource WeatherSource
}

// NewWeatherService WeatherService
func NewWeatherService(cache CacheService) *WeatherService {
	return &WeatherService{
		cache:         cache,
		weatherSource: NewWeatherAPI(),
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

func (wc *WeatherService) getForecast(arg string) (*WeatherResponseForecasts, error) {
	if cachedForecast, _ := wc.cache.read("forecast_weather" + arg); cachedForecast != nil {
		forecast := &WeatherResponseForecasts{}
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

func (wc *WeatherService) formatCurrentWeather(weather *CurrentWeather) (string, telegram.InlineKeyboardMarkup) {
	var inlineKeyboard = [][]telegram.InlineKeyboardButton{}
	messageText := fmt.Sprintf("%g °C, %s %g м/с, %s %s\n [pogoda.ngs.ru](https://pogoda.ngs.ru/%s)",
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

func (wc *WeatherService) formatFullCurrentWeather(weather *CurrentWeather) (string, telegram.InlineKeyboardMarkup) {
	var inlineKeyboard = [][]telegram.InlineKeyboardButton{}
	// strings.Replace(weather.IconPath, "small", "big-icons", -1),
	messageText := fmt.Sprintf(`*Температура:* %g°C, ощущается как  %g°C.
*Ветер:* %s %g м/с. 
%s, %s.
%s, %s - %s
%s
[pogoda.ngs.ru](https://pogoda.ngs.ru/%s)`,
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

func (wc *WeatherService) formatForecasttWeather(weather *WeatherResponseForecasts) string {
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

func formatForecast(hourForecast *HourForecast) string {
	return fmt.Sprintf("%v °C, Ветер: %v м/с, %s %s %s \n",
		hourForecast.Temperature.Avg,
		hourForecast.Wind.Speed.Avg,
		weatherIcons.getWind(hourForecast.Wind.Direction.Value),
		weatherIcons.getClouds(hourForecast.Cloud.Value),
		weatherIcons.getPrecipitations(hourForecast.Precipitation.Value),
	)
}
