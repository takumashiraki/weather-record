package apis

import (
	"fmt"
	"strconv"
)

type ResponseJson struct {
	Response map[string]interface{}
}

type Reagion struct {
	Id   string
	Name string
	Pref string
}

type Temperature struct {
	Id   string
	Date string  `json:"date"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
}

type Detail struct {
	Id      string
	Date    string `json:"date"`
	Weather string `json:"weather"`
	Wind    string `json:"wind"`
	Wave    string `json:"wave"`
	Telop   string `json:"telop"`
}

type ChanceOfRain struct {
	Id     string
	Date   string `json:"date"`
	T00_06 string `json:"T00_06"`
	T06_12 string `json:"T06_12"`
	T12_18 string `json:"T12_18"`
	T18_24 string `json:"T18_24"`
}

// PutDataDB関数の修正
func PutDataDB(response ResponseJson) {
	fmt.Println("hogehoge 😃")

	forecasts, ok := response.Response["forecasts"].([]interface{})
	if !ok {
		fmt.Println("Error: forecasts is not an array")
		return
	}

	for _, forecastInterface := range forecasts {
		// 構造体を初期化する
		d := new(Detail)
		t := new(Temperature)
		c := new(ChanceOfRain)

		forecast, ok := forecastInterface.(map[string]interface{})
		if !ok {
			fmt.Println("Error: forecast is not a map[string]interface{}")
			continue
		}

		dateLabel, ok := forecast["dateLabel"].(string)
		if !ok {
			fmt.Println("Error: dateLabel is not a string")
			continue
		}
		fmt.Printf("Date Label: %s\n", dateLabel)
		if dateLabel != "明日" {
			fmt.Printf("Date Labelが%sでした\n", dateLabel)
			fmt.Println("--------")
			continue
		}

		detail, ok := forecast["detail"].(map[string]interface{})
		if !ok {
			fmt.Println("detailのデータ取得に失敗")
			break
		}
		wave := detail["wave"].(string)
		weather := detail["weather"].(string)
		wind := detail["wind"].(string)
		d.Wave = wave
		d.Weather = weather
		d.Wind = wind

		fmt.Printf("Wave: %s\nWeather: %s\nWind: %s\n", d.Wave, d.Weather, d.Wind)

		telop, ok := forecast["telop"].(string)
		if !ok {
			fmt.Println("telopのデータ取得に失敗")
			break
		}
		d.Telop = telop
		fmt.Printf("Telop: %s\n", d.Telop)

		temperature, ok := forecast["temperature"].(map[string]interface{})
		if !ok {
			fmt.Println("temperatureのデータ取得に失敗")
			break
		}
		maxTemperature, ok := temperature["max"].(map[string]interface{})
		if !ok {
			fmt.Println("max temperatureのデータ取得に失敗")
			break
		}
		maxCelsiusInterface, ok := maxTemperature["celsius"]
		if !ok {
			fmt.Println("Error: maxCelsiusInterface is not a interface{}")
			break
		}
		maxCelsiusStr := fmt.Sprintf("%v", maxCelsiusInterface)
		maxCelsius, maxErr := strconv.ParseFloat(maxCelsiusStr, 64)
		if maxErr != nil {
			fmt.Println("Error:", maxErr)
			break
		}
		t.Max = maxCelsius
		fmt.Printf("Max Temperature: %v\n", t.Max)

		minTemperature, ok := temperature["min"].(map[string]interface{})
		if !ok {
			fmt.Println("min temperatureのデータ取得に失敗")
			break
		}
		minCelsiusInterface, ok := minTemperature["celsius"]
		if !ok {
			fmt.Println("Error: maxCelsiusInterface is not a interface{}")
			break
		}
		minCelsiusStr := fmt.Sprintf("%v", minCelsiusInterface)
		minCelsius, minErr := strconv.ParseFloat(minCelsiusStr, 64)
		if minErr != nil {
			fmt.Println("Error:", minErr)
			break
		}
		// fmt.Printf("minCelsius as float: %s\n", minCelsius)
		t.Min = minCelsius
		fmt.Printf("Min Temperature: %v\n", t.Min)

		chanceOfRain, ok := forecast["chanceOfRain"].(map[string]interface{})
		if !ok {
			fmt.Println("Error:", minErr)
			break
		}
		t00_06 := chanceOfRain["T00_06"].(string)
		t06_12 := chanceOfRain["T06_12"].(string)
		t12_18 := chanceOfRain["T12_18"].(string)
		t18_24 := chanceOfRain["T18_24"].(string)

		c.T00_06 = t00_06
		c.T06_12 = t06_12
		c.T12_18 = t12_18
		c.T18_24 = t18_24
		fmt.Printf("Chance of Rain (T00_06): %s\n", c.T00_06)
		fmt.Printf("Chance of Rain (T06_12): %s\n", c.T06_12)
		fmt.Printf("Chance of Rain (T12_18): %s\n", c.T12_18)
		fmt.Printf("Chance of Rain (T18_24): %s\n", c.T18_24)

		fmt.Println("--------")
	}
	// ... (後略)

	return
}
