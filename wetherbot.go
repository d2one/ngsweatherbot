package main

//TODO кэширование ответов апи новостей, относительно прогноза на 10 минут
//TODO сделать нотификации пользователям, о выбранной погоде, по времени
//TODO Сделать вывод выбранного текущего города, и так же в прогнозе
//TODO Няшные менюшки и вообще навигация

import (
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
	"log"
	"flag"
)

func main() {
	var telegramKey string
	var debugApi  bool
	flag.StringVar(&telegramKey, "k", "", "sekret telegram api key")
	flag.BoolVar(&debugApi, "d", false, "enable api debug mode")
	flag.Parse()
	if telegramKey == "" {
		panic("Secret telegram key not setted")
	}

	db, err := initAppDb()
	if err != nil {
		panic(err)
	}
	api := telegram.New(telegramKey)
	api.Debug(debugApi)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(func(ctx context.Context) error {
		update := telebot.GetUpdate(ctx) // take update from context
		if update.Message == nil {
			return nil
		}
		var textMessage string
		city, err := getCity(update.Message.Text)
		if err != nil {
			_, err := api.SendMessage(ctx,
				telegram.NewMessagef(update.Chat().ID, err.Error()))
			return err
		}

		currentWeather, err := getCurrentWeather(city.Alias)
		if err != nil {
			textMessage = "Cant get current weather. Try later"
		} else {textMessage = formatCurrentWeather(currentWeather)}
		api := telebot.GetAPI(ctx) // take api from context
		msg := telegram.NewMessage(update.Chat().ID, textMessage)
		if _, err2 := api.Send(ctx, msg); err2 != nil {
			return err2
		}
		return nil
	})

	// Use command middleware, that helps to work with commands
	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {
				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				_, err := api.SendMessage(ctx,
					telegram.NewMessagef(update.Chat().ID,
						"received start with arg %s", arg,
					))

				return err
			}),
		"current": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {
				update := telebot.GetUpdate(ctx)
				textMessage := "No selected city. Select city with command \n/city {cityName}"
				userCity, _ := getUserCity(db, update.From().ID)
				if userCity != (UserCity{}) {
					if currentWeather, err := getCurrentWeather(userCity.city_alias); err == nil {
						textMessage = formatCurrentWeather(currentWeather)
					}
				}
				api := telebot.GetAPI(ctx) // take api from context
				msg := telegram.NewMessage(update.Chat().ID, textMessage)
				_, err2 := api.Send(ctx, msg)
				return err2
			}),
		"help": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {

				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				_, err := api.SendMessage(ctx,
					telegram.NewMessagef(update.Chat().ID,
						"It's Ngs weather bot. \nIt can show you weather in city. To start messaging, you can send city name in your message.\nCommands:\n/city {cityName} - set your prefered city to show by /curent command\n/current - show's the weather in prefered city by /city command\n/forecast - show's forecast by prefered city/help",
					))
				return err
			}),
		"city": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {
				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				user := update.From()

				cities, err := getCities(arg)
				if err != nil {
					_, err := api.SendMessage(ctx,
						telegram.NewMessagef(update.Chat().ID, err.Error()))
					return err
				}

				if len(cities) > 1 {
					textMessage := "Please, set city:\n"
					for index := range cities {
						textMessage += "/city " + cities[index].Alias + "\n"
					}
					api.SendMessage(ctx,
						telegram.NewMessagef(update.Chat().ID,
							textMessage,
						))
					return nil
				}

				city := cities[0]
				userCity := UserCity{
					user_id: user.ID,
					city_alias: city.Alias,
					chat_id: update.Chat().ID,
				}

				textMessage := "City selected: " + city.Title
				if err := saveUserCity(db, userCity); err != nil {
					textMessage = "Cant set selected city " + city.Title + ". Try later."
				}

				api.SendMessage(ctx, telegram.NewMessagef(userCity.chat_id, textMessage))
				return nil
			}),
	}))

	err = bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}
