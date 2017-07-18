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

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"github.com/jasonlvhit/gocron"
	"golang.org/x/net/context"
)

var db DB
var weatherAPI WeatherAPI
var err error

func main() {
	var telegramKey string
	var debugAPI bool

	//TODO rebuild on env params
	flag.StringVar(&telegramKey, "k", "", "sekret telegram api key")
	flag.BoolVar(&debugAPI, "d", false, "enable api debug mode")
	flag.Parse()
	if telegramKey == "" {
		log.Println("Secret telegram key not setted")
		os.Exit(1)
	}

	err = db.init()
	if err != nil {
		log.Println("error init database")
		os.Exit(1)
	}

	weatherAPI.init()

	api := telegram.New(telegramKey)
	api.Debug(debugAPI)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(defaultCommand)

	// Use command middleware, that helps to work with commands
	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start":   telebot.CommandFunc(startCommand),
		"current": telebot.CommandFunc(currentCommand),
		"help":    telebot.CommandFunc(helpCommand),
		"city":    telebot.CommandFunc(cityCommand),
	}))

	gocron.Every(5).Seconds().Do(runNotificationTasks, netCtx, api)
	//<-gocron.Start()

	log.Fatal(bot.Serve(netCtx))
}
