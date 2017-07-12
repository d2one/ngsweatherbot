package main

import (
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
	"log"
	"flag"
)

func main() {
	var telegramKey string
	flag.StringVar(&telegramKey, "k", "", "sekret telegram api key")
	flag.Parse()
	if telegramKey == "" {
		panic("Secret telegram key not setted")
	}


	db :=InitDB("db.sqlite3")
	CreateTable(db)
	api := telegram.New(telegramKey)
	api.Debug(true)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(func(ctx context.Context) error {
		update := telebot.GetUpdate(ctx) // take update from context
		if update.Message == nil {
			return nil
		}
		city, error := getCity(update.Message.Text)

		if error != "" {
			_, err := api.SendMessage(ctx,
				telegram.NewMessagef(update.Chat().ID, error))
			return err
		}

		textMessage := getStations(city)
		api := telebot.GetAPI(ctx) // take api from context
		msg := telegram.NewMessage(update.Chat().ID, textMessage)
		_, err := api.Send(ctx, msg)
		return err

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


				userData := ReadItem(db, update.From().ID)
				if userData.city_alias != "" {
					textMessage = getStations(userData.city_alias)
				}
				api := telebot.GetAPI(ctx) // take api from context
				msg := telegram.NewMessage(update.Chat().ID, textMessage)
				_, err := api.Send(ctx, msg)
				return err
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

				city, error := getCities(arg)

				if error != "" {
					_, err := api.SendMessage(ctx,
						telegram.NewMessagef(update.Chat().ID, error))
					return err
				}

				if len(city.Cities) > 1 {
					textMessage := "Please, set city:\n"
					for index := range city.Cities {
						log.Println(index)
						log.Println(city.Cities[index].Alias)
						textMessage += "/city " + city.Cities[index].Alias + "\n"
					}
					api.SendMessage(ctx,
						telegram.NewMessagef(update.Chat().ID,
							textMessage,
						))
					return nil
				}

				userData := UserCity{
					user_id: user.ID,
					city_alias: city.Cities[0].Alias,
					chat_id: update.Chat().ID,
				}

				StoreItem(db, []UserCity{userData})
				api.SendMessage(ctx,
					telegram.NewMessagef(update.Chat().ID,
						"City selected: %s", city.Cities[0].Title,
					))
				return nil
			}),
	}))

	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}

