package main

import (
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

func getDefaultKeyboard() telegram.ReplyMarkup {
	return buildKeybopard([]string{
		wi.getButtons("current"),
		wi.getButtons("forecast"),
		wi.getButtons("settings"),
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
	msg := telegram.NewMessage(update.Chat().ID, "Настройки:")
	msg.ReplyMarkup = buildKeybopard([]string{
		wi.getButtons("default_city"),
		wi.getButtons("notifications"),
		wi.getButtons("back"),
	}, 3, false)
	_, err := api.Send(ctx, msg)
	return err
}

func currentCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	textMessage := "Город не выбран. Выберете город в настройках."
	userCity, err := ds.getUserCity(update.From().ID)
	if err != nil || userCity == nil {
		return sendMessage(ctx, update.Chat().ID, textMessage, nil)
	}
	log.Println("current")
	if currentWeather, err := ws.getCurrentWeather(userCity.CityAlias); err == nil {
		textMessage = "*" + userCity.CityTitle + "*\n " + ws.formatCurrentWeather(currentWeather)
	}
	return sendMessage(ctx, update.Chat().ID, textMessage, nil)
}

func commandForecast(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	textMessage := "Город не выбран. Выберете город в настройках."
	userCity, err := ds.getUserCity(update.From().ID)
	if err != nil || userCity == nil {
		return sendMessage(ctx, update.Chat().ID, textMessage, nil)
	}
	if forecast, err := ws.getForecast(userCity.CityAlias); err == nil {
		textMessage = userCity.CityTitle + "\n" + ws.formatForecasttWeather(forecast)
	}
	return sendMessage(ctx, update.Chat().ID, textMessage, nil)
}

func commandNotifications(ctx context.Context, arg string) error {
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
	un := &UserNotification{
		UserID:  update.Chat().ID,
		ChatID:  update.Chat().ID,
		NextRun: time.Unix(),
	}
	if err = ds.saveUserNotification(*un); err != nil {
		logWork(err)
	}
	cache.delete(strconv.Itoa(int(update.Chat().ID)))
	return sendMessage(ctx, update.Chat().ID, "Установлено: "+time.Format("15:04"), getDefaultKeyboard())
}

func sendMessage(ctx context.Context, userID int64, textMessage string, markup telegram.ReplyMarkup) error {
	api := telebot.GetAPI(ctx) // take api from context
	msg := telegram.NewMessage(userID, textMessage)
	msg.ParseMode = "markdown"
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	_, err := api.Send(ctx, msg)
	return err
}

func cityCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	user := update.From()
	log.Println(arg)
	cities, err := ws.getCities(arg)
	if err != nil {
		return sendMessage(ctx, update.Chat().ID, err.Error(), getDefaultKeyboard())
	}

	if len(cities) > 1 {
		msg := telegram.NewMessage(update.Chat().ID, "Пожалуйста, выберете город:")
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
			// TODO проверку на количество
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
	textMessage := "Город выбран: " + city.Title

	if err := ds.saveUserCity(user.ID, city); err != nil {
		textMessage = "Не удаолось выбрать город " + city.Title + ". Попробуйте позже."
	}
	return sendMessage(ctx, update.Chat().ID, textMessage, getDefaultKeyboard())
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

	//TODO выделить эту штуку в отдельную ф-ю
	switch update.Message.Text {
	case wi.getButtons("current"):
		return currentCommand(ctx, "")
	case wi.getButtons("settings"):
		return settingsCommand(ctx, "")
	case wi.getButtons("default_city"):
		cache.write(strconv.Itoa(int(update.Chat().ID)), "default_city", time.Now().Unix()+60*10)
		return sendMessage(ctx, update.Chat().ID, "Введите город по умолчанию:", nil)
	case wi.getButtons("forecast"):
		return commandForecast(ctx, "")
	case wi.getButtons("notifications"):
		cache.write(strconv.Itoa(int(update.Chat().ID)), "notification_user", time.Now().Unix()+60*10)
		return sendMessage(ctx, update.Chat().ID, "Введите время для нотификаций:", nil)
	case wi.getButtons("back"):
		cache.delete(strconv.Itoa(int(update.Chat().ID)))
		return startCommand(ctx, "")
	}

	if lastCommand, _ := cache.read(strconv.Itoa(int(update.Chat().ID))); lastCommand != nil {
		switch lastCommand.CacheValue {
		case "default_city":
			cache.delete(strconv.Itoa(int(update.Chat().ID)))
			return cityCommand(ctx, update.Message.Text)
		case "notification_user":
			return commandNotifications(ctx, update.Message.Text)
		}
	}

	city, err := ws.getCity(update.Message.Text)
	if city == nil {
		return sendMessage(ctx, update.Chat().ID, err.Error(), nil)
	}

	textMessage := "Не удалось получить текущую погоду. Попробуйте позже"
	if currentWeather, err := ws.getCurrentWeather(city.Alias); err == nil {
		textMessage = "*" + city.Title + "*\n " + ws.formatCurrentWeather(currentWeather)
	}

	return sendMessage(ctx, update.Chat().ID, textMessage, getDefaultKeyboard())
}

func helpCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	return sendMessage(ctx, update.Chat().ID, "It's Ngs weather bot. \nIt can show you weather in city. To start messaging, you can send city name in your message.\nCommands:\n/city {cityName} - set your prefered city to show by /curent command\n/current - show's the weather in prefered city by /city command\n/forecast - show's forecast by prefered city/help", nil)
}
