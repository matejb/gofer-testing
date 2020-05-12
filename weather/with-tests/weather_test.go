package main_test

import (
	"testing"

	main "github.com/matejb/gofer-testing/weather/with-tests"
)

func TestIsForShorts(t *testing.T) {
	notForShortsData := mockWeatherSource{temperature: 10}

	isForShorts, err := main.IsForShorts(notForShortsData)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}

	if isForShorts != false {
		t.Errorf("expected %v got %v", false, isForShorts)
	}
}

type mockWeatherSource struct {
	temperature int
	err         error
}

func (m mockWeatherSource) CurrentTemperature() (int, error) {
	return m.temperature, m.err
}
