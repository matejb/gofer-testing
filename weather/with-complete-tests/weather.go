package main

import (
	"fmt"
	"log"

	"github.com/matejb/gofer-testing/weather/with-complete-tests/wittr"
)

func main() {
	shorts, err := IsForShorts(wittr.New("Zagreb"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("U %s je za kratke rukave: %v\n", "Zagreb", shorts)
}

type WeatherSource interface {
	CurrentTemperature() (int, error)
}

func IsForShorts(source WeatherSource) (bool, error) {
	temperature, err := source.CurrentTemperature()
	if err != nil {
		return false, err
	}

	forShorts := temperature > 20

	return forShorts, nil
}
