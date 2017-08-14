package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"

	"strconv"

	"encoding/json"
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
)

func startCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	return sendMessage(ctx, update.Chat().ID, "Что будет делать??", getDefaultKeyboard())
}

func initCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	dataStore.initUser(update.Chat().ID)
	messageText := fmt.Sprintf(
		"Привет *%s %s*. Я погодный бот. Ты сможешь получать от меня данные о погоде от pogoda.ngs.ru",
		update.Chat().FirstName,
		update.Chat().LastName,
	)
	return sendMessage(ctx, update.Chat().ID, messageText, getDefaultKeyboard())
}

func settingsCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	api := telebot.GetAPI(ctx) // take api from context
	textMessage := fmt.Sprintf(
		"*Настройки:*\n%s - выбрать свой город\n%s - настроить уведомления о погоде",
		weatherIcons.getButtons("default_city"),
		weatherIcons.getButtons("notifications"),
	)
	msg := telegram.NewMessage(update.Chat().ID, textMessage)
	msg.ReplyMarkup = getSettingKeyboard()
	msg.ParseMode = "markdown"
	_, err := api.Send(ctx, msg)
	return err
}

func currentCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	userData, err := dataStore.getUserData(update.From().ID)
	if err != nil || !userData.CityAlias.Valid {
		return sendMessage(ctx, update.Chat().ID, "Город не выбран. Выберете город в настройках.", nil)
	}
	city := &City{
		Alias: userData.CityAlias.String,
		Title: userData.CityTitle.String,
	}
	textMessage, replyMarkup := currentWeather(city, userData.ForecastType)
	return sendMessage(ctx, update.Chat().ID, textMessage, replyMarkup)
}

func commandForecast(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	textMessage := "Город не выбран. Выберете город в настройках."
	userData, err := dataStore.getUserData(update.From().ID)
	if err != nil || !userData.CityAlias.Valid {
		return sendMessage(ctx, update.Chat().ID, textMessage, nil)
	}
	if forecast, err := weatherService.getForecast(userData.CityAlias.String); err == nil {
		textMessage = userData.CityTitle.String + "\n" + weatherService.formatForecastWeather(forecast)
	}
	return sendMessage(ctx, update.Chat().ID, textMessage, nil)
}

func commandSetNotifications(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	t := time.Now()
	s := strings.Split(arg, ":")
	errorText := "Неверный формат времени"
	if len(s) != 2 {
		return sendMessage(ctx, update.Chat().ID, errorText, nil)
	}

	hours, err := strconv.Atoi(s[0])
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, errorText, nil)
	}
	minutes, err := strconv.Atoi(s[1])
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, errorText, nil)
	}
	notificationTime := time.Date(t.Year(), t.Month(), t.Day(), hours, minutes, int(0), int(0), time.Local)

	if t.Unix() > notificationTime.Unix() {
		notificationTime = notificationTime.AddDate(0, 0, 1)
	}

	if err = dataStore.saveUserNotification(update.Chat().ID, notificationTime.Unix()); err != nil {
		logWork(err)
	}
	cache.delete(strconv.Itoa(int(update.Chat().ID)))
	return sendMessage(ctx, update.Chat().ID, "Установлено: "+notificationTime.Format("15:04"), getDefaultKeyboard())
}

func commandNotifications(ctx context.Context, userID int64) error {
	userData, _ := dataStore.getUserData(int64(userID))
	textMessage := "Уведомления выключены"
	if userData.NotificationsNextRun.Valid {
		t := time.Unix(userData.NotificationsNextRun.Int64, 0)
		textMessage = "У вас включены уведомления: *" + t.Format("15:04") + "*"
	}
	keyboard := buildKeyboard([]string{
		weatherIcons.getButtons("notifications_set"),
		weatherIcons.getButtons("notifications_remove"),
		weatherIcons.getButtons("back"),
	}, 3, false)
	textMessage += fmt.Sprintf("\n%s - установить время уведомления\n%s - отключить уведомления",
		weatherIcons.getButtons("notifications_set"),
		weatherIcons.getButtons("notifications_remove"),
	)
	return sendMessage(ctx, userID, textMessage, keyboard)
}

func cityCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	user := update.From()
	cities, err := weatherService.getCities(arg)
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, err.Error(), getDefaultKeyboard())
	}

	if len(cities) > 1 {
		return buildInlineCityKeyboard(ctx, cities, "/city:")
	}

	city := cities[0]
	textMessage := "Город выбран: " + city.Title

	if err := dataStore.saveUserCity(user.ID, city); err != nil {
		textMessage = "Не удаолось выбрать город " + city.Title + ". Попробуйте позже."
	}
	return sendMessage(ctx, update.Chat().ID, textMessage, getDefaultKeyboard())
}

func helpCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	return sendMessage(ctx, update.Chat().ID, "It's Ngs weather bot. \nIt can show you weather in city. To start messaging, you can send city name in your message.\nCommands:\n/city {cityName} - set your prefered city to show by /curent command\n/current - show's the weather in prefered city by /city command\n/forecast - show's forecast by prefered city/help", nil)
}

func defaultCommand(ctx context.Context) error {
	update := telebot.GetUpdate(ctx) // take update from context
	if update.Message == nil {
		return nil
	}

	messageText := update.Message.Text

	if err, isCommandRun := inProcessCommand(ctx, messageText); isCommandRun {
		return err
	}
	cities, err := weatherService.getCities(messageText)
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, err.Error(), getDefaultKeyboard())
	}

	if len(cities) > 1 {
		return buildInlineCityKeyboard(ctx, cities, "/show_city:")
	}
	city := cities[0]
	userData, err := dataStore.getUserData(update.From().ID)
	textMessage, replyMarkup := currentWeather(city, userData.ForecastType)
	cityJson, _ := json.Marshal(city)
	cache.write("last_city"+strconv.Itoa(int(update.Chat().ID)), string(cityJson), time.Now().Unix()+60*10)
	return sendMessage(ctx, update.Chat().ID, textMessage, replyMarkup)
}

func buildInlineCityKeyboard(ctx context.Context, cities []*City, callback string) error {
	update := telebot.GetUpdate(ctx)
	msg := telegram.NewMessage(update.Chat().ID, "Пожалуйста, выберете город:")
	var keyboardText = [][]telegram.InlineKeyboardButton{}
	keyboardRow := []telegram.InlineKeyboardButton{}
	for index, city := range cities {
		keyboardRow = append(
			keyboardRow,
			telegram.InlineKeyboardButton{
				Text:         city.Title,
				CallbackData: callback + city.Title,
			},
		)
		log.Println(callback + city.Title)
		if (index+1)%3 == 0 {
			keyboardText = append(keyboardText, keyboardRow)
			keyboardRow = []telegram.InlineKeyboardButton{}
		}
	}
	keyboardText = append(keyboardText, keyboardRow)

	msg.ReplyMarkup = telegram.InlineKeyboardMarkup{
		InlineKeyboard: keyboardText,
	}

	api := telebot.GetAPI(ctx) // take api from context
	_, err := api.SendMessage(ctx, msg)
	return err
}

func callbackWeatherType(ctx context.Context, forecastType string) error {
	update := telebot.GetUpdate(ctx) // take update from context
	dataStore.saveForecastType(update.Chat().ID, forecastType)
	if cityData, _ := cache.read("last_city" + strconv.Itoa(int(update.Chat().ID))); cityData != nil {
		var city = &City{}
		json.Unmarshal([]byte(cityData.CacheValue), &city)
		textMessage, replyMarkup := currentWeather(city, forecastType)
		return sendMessage(ctx, update.Chat().ID, textMessage, replyMarkup)
	}

	return currentCommand(ctx, "")
}

func sendMessage(ctx context.Context, userID int64, textMessage string, markup telegram.ReplyMarkup) error {
	api := telebot.GetAPI(ctx) // take api from context
	msg := telegram.NewMessage(userID, textMessage)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	_, err := api.Send(ctx, msg)
	return err
}

func currentWeather(city *City, forecastType string) (string, telegram.InlineKeyboardMarkup) {
	var replyMarkup telegram.InlineKeyboardMarkup
	var textMessage string
	if currentWeather, err := weatherService.getCurrentWeather(city.Alias); err == nil {
		log.Println(currentWeather)
		if forecastType == "full" {
			textMessage, replyMarkup = weatherService.formatFullCurrentWeather(currentWeather)
		} else {
			textMessage, replyMarkup = weatherService.formatCurrentWeather(currentWeather)
		}
	}
	textMessage = "*" + city.Title + "*\n" + textMessage
	return textMessage, replyMarkup
}

func getDefaultKeyboard() telegram.ReplyMarkup {
	return buildKeyboard([]string{
		weatherIcons.getButtons("current"),
		weatherIcons.getButtons("forecast"),
		weatherIcons.getButtons("settings"),
	}, 3, false)
}

func getSettingKeyboard() telegram.ReplyMarkup {
	return buildKeyboard([]string{
		weatherIcons.getButtons("default_city"),
		weatherIcons.getButtons("notifications"),
		weatherIcons.getButtons("back"),
	}, 3, false)
}

func buildKeyboard(menu []string, rows int, oneTimeKeyboard bool) telegram.ReplyKeyboardMarkup {
	var keyboard = [][]string{}
	var keyboardRow = []string{}

	for index, menuItem := range menu {
		keyboardRow = append(keyboardRow, menuItem)
		if (index+1)%rows == 0 {
			keyboard = append(keyboard, keyboardRow)
			keyboardRow = []string{}
		}
	}
	return telegram.ReplyKeyboardMarkup{
		Keyboard:        telegram.NewKeyboard(keyboard),
		OneTimeKeyboard: oneTimeKeyboard,
		ResizeKeyboard:  true,
	}
}

func inProcessCommand(ctx context.Context, commandText string) (error, bool) {
	update := telebot.GetUpdate(ctx) // take update from context
	switch commandText {
	case weatherIcons.getButtons("current"):
		return currentCommand(ctx, ""), true
	case weatherIcons.getButtons("settings"):
		return settingsCommand(ctx, ""), true
	case weatherIcons.getButtons("default_city"):
		cache.write(strconv.Itoa(int(update.Chat().ID)), "default_city", time.Now().Unix()+60*10)
		return sendMessage(ctx, update.Chat().ID, "Введите город по умолчанию:", nil), true
	case weatherIcons.getButtons("forecast"):
		return commandForecast(ctx, ""), true
	case weatherIcons.getButtons("notifications_set"):
		cache.write(strconv.Itoa(int(update.Chat().ID)), "notification_user", time.Now().Unix()+60*10)
		return sendMessage(ctx, update.Chat().ID, "Введите время", nil), true
	case weatherIcons.getButtons("notifications_remove"):
		dataStore.deleteUserNotification(update.Chat().ID)
		return sendMessage(ctx, update.Chat().ID, "Уведомления отключены", getDefaultKeyboard()), true
	case weatherIcons.getButtons("notifications"):
		return commandNotifications(ctx, update.Chat().ID), true
	case weatherIcons.getButtons("back"):
		cache.delete(strconv.Itoa(int(update.Chat().ID)))
		return startCommand(ctx, ""), true
	}

	if lastCommand, _ := cache.read(strconv.Itoa(int(update.Chat().ID))); lastCommand != nil {
		switch lastCommand.CacheValue {
		case "default_city":
			cache.delete(strconv.Itoa(int(update.Chat().ID)))
			return cityCommand(ctx, update.Message.Text), true
		case "notification_user":
			return commandSetNotifications(ctx, update.Message.Text), true
		}
	}
	return nil, false
}
