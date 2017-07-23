package main

// UserCity user selected city
type UserCity struct {
	ID        int64
	UserID    int64
	ChatID    int64
	CityAlias string
	CityTitle string
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
	IconPath string `json:"icon_path"`
}

// WeatherResponce weather response
type WeatherResponce struct {
	Forecasts []*CurrentWeather `json:"forecasts"`
}

type WeatherResponceForecasts struct {
	Forecasts []struct {
		Date  string `json:"date"`
		Hours []struct {
			Hour        int `json:"hour"`
			Temperature struct {
				Avg int `json:"avg"`
			} `json:"temperature"`
			Pressure struct {
				Avg int `json:"avg"`
			} `json:"pressure"`
			Wind struct {
				Speed struct {
					Avg int `json:"avg"`
				} `json:"speed"`
				Direction struct {
					Title string `json:"title"`
				} `json:"direction"`
			} `json:"wind"`
			Cloud struct {
				Title string `json:"title"`
			} `json:"cloud"`
			Precipitation struct {
				Title string `json:"title"`
			} `json:"precipitation"`
		} `json:"hours"`
	} `json:"forecasts"`
}
