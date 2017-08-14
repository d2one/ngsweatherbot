package main

import (
	"context"
	"github.com/bot-api/telegram"
)

func runNotificationTasks(ctx context.Context, api *telegram.API) error {
	usersData, err := dataStore.getCronUsersNotifications()
	if err != nil {
		logWork(err)
		return err
	}

	for _, userData := range usersData {
		if !userData.CityAlias.Valid {
			logWork(err)
			continue
		}

		forecast, err := weatherService.getForecast(userData.CityAlias.String)
		if err != nil {
			continue
		}
		textMessage := userData.CityTitle.String + "\n" + weatherService.formatForecastWeather(forecast)
		msg := telegram.NewMessage(userData.ChatID, textMessage)
		msg.ParseMode = "markdown"
		if _, err := api.Send(ctx, msg); err != nil {
			logWork(err)
		}
		nextRun := 60*60*24 + userData.NotificationsNextRun.Int64
		if err = dataStore.saveUserNotification(userData.ChatID, nextRun); err != nil {
			logWork(err)
		}
	}
	return nil
}
