package main

import (
	"log"

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

	userCity, err := ds.getUserCity(update.From().ID)
	if err != nil || userCity == nil {
		return sendMessage(ctx, update.Chat().ID, textMessage)
	}

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

	if len(cities) > 1 {
		textMessage := "Please, set city:\n"
		for index := range cities {
			textMessage += "/city " + cities[index].Alias + "\n"
		}
		return sendMessage(ctx, update.Chat().ID, textMessage)
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
	if update.Message == nil {
		return nil
	}
	var textMessage string
	city, err := ws.getCity(update.Message.Text)
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, err.Error())
	}

	textMessage = "Cant get current weather. Try later"
	currentWeather, err := ws.getCurrentWeather(city.Alias)
	if err == nil {
		textMessage = ws.formatCurrentWeather(currentWeather)
	}

	return sendMessage(ctx, update.Chat().ID, textMessage)
}

func helpCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	return sendMessage(ctx, update.Chat().ID, "It's Ngs weather bot. \nIt can show you weather in city. To start messaging, you can send city name in your message.\nCommands:\n/city {cityName} - set your prefered city to show by /curent command\n/current - show's the weather in prefered city by /city command\n/forecast - show's forecast by prefered city/help")
}
