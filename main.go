package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/bot-api/telegram/telebot"
)

var dataStore *DataStore
var weatherService *WeatherService
var cache *Cache
var weatherIcons *WeatherIcons
var err error

func main() {

	args, err := getArgs()
	if err != nil {
		panic(err)
	}

	dataStore = NewDataStore(args.dbPath)
	weatherIcons = NewWeatherIcons()
	cache = NewCache()
	weatherService = NewWeatherService(cache)
	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	weatherBot := NewWeatherBot(netCtx, args.telegramKey, dataStore, weatherService)
	weatherBot.AddCommands(map[string]telebot.CommandFunc{
		"start":   initCommand,
		"current": currentCommand,
		"help":    helpCommand,
		"city":    cityCommand,
	})

	weatherBot.AddCallbacks(map[string]telebot.CallbackFunc{
		"/city":                     cityCommand,
		"/show_city":                currentCommand,
		"/weather_forecast_current": callbackWeatherType,
	})

	weatherBot.start()

	go weatherBot.startCron()

}

// Args ds
type Args struct {
	telegramKey string
	debugAPI    bool
	dbPath      string
}

func getArgs() (*Args, error) {
	args := &Args{
		debugAPI: false,
	}

	args.telegramKey = os.Getenv("NGS_WEATHER_BOT_TELEGRAM_KEY")
	if args.telegramKey == "" {
		return nil, errors.New("Empty telegram API KEY")
	}

	if debugAPIStatus := os.Getenv("NGS_WEATHER_BOT_DEBUG"); debugAPIStatus == "true" {
		args.debugAPI = true
	}

	args.dbPath = os.Getenv("NGS_WEATHER_BOT_DB_PATH")
	if args.dbPath == "" {
		return nil, errors.New("Empty DB PATH")
	}
	return args, nil
}

func logWork(err error) {
	log.Println(err)
}
