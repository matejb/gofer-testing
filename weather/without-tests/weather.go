package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	shorts, err := IsForShorts("Zagreb")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("U %s je za kratke rukave: %v\n", "Zagreb", shorts)
}

func IsForShorts(city string) (bool, error) {
	res, err := http.Get(`https://wttr.in/Zagreb?format=j1`)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// Data sample:
	// {
	// "current_condition": [
	//         {
	//             "FeelsLikeC": "17",
	//             "FeelsLikeF": "63",
	//             "cloudcover": "0",
	//             "humidity": "59",
	// .....

	type weatherData struct {
		CurrentCondition []struct {
			FeelsLikeC string
		} `json:"current_condition"`
	}
	var data weatherData

	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}

	if len(data.CurrentCondition) < 1 {
		return false, fmt.Errorf("current condition not given from remote source")
	}

	temperature, err := strconv.Atoi(data.CurrentCondition[0].FeelsLikeC)
	if err != nil {
		return false, err
	}

	forShorts := temperature > 20

	return forShorts, nil
}
