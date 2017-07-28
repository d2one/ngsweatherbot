package main

import (
	"log"

	"github.com/bot-api/telegram"
	"golang.org/x/net/context"
)

func runNotificationTasks(ctx context.Context, api *telegram.API) error {
	log.Println("start job")
	userNotifications, err := ds.getCronUserNotification()
	if err != nil {
		logWork(err)
		return err
	}

	for _, userNotification := range userNotifications {
		userCity, err := ds.getUserCity(userNotification.UserID)
		if userCity == nil {
			logWork(err)
			continue
		}

		forecast, err := ws.getForecast(userCity.CityAlias)
		if err != nil {
			continue
		}
		textMessage := userCity.CityTitle + "\n" + ws.formatForecasttWeather(forecast)
		msg := telegram.NewMessage(userNotification.ChatID, textMessage)
		if _, err := api.Send(ctx, msg); err != nil {
			logWork(err)
		}
		userNotification.NextRun += 60 * 60 * 24
		if err = ds.saveUserNotification(*userNotification); err != nil {
			logWork(err)
		}
	}
	return nil
}
