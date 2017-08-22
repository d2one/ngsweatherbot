package main

import (
	"context"
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"log"
	"os"
	"time"
)

var dataStore *DataStore
var weatherService *WeatherService
var cache *Cache
var weatherIcons *WeatherIcons
var debugAPI bool
var err error

func main() {
	var telegramKey string
	telegramKey = os.Getenv("NGS_WEATHER_BOT_TELEGRAM_KEY")
	if telegramKey == "" {
		panic("Empty telegram API KEY")
	}
	if debugAPIStatus := os.Getenv("NGS_WEATHER_BOT_DEBUG"); debugAPIStatus == "true" {
		debugAPI = true
	}

	dbPath := os.Getenv("NGS_WEATHER_BOT_DB_PATH")
	if dbPath == "" {
		panic("Empty DB PATH")
	}

	dataStore = NewDataStore(dbPath)
	weatherIcons = NewWeatherIcons()
	cache = NewCache()
	weatherService = NewWeatherService(cache)

	api := telegram.New(telegramKey)
	api.Debug(debugAPI)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic
	bot.HandleFunc(defaultCommand)
	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Use command middleware, that helps to work with commands
	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start":   telebot.CommandFunc(initCommand),
		"current": telebot.CommandFunc(currentCommand),
		"help":    telebot.CommandFunc(helpCommand),
		"city":    telebot.CommandFunc(cityCommand),
	}))

	bot.Use(telebot.Callbacks(map[string]telebot.InlineCallback{
		"/city":                     telebot.CallbackFunc(cityCommand),
		"/show_city":                telebot.CallbackFunc(currentCommand),
		"/weather_forecast_current": telebot.CallbackFunc(callbackWeatherType),
	}))
	go runCronCommands(netCtx, api)
	log.Fatal(bot.Serve(netCtx))
}

func logWork(err error) {
	log.Println(err)
}

func runCronCommands(netCtx context.Context, api *telegram.API) {
	for {
		runNotificationTasks(netCtx, api)
		time.Sleep(time.Minute)
	}
}
