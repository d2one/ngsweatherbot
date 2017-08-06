package main

// UserCity user selected city
type UserCity struct {
	ChatID       int64
	CityAlias    string
	CityTitle    string
	ForecastType string
}

// WeatherSource WeatherSource
type WeatherSource interface {
	getCities(arg string) ([]*City, error)
	getCity(arg string) (*City, error)
	getCurrentWeather(arg string) (*CurrentWeather, error)
	getForecast(arg string) (*WeatherResponceForecasts, error)
}

// CacheService CacheService
type CacheService interface {
	read(cacheKey string) (*CachedItem, error)
	write(cacheKey string, cacheValue string, ttl int64) error
	delete(cacheKey string) error
}

// UserNotification user notification
type UserNotification struct {
	ID      int64
	UserID  int64
	ChatID  int64
	NextRun int64
}

// CachedItem cache
type CachedItem struct {
	ID         int64
	CacheKey   string
	CacheValue string
	TTL        int64
	TTLLock    int64
}

// City ngs weather api city responce
type City struct {
	Alias string `json:"alias"`
	Title string `json:"title"`
}

// WeatherCities cities
type WeatherCities struct {
	Cities []*City `json:"cities"`
	Errors struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

// CurrentWeather current weather
type CurrentWeather struct {
	Astronomy struct {
		LengthDayHuman  string `json:"length_day_human"`
		MoonIlluminated int    `json:"moon_illuminated"`
		MoonPhase       string `json:"moon_phase"`
		Sunrise         string `json:"sunrise"`
		Sunset          string `json:"sunset"`
	} `json:"astronomy"`
	FeelLikeTemperature float64 `json:"feel_like_temperature"`
	Cloud               struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"cloud"`
	Precipitation struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"precipitation"`
	Temperature float64 `json:"temperature"`
	Wind        struct {
		Direction struct {
			Value string `json:"value"`
			Title string `json:"title"`
		} `json:"direction"`
		Speed float64 `json:"speed"`
	} `json:"wind"`
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
	IconPath string `json:"icon_path"`
	Links    struct {
		City string `json:"city"`
	} `json:"links"`
}

// WeatherResponce weather response
type WeatherResponce struct {
	Forecasts []*CurrentWeather `json:"forecasts"`
}

type HourForecat struct {
	Hour        int `json:"hour"`
	Temperature struct {
		Avg float64 `json:"avg"`
	} `json:"temperature"`
	Wind struct {
		Speed struct {
			Avg float64 `json:"avg"`
		} `json:"speed"`
		Direction struct {
			Value string `json:"value"`
			Title string `json:"title"`
		} `json:"direction"`
	} `json:"wind"`
	Cloud struct {
		Value string `json:"value"`
		Title string `json:"title"`
	} `json:"cloud"`
	Precipitation struct {
		Value string `json:"value"`
		Title string `json:"title"`
	} `json:"precipitation"`
}

type WeatherResponceForecasts struct {
	Forecasts []struct {
		Date  string         `json:"date"`
		Hours []*HourForecat `json:"hours"`
	} `json:"forecasts"`
}
