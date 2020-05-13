package wittr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// New returns WeatherSource preset with default values.
func New(city string) *WeatherSource {
	return &WeatherSource{
		BaseURL: "https://wttr.in",
		city:    city,
	}
}

// WeatherSource is source of weather data which polls it from wttr service.
type WeatherSource struct {
	// BaseURL is URL of where wttr service is located.
	// Default value is "https://wttr.in".
	BaseURL string

	city string
}

// CurrentTemperature returns the current temperature in °C.
func (ws *WeatherSource) CurrentTemperature() (int, error) {
	res, err := http.Get(ws.BaseURL + "/" + ws.city + "?format=j1")
	if err != nil {
		return 0, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return 0, fmt.Errorf("service call failed with %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	res.Body.Close()

	type weatherData struct {
		CurrentCondition []struct {
			FeelsLikeC *string
		} `json:"current_condition"`
	}
	var data weatherData

	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("decoding weather data failed: %w", err)
	}

	currentTempReceived := len(data.CurrentCondition) > 0 && data.CurrentCondition[0].FeelsLikeC != nil
	if !currentTempReceived {
		return 0, fmt.Errorf("current temperature not received")
	}

	temperature, err := strconv.Atoi(*data.CurrentCondition[0].FeelsLikeC)
	if err != nil {
		return 0, fmt.Errorf("bad value for current temperature received: %w", err)
	}

	return temperature, nil
}
