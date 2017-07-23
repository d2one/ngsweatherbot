package main

//TODO кэширование ответов апи новостей, относительно прогноза на 10 минут
//TODO сделать нотификации пользователям, о выбранной погоде, по времени
//TODO Сделать вывод выбранного текущего города, и так же в прогнозе
//TODO Няшные менюшки и вообще навигация
//TODO start using docopt

import (
	"flag"
	"log"
	"os"

	"golang.org/x/net/context"

	"time"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
)

var ds *DataStore
var ws *WeatherService
var debugAPI bool
var err error
var cache CacheService

func main() {
	var telegramKey string

	//TODO rebuild on env params
	flag.StringVar(&telegramKey, "k", "", "sekret telegram api key")
	flag.BoolVar(&debugAPI, "d", false, "enable api debug mode")
	flag.Parse()
	if telegramKey == "" {
		log.Println("Secret telegram key not setted")
		os.Exit(1)
	}
	ds = NewDataStore()
	cache = NewCache()
	apiW := NewWeatherAPI()
	ws = NewWeatherService(cache, apiW)

	api := telegram.New(telegramKey)
	api.Debug(debugAPI)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic
	bot.HandleFunc(defaultCommand)
	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Use command middleware, that helps to work with commands
	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start":   telebot.CommandFunc(startCommand),
		"current": telebot.CommandFunc(currentCommand),
		"help":    telebot.CommandFunc(helpCommand),
		"city":    telebot.CommandFunc(cityCommand),
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
