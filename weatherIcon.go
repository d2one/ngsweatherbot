package main

//WeatherIcons ds
type WeatherIcons struct {
	icons map[string]map[string]string
}

// NewWeatherIcons ds
func NewWeatherIcons() *WeatherIcons {
	icons := map[string]map[string]string{
		"buttons": {
			"back":                 "\xF0\x9F\x94\x99",
			"default_city":         "\xF0\x9F\x8F\xA4",
			"notifications":        "\xF0\x9F\x94\x94",
			"notifications_set":    "\xE2\x8F\xB0",
			"notifications_remove": "\xE2\x9D\x8C",
			"settings":             "\xF0\x9F\x94\xA7",
			"forecast":             "\xF0\x9F\x93\xAF Прогноз",
			"current":              "\xF0\x9F\x94\x86 Сейчас",
		},
		"wind": {
			"north_west": "\xE2\x86\x96",
			"north_east": "\xE2\x86\x97",
			"south_west": "\xE2\x86\x99",
			"south_east": "\xE2\x86\x98",
			"north":      "\xE2\xAC\x87",
			"east":       "\xE2\xAC\x86",
			"south":      "\xE2\x9E\xA1",
			"west":       "\xE2\xAC\x85",
		},
		"icons": {
			"north_west": "\xE2\x86\x96",
			"north_east": "\xE2\x86\x97",
			"south_west": "\xE2\x86\x99",
			"south_east": "\xE2\x86\x98",
			"north":      "\xE2\xAC\x87",
			"east":       "\xE2\xAC\x86",
			"south":      "\xE2\x9E\xA1",
			"west":       "\xE2\xAC\x85",
		},
		"clouds": {
			"sunshine":      "\xF0\x9F\x94\x86",
			"partly_cloudy": "\xE2\x9B\x85",
			"mostly_cloudy": "\xE2\x9B\x85",
			"cloudy":        "\xE2\x98\x81",
			"light_cloudy":  "\xF0\x9F\x8C\x87",
		},
		"precipitations": {
			"rain":           "\xE2\x98\x94",
			"rainlight_rain": "\xE2\x98\x94",
			"heavy_rain":     "\xE2\x98\x94\xE2\x98\x94",
			"snow":           "\xE2\x9D\x84",
			"light_snow":     "\xE2\x9D\x84",
			"heavy_snow":     "\xE2\x9D\x84",
			"sleet":          "\xE2\x9D\x84",
			"thunderstorm":   "\xE2\x9A\xA1",
		},
		"astronomy": {
			"sunrise": "\xF0\x9F\x8C\x85",
			"sunset":  "\xF0\x9F\x8C\x87",
		},
	}
	return &WeatherIcons{
		icons: icons,
	}
}

func (icons *WeatherIcons) getWind(key string) string {
	if val, ok := icons.icons["wind"][key]; ok {
		return val
	}
	return ""
}

func (icons *WeatherIcons) getAstronomy(key string) string {
	if val, ok := icons.icons["astronomy"][key]; ok {
		return val
	}
	return ""
}

func (icons *WeatherIcons) getWeatherIcons(key string) string {
	if val, ok := icons.icons["icons"][key]; ok {
		return val
	}
	return ""
}

func (icons *WeatherIcons) getClouds(key string) string {
	if val, ok := icons.icons["clouds"][key]; ok {
		return val
	}
	return ""
}

func (icons *WeatherIcons) getPrecipitations(key string) string {
	if val, ok := icons.icons["precipitations"][key]; ok {
		return val
	}
	return ""
}

func (icons *WeatherIcons) getButtons(key string) string {
	if val, ok := icons.icons["buttons"][key]; ok {
		return val
	}
	return ""
}
