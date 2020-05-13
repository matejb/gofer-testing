package wittr_test

import (
	"fmt"
	"log"

	"github.com/matejb/gofer-testing/weather/with-complete-tests/wittr"
)

func ExampleWeatherSource_CurrentTemperature() {
	wtr := wittr.New("Zagreb")

	cleanup := setupFakeConnection(wtr) // not needed in real code
	defer cleanup()                     // not needed in real code

	temperature, err := wtr.CurrentTemperature()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temperature in Zagreb is %d °C", temperature)

	// Output:
	// Temperature in Zagreb is 20 °C
}
