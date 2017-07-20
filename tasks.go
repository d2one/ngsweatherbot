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
		log.Printf("%s", err)
		return err
	}

	for _, userNotification := range userNotifications {
		if userCity, _ := ds.getUserCity(userNotification.UserID); userCity != nil {
			if currentWeather, err := ws.getCurrentWeather(userCity.CityAlias); err == nil {
				msg := telegram.NewMessage(userNotification.ChatID, ws.formatCurrentWeather(currentWeather))
				if _, err := api.Send(ctx, msg); err != nil {
					log.Println(err)
				}
				userNotification.NextRun += 60 * 60 * 24
				if err = ds.saveUserNotification(*userNotification); err != nil {
					log.Println(err)
				}
			}
		}
	}
	return nil
}
