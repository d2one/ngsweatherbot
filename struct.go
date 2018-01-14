package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
)

// WeatherBot ds
type WeatherBot struct {
	api            *telegram.API
	Bot            *telebot.Bot
	DataStore      *DataStore
	Icons          *WeatherIcons
	WeatherService *WeatherService
	Debug          bool
	ctx            context.Context
}

// NewWeatherBot dsds
func NewWeatherBot(netCtx context.Context, telegramKey string, store *DataStore, service *WeatherService) *WeatherBot {
	api := telegram.New(telegramKey)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic
	bot.HandleFunc(defaultCommand)

	weatherBot := &WeatherBot{
		api:            api,
		Bot:            bot,
		DataStore:      store,
		Icons:          NewWeatherIcons(),
		WeatherService: service,
		Debug:          false,
		ctx:            netCtx,
	}
	return weatherBot
}

func (weatherBot WeatherBot) setDebug(enable bool) {
	weatherBot.Debug = enable
	weatherBot.api.Debug(enable)
}

//AddCommands ds
func (weatherBot WeatherBot) AddCommands(commands map[string]telebot.CommandFunc) {
	tCommands := make(map[string]telebot.Commander)
	log.Println(commands)
	for k, v := range commands {
		tCommands[k] = telebot.CommandFunc(v)
	}

	weatherBot.Bot.Use(telebot.Commands(tCommands))
}

//AddCallbacks ds
func (weatherBot WeatherBot) AddCallbacks(commands map[string]telebot.CallbackFunc) {
	tCommands := make(map[string]telebot.InlineCallback)
	for k, v := range commands {
		tCommands[k] = telebot.CallbackFunc(v)
	}
	weatherBot.Bot.Use(telebot.Callbacks(tCommands))
}

func (weatherBot WeatherBot) start() {
	log.Fatal(weatherBot.Bot.Serve(weatherBot.ctx))
}

func (weatherBot WeatherBot) startCron() {
	for {
		runNotificationTasks(weatherBot.ctx, weatherBot.api)
		time.Sleep(time.Minute)
	}
}

// UserCity user selected city
type UserCity struct {
	ChatID       int64
	CityAlias    string
	CityTitle    string
	ForecastType string
}

// WeatherSource WeatherSource
type WeatherSource interface {
	getCities(arg string) ([]*City, error)
	getCity(arg string) (*City, error)
	getCurrentWeather(arg string) (*CurrentWeather, error)
	getForecast(arg string) (*WeatherResponseForecasts, error)
}

// UserData ds
type UserData struct {
	ID                   int64
	ChatID               int64
	CityAlias            sql.NullString
	CityTitle            sql.NullString
	NotificationsNextRun sql.NullInt64
	ForecastType         string
	CreatedAt            int64
}

// CacheService CacheService
type CacheService interface {
	read(cacheKey string) (*CachedItem, error)
	write(cacheKey string, cacheValue string, ttl int64) error
	delete(cacheKey string) error
}

// UserNotification user notification
type UserNotification struct {
	ID      int64
	UserID  int64
	ChatID  int64
	NextRun int64
}

// CachedItem cache
type CachedItem struct {
	ID         int64
	CacheKey   string
	CacheValue string
	TTL        int64
	TTLLock    int64
}

// City ngs weather api city responce
type City struct {
	Alias string `json:"alias"`
	Title string `json:"title"`
}

// WeatherCities cities
type WeatherCities struct {
	Cities []*City `json:"cities"`
	Errors struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

// CurrentWeather current weather
type CurrentWeather struct {
	Astronomy struct {
		LengthDayHuman  string  `json:"length_day_human"`
		MoonIlluminated float64 `json:"moon_illuminated"`
		MoonPhase       string  `json:"moon_phase"`
		Sunrise         string  `json:"sunrise"`
		Sunset          string  `json:"sunset"`
	} `json:"astronomy"`
	MagneticStatus      string  `json:"magnetic_status"`
	FeelLikeTemperature float64 `json:"feel_like_temperature"`
	Cloud               struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"cloud"`
	Precipitation struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"precipitation"`
	Temperature float64 `json:"temperature"`
	Wind        struct {
		Direction struct {
			Value string `json:"value"`
			Title string `json:"title"`
		} `json:"direction"`
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Water []struct {
		Level struct {
			Hint  string  `json:"hint"`
			Trend float64 `json:"trend"`
			Value float64 `json:"value"`
		} `json:"level"`
		Temperature interface{} `json:"temperature"`
		Title       string      `json:"title"`
		WaveHeight  interface{} `json:"wave_height"`
	} `json:"water"`
	IconPath string `json:"icon_path"`
	Links    struct {
		City string `json:"city"`
	} `json:"links"`
}

// WeatherResponse weather response
type WeatherResponse struct {
	Forecasts []*CurrentWeather `json:"forecasts"`
}

// HourForecast ds
type HourForecast struct {
	Hour        int `json:"hour"`
	Temperature struct {
		Avg float64 `json:"avg"`
	} `json:"temperature"`
	Wind struct {
		Speed struct {
			Avg float64 `json:"avg"`
		} `json:"speed"`
		Direction struct {
			Value string `json:"value"`
			Title string `json:"title"`
		} `json:"direction"`
	} `json:"wind"`
	Cloud struct {
		Value string `json:"value"`
		Title string `json:"title"`
	} `json:"cloud"`
	Precipitation struct {
		Value string `json:"value"`
		Title string `json:"title"`
	} `json:"precipitation"`
}

// WeatherResponseForecasts ds
type WeatherResponseForecasts struct {
	Forecasts []struct {
		Date  string          `json:"date"`
		Hours []*HourForecast `json:"hours"`
		Links struct {
			City string `json:"city"`
		} `json:"links"`
	} `json:"forecasts"`
}
