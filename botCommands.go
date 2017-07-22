package main

import (
	"log"
	"strings"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
)

func startCommand(ctx context.Context, arg string) error {
	api := telebot.GetAPI(ctx)
	update := telebot.GetUpdate(ctx)
	_, err := api.SendMessage(ctx,
		telegram.NewMessagef(update.Chat().ID,
			"received start with arg %s", arg,
		))

	return err
}

func currentCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	textMessage := "No selected city. Select city with command \n/city {cityName}"
	log.Println("current1")
	userCity, err := ds.getUserCity(update.From().ID)
	if err != nil || userCity == nil {
		return sendMessage(ctx, update.Chat().ID, textMessage)
	}
	log.Println("current")
	if currentWeather, err := ws.getCurrentWeather(userCity.CityAlias); err == nil {
		textMessage = ws.formatCurrentWeather(currentWeather)
	}
	return sendMessage(ctx, update.Chat().ID, textMessage)
}

func sendMessage(ctx context.Context, userID int64, textMessage string) error {
	api := telebot.GetAPI(ctx) // take api from context
	msg := telegram.NewMessage(userID, textMessage)
	_, err := api.Send(ctx, msg)
	return err
}

func cityCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	user := update.From()
	log.Println(arg)
	cities, err := ws.getCities(arg)
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, err.Error())
	}

	if update.CallbackQuery != nil {
		log.Println("cityCommand")
		log.Println(update.CallbackQuery)
	}

	if len(cities) > 1 {
		msg := telegram.NewMessage(update.Chat().ID, "Please, set city:")
		var keyboardText = [][]telegram.InlineKeyboardButton{}
		keyboardRow := []telegram.InlineKeyboardButton{}
		for index, city := range cities {
			log.Println(index)
			log.Println(city.Title)
			keyboardRow = append(
				keyboardRow,
				telegram.InlineKeyboardButton{
					Text:         city.Title,
					CallbackData: "/city " + city.Alias,
				},
			)

			if (index+1)%3 == 0 {
				keyboardText = append(keyboardText, keyboardRow)
				keyboardRow = []telegram.InlineKeyboardButton{}
			}
		}

		msg.ReplyMarkup = telegram.InlineKeyboardMarkup{
			InlineKeyboard: keyboardText,
		}

		api := telebot.GetAPI(ctx) // take api from context
		_, err := api.SendMessage(ctx, msg)
		return err
	}

	city := cities[0]
	textMessage := "City selected: " + city.Title

	if err := ds.saveUserCity(user.ID, city.Alias); err != nil {
		textMessage = "Cant set selected city " + city.Title + ". Try later."
	}
	return sendMessage(ctx, update.Chat().ID, textMessage)
}

func defaultCommand(ctx context.Context) error {
	update := telebot.GetUpdate(ctx) // take update from context

	if update.CallbackQuery != nil {
		data := update.CallbackQuery.Data
		if strings.HasPrefix(data, "/city ") {
			city := strings.Replace(data, "/city ", "", -1)
			cityCommand(ctx, city)
		}
	}

	if update.Message == nil {
		return nil
	}

	city, err := ws.getCity(update.Message.Text)
	if city == nil {
		return sendMessage(ctx, update.Chat().ID, err.Error())
	}

	textMessage := "Cant get current weather. Try later"
	if currentWeather, err := ws.getCurrentWeather(city.Alias); err == nil {
		textMessage = ws.formatCurrentWeather(currentWeather)
	}
	return sendMessage(ctx, update.Chat().ID, textMessage)
}

func helpCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	return sendMessage(ctx, update.Chat().ID, "It's Ngs weather bot. \nIt can show you weather in city. To start messaging, you can send city name in your message.\nCommands:\n/city {cityName} - set your prefered city to show by /curent command\n/current - show's the weather in prefered city by /city command\n/forecast - show's forecast by prefered city/help")
}
