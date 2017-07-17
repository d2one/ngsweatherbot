package main

type UserCity struct {
	id         int64
	user_id    int64
	chat_id    int64
	city_alias string
}

type UserNotification struct {
	id       int64
	user_id  int64
	chat_id  int64
	next_run int64
}

type CachedItem struct {
	id          int64
	cache_key   string
	cache_value string
	ttl         int64
	ttl_lock    int64
}

type City struct {
	Alias string `json:"alias"`
	Title string `json:"title"`
}

type WeatherCitys struct {
	Cities []City `json:"cities"`
	Errors struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

type CurrentWeather struct {
	Cloud struct {
		Title string `json:"title"`
	} `json:"cloud"`
	Precipitation struct {
		Title string `json:"title"`
	} `json:"precipitation"`
	Temperature float64 `json:"temperature"`
	Wind        struct {
		Direction struct {
			Title string `json:"title"`
		} `json:"direction"`
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

type WeatherResponce struct {
	Forecasts []CurrentWeather `json:"forecasts"`
	Metadata  struct {
		Resultset struct {
			Count int `json:"count"`
		} `json:"resultset"`
	} `json:"metadata"`
}
