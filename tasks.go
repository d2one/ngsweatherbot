package main

import (
	"github.com/bot-api/telegram"
	"golang.org/x/net/context"
	"log"
)

func runNotificationTasks(api *telegram.API, ctx context.Context) error {
	log.Println("start job")
	userNotifications, err := getCronUserNotification()
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	for _, userNotification := range userNotifications {
		userCity, _ := getUserCity(userNotification.user_id)
		if userCity != (UserCity{}) {
			if currentWeather, err := getCurrentWeather(userCity.city_alias); err == nil {
				msg := telegram.NewMessage(userNotification.chat_id, formatCurrentWeather(currentWeather))
				_, err := api.Send(ctx, msg)
				if err != nil {
					log.Println(err)
				}
				userNotification.next_run += 60 * 60 * 24
				err = saveUserNotification(userNotification)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	return nil
}
