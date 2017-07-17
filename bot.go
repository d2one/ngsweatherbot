package main

//TODO кэширование ответов апи новостей, относительно прогноза на 10 минут
//TODO сделать нотификации пользователям, о выбранной погоде, по времени
//TODO Сделать вывод выбранного текущего города, и так же в прогнозе
//TODO Няшные менюшки и вообще навигация
//TODO start using docopt

import (
	"database/sql"
	"flag"
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"github.com/jasonlvhit/gocron"
	"golang.org/x/net/context"
	"log"
)

var db *sql.DB
var err error

func main() {
	var telegramKey string
	var debugApi bool

	flag.StringVar(&telegramKey, "k", "", "sekret telegram api key")
	flag.BoolVar(&debugApi, "d", false, "enable api debug mode")
	flag.Parse()
	if telegramKey == "" {
		panic("Secret telegram key not setted")
	}

	db, err = initAppDb()
	if err != nil {
		panic(err)
	}
	api := telegram.New(telegramKey)
	api.Debug(debugApi)
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

	gocron.Every(5).Seconds().Do(runNotificationTasks, api, netCtx)
	<-gocron.Start()

	err = bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}
