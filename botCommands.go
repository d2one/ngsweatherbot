package main

import (
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
	userCity, _ := getUserCity(update.From().ID)
	if userCity != (UserCity{}) {
		if currentWeather, err := getCurrentWeather(userCity.city_alias); err == nil {
			textMessage = formatCurrentWeather(currentWeather)
		}
	}
	api := telebot.GetAPI(ctx) // take api from context
	msg := telegram.NewMessage(update.Chat().ID, textMessage)
	_, err2 := api.Send(ctx, msg)
	return err2
}

func cityCommand(ctx context.Context, arg string) error {
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
		user_id:    user.ID,
		city_alias: city.Alias,
		chat_id:    update.Chat().ID,
	}

	textMessage := "City selected: " + city.Title
	if err := saveUserCity(userCity); err != nil {
		textMessage = "Cant set selected city " + city.Title + ". Try later."
	}
	api.SendMessage(ctx, telegram.NewMessagef(userCity.chat_id, textMessage))
	return nil
}

func defaultCommand(ctx context.Context) error {
	update := telebot.GetUpdate(ctx) // take update from context
	if update.Message == nil {
		return nil
	}
	var textMessage string
	city, err := getCity(update.Message.Text)
	api := telebot.GetAPI(ctx) // take api from context
	if err != nil {
		_, err := api.SendMessage(ctx,
			telegram.NewMessagef(update.Chat().ID, err.Error()))
		return err
	}

	textMessage = "Cant get current weather. Try later"
	currentWeather, err := getCurrentWeather(city.Alias)
	if err == nil {
		textMessage = formatCurrentWeather(currentWeather)
	}

	msg := telegram.NewMessage(update.Chat().ID, textMessage)
	if _, err2 := api.Send(ctx, msg); err2 != nil {
		return err2
	}
	return nil
}

func helpCommand(ctx context.Context, arg string) error {
	api := telebot.GetAPI(ctx)
	update := telebot.GetUpdate(ctx)
	_, err := api.SendMessage(ctx,
		telegram.NewMessagef(update.Chat().ID,
			"It's Ngs weather bot. \nIt can show you weather in city. To start messaging, you can send city name in your message.\nCommands:\n/city {cityName} - set your prefered city to show by /curent command\n/current - show's the weather in prefered city by /city command\n/forecast - show's forecast by prefered city/help",
		))
	return err
}
