package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	api := telegram.New("TOKEN")
	api.Debug(true)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic
	boltStorage := NewBoltStorage("bolt.db")
	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(func(ctx context.Context) error {
		update := telebot.GetUpdate(ctx) // take update from context
		if update.Message == nil {
			return nil
		}
		city, error := getCity(update.Message.Text)

		if error != "" {
			_, err := api.SendMessage(ctx,
				telegram.NewMessagef(update.Chat().ID, error))
			return err
		}

		textMessage := getStations(city)
		api := telebot.GetAPI(ctx) // take api from context
		msg := telegram.NewMessage(update.Chat().ID, textMessage)
		_, err := api.Send(ctx, msg)
		return err

	})
	// Use command middleware, that helps to work with commands
	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {

				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				_, err := api.SendMessage(ctx,
					telegram.NewMessagef(update.Chat().ID,
						"received start with arg %s", arg,
					))

				return err
			}),
		"current": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {
				update := telebot.GetUpdate(ctx)
				textMessage := "No selected city. Select city with command \n/city {cityName}"

				boltStorage.DB.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("users"))
					v := b.Get([]byte(strconv.FormatInt(update.From().ID, 10)))
					if v != nil {
						textMessage = getStations(string(v))
					}
					return nil
				})

				api := telebot.GetAPI(ctx) // take api from context
				msg := telegram.NewMessage(update.Chat().ID, textMessage)
				_, err := api.Send(ctx, msg)
				return err
			}),
		"help": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {

				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				_, err := api.SendMessage(ctx,
					telegram.NewMessagef(update.Chat().ID,
						"It's Ngs weather bot. \nIt can show you weather in city. To start messaging, you can send city name in your message.\nCommands:\n/city {cityName} - set your prefered city to show by /curent command\n/current - show's the weather in prefered city by /city command\n/forecast - show's forecast by prefered city/help",
					))
				return err
			}),
		"city": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {

				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				user := update.From()

				city, error := getCities(arg)

				if error != "" {
					_, err := api.SendMessage(ctx,
						telegram.NewMessagef(update.Chat().ID, error))
					return err
				}

				if len(city.Cities) > 1 {
					textMessage := "Please, set city:\n"
					for index := range city.Cities {
						log.Println(index)
						log.Println(city.Cities[index].Alias)
						textMessage += "/city " + city.Cities[index].Alias + "\n"
					}
					api.SendMessage(ctx,
						telegram.NewMessagef(update.Chat().ID,
							textMessage,
						))
					return nil
				}

				boltStorage.writerChan <- [3]interface{}{"users", strconv.FormatInt(user.ID, 10), []byte(city.Cities[0].Alias)}
				api.SendMessage(ctx,
					telegram.NewMessagef(update.Chat().ID,
						"City selected: %s", city.Cities[0].Title,
					))
				return nil
			}),
	}))

	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}

func getCities(arg string) (*WeatherCitys, string) {
	resp, err := http.Get("http://pogoda.ngs.ru/api/v1/cities?q=" + arg)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var s = new(WeatherCitys)
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Printf("whoops:", err)
	}

	if s.Errors.Message != "" {
		return s, s.Errors.Message
	}

	return s, ""
}

func getCity(arg string) (string, string) {
	s, err := getCities(arg)
	if err != "" {
		return "", err
	}
	return s.Cities[0].Alias, ""

}

func getStations(arg string) string {
	resp, err := http.Get("http://pogoda.ngs.ru/api/v1/forecasts/current?city=" + arg)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var s = new(WeatherResponce)
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Printf("whoops:", err)
	}

	messageText := strconv.FormatFloat(float64(s.Forecasts[0].Temperature), 'f', 1, 32) + "°C, " + "Ветер: " + strconv.FormatFloat(float64(s.Forecasts[0].Wind.Speed), 'f', 1, 32) + "м/с, " + s.Forecasts[0].Wind.Direction.Title + ", " + s.Forecasts[0].Cloud.Title + " " + s.Forecasts[0].Precipitation.Title

	return messageText
}

func (this *BoltStorage) writer() {
	for data := range this.writerChan {
		bucket := data[0].(string)
		keyId := data[1].(string)
		dataBytes := data[2].([]byte)
		err := this.DB.Update(func(tx *bolt.Tx) error {
			sesionBucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
			sesionBucket.Delete([]byte(keyId))
			return sesionBucket.Put([]byte(keyId), dataBytes)
		})
		if err != nil {
			// TODO: Handle instead of panic
			panic(err)
		}
	}
}

func NewBoltStorage(dbPath string) *BoltStorage {
	db, err := bolt.Open(dbPath, 0666, nil)
	writerChan := make(chan [3]interface{})
	boltStorage := &BoltStorage{DB: db, writerChan: writerChan}
	boltStorage.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		return nil
	})
	go boltStorage.writer()
	if err != nil {
		panic(err)
	}
	return boltStorage
}

type WeatherCitys struct {
	Cities []struct {
		Alias              string      `json:"alias"`
		ID                 int         `json:"id"`
		MobileURL          string      `json:"mobile_url"`
		Name               string      `json:"name"`
		ProjectName        string      `json:"project_name"`
		Region             int         `json:"region"`
		Timezone           string      `json:"timezone"`
		Title              string      `json:"title"`
		TitleDative        string      `json:"title_dative"`
		TitleForIos        interface{} `json:"title_for_ios"`
		TitlePrepositional string      `json:"title_prepositional"`
		URL                string      `json:"url"`
	} `json:"cities"`
	Errors struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Metadata struct {
		Resultset struct {
			Count  int `json:"count"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		} `json:"resultset"`
	} `json:"metadata"`
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

type BoltStorage struct {
	DB         *bolt.DB
	writerChan chan [3]interface{} //not so agnostic but enough now
}
