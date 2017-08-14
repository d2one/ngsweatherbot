package main

import (
	"context"
	"flag"
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

	//TODO rebuild on env params
	flag.StringVar(&telegramKey, "k", "", "secret telegram api key")
	flag.BoolVar(&debugAPI, "d", false, "enable api debug mode")
	flag.Parse()
	if telegramKey == "" {
		log.Println("Secret telegram key not setted")
		os.Exit(1)
	}
	dataStore = NewDataStore()
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
