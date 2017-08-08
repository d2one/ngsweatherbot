package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"

	"strconv"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
)

func startCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	return sendMessage(ctx, update.Chat().ID, "Что будет делать??", getDefaultKeyboard())
}

func initCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	ds.initUser(update.Chat().ID)
	messageText := fmt.Sprintf(
		"Привет *%s %s*. Я погодный бот. Ты сможешь получать от меня данные о погоде от pogoda.ngs.ru",
		update.Chat().FirstName,
		update.Chat().LastName,
	)
	return sendMessage(ctx, update.Chat().ID, messageText, getDefaultKeyboard())
}

func getDefaultKeyboard() telegram.ReplyMarkup {
	return buildKeybopard([]string{
		wi.getButtons("current"),
		wi.getButtons("forecast"),
		wi.getButtons("settings"),
	}, 3, false)
}

func getSettingKeyboard() telegram.ReplyMarkup {
	return buildKeybopard([]string{
		wi.getButtons("default_city"),
		wi.getButtons("notifications"),
		wi.getButtons("back"),
	}, 3, false)
}

func buildKeybopard(menu []string, rows int, oneTimeKeyboard bool) telegram.ReplyKeyboardMarkup {
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

func settingsCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	api := telebot.GetAPI(ctx) // take api from context
	textMessage := fmt.Sprintf(
		"*Настройки:*\n%s - выбрать свой город\n%s - настроить уведомления о погоде",
		wi.getButtons("default_city"),
		wi.getButtons("notifications"),
	)
	msg := telegram.NewMessage(update.Chat().ID, textMessage)
	msg.ReplyMarkup = getSettingKeyboard()
	msg.ParseMode = "markdown"
	_, err := api.Send(ctx, msg)
	return err
}

func currentCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	textMessage := "Город не выбран. Выберете город в настройках."
	userData, err := ds.getUserData(update.From().ID)
	if err != nil || !userData.CityAlias.Valid {
		return sendMessage(ctx, update.Chat().ID, textMessage, nil)
	}

	var replyMarkup = telegram.InlineKeyboardMarkup{}
	if currentWeather, err := ws.getCurrentWeather(userData.CityAlias.String); err == nil {
		if userData.ForecastType == "full" {
			textMessage, replyMarkup = ws.formatFullCurrentWeather(currentWeather)
		} else {
			textMessage, replyMarkup = ws.formatCurrentWeather(currentWeather)
		}

		textMessage = "*" + userData.CityTitle.String + "*\n" + textMessage
	}

	return sendMessage(ctx, update.Chat().ID, textMessage, replyMarkup)
}

func commandForecast(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	textMessage := "Город не выбран. Выберете город в настройках."
	userData, err := ds.getUserData(update.From().ID)
	if err != nil || !userData.CityAlias.Valid {
		return sendMessage(ctx, update.Chat().ID, textMessage, nil)
	}
	if forecast, err := ws.getForecast(userData.CityAlias.String); err == nil {
		textMessage = userData.CityTitle.String + "\n" + ws.formatForecasttWeather(forecast)
	}
	return sendMessage(ctx, update.Chat().ID, textMessage, nil)
}

func commandSetNotifications(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	t := time.Now()
	s := strings.Split(arg, ":")

	if len(s) != 2 {
		return sendMessage(ctx, update.Chat().ID, "Неверный формат времени", nil)
	}

	hours, err := strconv.Atoi(s[0])
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, "Неверный формат времени", nil)
	}
	minutes, err := strconv.Atoi(s[1])
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, "Неверный формат времени", nil)
	}
	time := time.Date(t.Year(), t.Month(), t.Day(), hours, minutes, int(0), int(0), time.Local)

	if t.Unix() > time.Unix() {
		time = time.AddDate(0, 0, 1)
	}

	if err = ds.saveUserNotification(update.Chat().ID, time.Unix()); err != nil {
		logWork(err)
	}
	cacheService.delete(strconv.Itoa(int(update.Chat().ID)))
	return sendMessage(ctx, update.Chat().ID, "Установлено: "+time.Format("15:04"), getDefaultKeyboard())
}

func commandNotifications(ctx context.Context, userID int64) error {
	userData, _ := ds.getUserData(int64(userID))
	textMessage := "Уведомления выключены"
	if userData.NotificationsNextRun.Valid {
		t := time.Unix(userData.NotificationsNextRun.Int64, 0)
		textMessage = "У вас включены уведомления: *" + t.Format("15:04") + "*"
	}
	keyboard := buildKeybopard([]string{
		wi.getButtons("notifications_set"),
		wi.getButtons("notifications_remove"),
		wi.getButtons("back"),
	}, 3, false)
	textMessage += fmt.Sprintf("\n%s - установить время уведомления\n%s - отключить уведомления",
		wi.getButtons("notifications_set"),
		wi.getButtons("notifications_remove"),
	)
	return sendMessage(ctx, userID, textMessage, keyboard)
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

func cityCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	user := update.From()
	cities, err := ws.getCities(arg)
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, err.Error(), getDefaultKeyboard())
	}

	if len(cities) > 1 {
		return buildInlineCityKeyboard(ctx, cities, "/city ")
	}

	city := cities[0]
	textMessage := "Город выбран: " + city.Title

	if err := ds.saveUserCity(user.ID, city); err != nil {
		textMessage = "Не удаолось выбрать город " + city.Title + ". Попробуйте позже."
	}
	return sendMessage(ctx, update.Chat().ID, textMessage, getDefaultKeyboard())
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

func runMenuCommand(ctx context.Context, commandText string) (error, bool) {
	update := telebot.GetUpdate(ctx) // take update from context
	switch commandText {
	case wi.getButtons("current"):
		return currentCommand(ctx, ""), true
	case wi.getButtons("settings"):
		return settingsCommand(ctx, ""), true
	case wi.getButtons("default_city"):
		cacheService.write(strconv.Itoa(int(update.Chat().ID)), "default_city", time.Now().Unix()+60*10)
		return sendMessage(ctx, update.Chat().ID, "Введите город по умолчанию:", nil), true
	case wi.getButtons("forecast"):
		return commandForecast(ctx, ""), true
	case wi.getButtons("notifications_set"):
		cacheService.write(strconv.Itoa(int(update.Chat().ID)), "notification_user", time.Now().Unix()+60*10)
		return sendMessage(ctx, update.Chat().ID, "Введите время", nil), true
	case wi.getButtons("notifications_remove"):
		ds.deleteUserNotification(update.Chat().ID)
		return sendMessage(ctx, update.Chat().ID, "Уведомления отключены", getDefaultKeyboard()), true
	case wi.getButtons("notifications"):
		return commandNotifications(ctx, update.Chat().ID), true
	case wi.getButtons("back"):
		cacheService.delete(strconv.Itoa(int(update.Chat().ID)))
		return startCommand(ctx, ""), true
	}

	if lastCommand, _ := cacheService.read(strconv.Itoa(int(update.Chat().ID))); lastCommand != nil {
		switch lastCommand.CacheValue {
		case "default_city":
			cacheService.delete(strconv.Itoa(int(update.Chat().ID)))
			return cityCommand(ctx, update.Message.Text), true
		case "notification_user":
			return commandSetNotifications(ctx, update.Message.Text), true
		}
	}
	return nil, false
}

func defaultCommand(ctx context.Context) error {
	update := telebot.GetUpdate(ctx) // take update from context
	var messageText string
	if update.CallbackQuery != nil {
		data := update.CallbackQuery.Data
		if strings.HasPrefix(data, "/city ") {
			city := strings.Replace(data, "/city ", "", -1)
			return cityCommand(ctx, city)
		}
		if strings.HasPrefix(data, "/current ") {
			forecastType := strings.Replace(data, "/current ", "", -1)
			ds.saveForecastType(update.Chat().ID, forecastType)
			return currentCommand(ctx, "")
		}

		if strings.HasPrefix(data, "/forecast ") {
			forecastType := strings.Replace(data, "/forecast ", "", -1)
			ds.saveForecastType(update.Chat().ID, forecastType)
			return currentCommand(ctx, "")
		}

		messageText = update.CallbackQuery.Data
	}

	if update.Message == nil && messageText == "" {
		return nil
	}
	if messageText == "" {
		messageText = update.Message.Text
	}

	if err, isCommandRun := runMenuCommand(ctx, messageText); isCommandRun {
		return err
	}

	//TODO выделить эту штуку в отдельную ф-ю

	cities, err := ws.getCities(messageText)
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, err.Error(), getDefaultKeyboard())
	}

	if len(cities) > 1 {
		return buildInlineCityKeyboard(ctx, cities, "")
	}
	city := cities[0]
	log.Println(city.Alias)
	textMessage := "Не удалось получить текущую погоду. Попробуйте позже"
	var replyMarkup = telegram.InlineKeyboardMarkup{}
	if currentWeather, err := ws.getCurrentWeather(city.Alias); err == nil {
		textMessage, replyMarkup = ws.formatCurrentWeather(currentWeather)
		textMessage = "*" + city.Title + "*\n " + textMessage
	}
	return sendMessage(ctx, update.Chat().ID, textMessage, replyMarkup)
}

func helpCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	return sendMessage(ctx, update.Chat().ID, "It's Ngs weather bot. \nIt can show you weather in city. To start messaging, you can send city name in your message.\nCommands:\n/city {cityName} - set your prefered city to show by /curent command\n/current - show's the weather in prefered city by /city command\n/forecast - show's forecast by prefered city/help", nil)
}
