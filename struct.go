package main

type UserCity struct {
	id	int64
	user_id	int64
	chat_id	int64
	city_alias	string
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
		Links struct {
			City string `json:"city"`
		} `json:"links"`
		MagneticStatus string `json:"magnetic_status"`
		Precipitation struct {
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
		Water []struct {
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
