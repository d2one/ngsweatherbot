package main

import (
	"encoding/json"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		resp, err := http.Get("http://pogoda.ngs.ru/api/v1/forecasts/current?city=novosibirsk")
		if err != nil {
			// handle error
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		messageText := getStations([]byte(body))

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
		bot.Send(msg)

		// icon := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "img/cloudy_none_day_big.png")
		// icon.Caption = messageText
		// bot.Send(icon)

	}
}

func getStations(body []byte) string {
	var s = new(WeatherResponce)
	err := json.Unmarshal(body, &s)
	if err != nil {
		log.Printf("whoops:", err)
	}

	messageText := strconv.FormatFloat(float64(s.Forecasts[0].Temperature), 'f', 1, 32) + "°C, " + "Ветер: " + strconv.FormatFloat(float64(s.Forecasts[0].Wind.Speed), 'f', 1, 32) + "м/с, " + s.Forecasts[0].Wind.Direction.Title + ", " + s.Forecasts[0].Cloud.Title + " " + s.Forecasts[0].Precipitation.Title

	return messageText
}

type WeatherResponce struct {
	Forecasts []struct {
		Astronomy struct {
			LengthDayHuman  string `json:"length_day_human"`
			MoonIlluminated int    `json:"moon_illuminated"`
			MoonPhase       string `json:"moon_phase"`
			Sunrise         string `json:"sunrise"`
			Sunset          string `json:"sunset"`
		} `json:"astronomy"`
		Cloud struct {
			Name  string `json:"name"`
			Title string `json:"title"`
			Value string `json:"value"`
		} `json:"cloud"`
		Date                string      `json:"date"`
		EcologicalSituation interface{} `json:"ecological_situation"`
		FeelLikeTemperature float64     `json:"feel_like_temperature"`
		Humidity            int         `json:"humidity"`
		Icon                string      `json:"icon"`
		IconPath            string      `json:"icon_path"`
		Links               struct {
			City string `json:"city"`
		} `json:"links"`
		MagneticStatus string `json:"magnetic_status"`
		Precipitation  struct {
			DayValue int    `json:"day_value"`
			Title    string `json:"title"`
			Value    string `json:"value"`
		} `json:"precipitation"`
		Pressure         int     `json:"pressure"`
		SolarRadiation   int     `json:"solar_radiation"`
		Source           string  `json:"source"`
		Temperature      float64 `json:"temperature"`
		TemperatureTrend float64 `json:"temperature_trend"`
		UpdateDate       string  `json:"update_date"`
		UvIndex          float64 `json:"uv_index"`
		Water            []struct {
			Level struct {
				Hint  string `json:"hint"`
				Trend int    `json:"trend"`
				Value int    `json:"value"`
			} `json:"level"`
			Temperature interface{} `json:"temperature"`
			Title       string      `json:"title"`
			WaveHeight  interface{} `json:"wave_height"`
		} `json:"water"`
		Wind struct {
			Direction struct {
				Title       string `json:"title"`
				TitleLetter string `json:"title_letter"`
				TitleShort  string `json:"title_short"`
				Value       string `json:"value"`
			} `json:"direction"`
			Speed float64 `json:"speed"`
		} `json:"wind"`
	} `json:"forecasts"`
	Metadata struct {
		Resultset struct {
			Count int `json:"count"`
		} `json:"resultset"`
	} `json:"metadata"`
}
