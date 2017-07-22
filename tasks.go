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

		currentWeather, err := ws.getCurrentWeather(userCity.CityAlias)
		if err != nil {
			logWork(err)
			continue
		}
		msg := telegram.NewMessage(userNotification.ChatID, ws.formatCurrentWeather(currentWeather))
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
