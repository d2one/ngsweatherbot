package main

type WeatherIcon struct {
	icons map[string]map[string]string
}

// dasdas
func NewWeatherIcon() *WeatherIcon {
	icons := map[string]map[string]string{
		"buttons": map[string]string{
			"back":          "\xF0\x9F\x94\x99",
			"default_city":  "\xF0\x9F\x8F\xA2",
			"notifications": "\xF0\x9F\x94\x94",
			"settings":      "\xF0\x9F\x94\xA7",
			"forecast":      "\xF0\x9F\x93\xAF Прогноз",
			"current":       "\xF0\x9F\x94\x86 Сейчас",
		},
		"wind": map[string]string{
			"north_west": "\xE2\x86\x96",
			"north_east": "\xE2\x86\x97",
			"south_west": "\xE2\x86\x99",
			"south_east": "\xE2\x86\x98",
			"north":      "\xE2\xAC\x87",
			"east":       "\xE2\xAC\x86",
			"south":      "\xE2\x9E\xA1",
			"west":       "\xE2\xAC\x85",
		},
		"clouds": map[string]string{
			"sunshine":      "\xF0\x9F\x94\x86",
			"partly_cloudy": "\xF0\x9F\x8C\x84",
			"mostly_cloudy": "\xE2\x9B\x85",
			"cloudy":        "\xE2\x98\x81",
			"light_cloudy":  "\xF0\x9F\x8C\x87",
		},
		"precipitations": map[string]string{
			"rain":           "\xE2\x98\x94",
			"rainlight_rain": "\xE2\x98\x94",
			"heavy_rain":     "\xE2\x98\x94",
			"snow":           "\xE2\x9D\x84",
			"light_snow":     "\xE2\x9D\x84",
			"heavy_snow":     "\xE2\x9D\x84",
			"sleet":          "\xE2\x9D\x84",
			"thunderstorm":   "\xE2\x9A\xA1",
		},
	}
	return &WeatherIcon{
		icons: icons,
	}
}

func (wi *WeatherIcon) getWind(key string) string {
	icon, ok := wi.icons["wind"][key]
	if ok {
		return icon
	}
	return ""
}
func (wi *WeatherIcon) getClouds(key string) string {
	icon, ok := wi.icons["clouds"][key]
	if ok {
		return icon
	}
	return ""
}

func (wi *WeatherIcon) getPrecipitations(key string) string {
	icon, ok := wi.icons["precipitations"][key]
	if ok {
		return icon
	}
	return ""
}
func (wi *WeatherIcon) getButtons(key string) string {
	icon, ok := wi.icons["buttons"][key]
	if ok {
		return icon
	}
	return ""
}
