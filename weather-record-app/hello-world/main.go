// 取得できるやつ

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

// --- API Responseを受け取るための構造体 ---
// Forecast は天気予報の情報を表す構造体
type Forecast struct {
	Date           string `json:"date"`
	DateLabel      string `json:"dateLabel"`
	Telop          string `json:"telop"`
	Weather        string `json:"weather"`
	Wind           string `json:"wind"`
	Wave           string `json:"wave"`
	MinTemperature *struct {
		Celsius    *float64 `json:"celsius"`
		Fahrenheit *float64 `json:"fahrenheit"`
	} `json:"temperature"`
	MaxTemperature *struct {
		Celsius    *float64 `json:"celsius"`
		Fahrenheit *float64 `json:"fahrenheit"`
	} `json:"temperature"`
	ChanceOfRain struct {
		T00_06 string `json:"T00_06"`
		T06_12 string `json:"T06_12"`
		T12_18 string `json:"T12_18"`
		T18_24 string `json:"T18_24"`
	} `json:"chanceOfRain"`
	Image struct {
		Title  string `json:"title"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image"`
}

// Location は地域情報を表す構造体
type Location struct {
	Area       string `json:"area"`
	Prefecture string `json:"prefecture"`
	District   string `json:"district"`
	City       string `json:"city"`
}

// Copyright は著作権情報を表す構造体
type Copyright struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Image struct {
		Title  string `json:"title"`
		Link   string `json:"link"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image"`
	Provider []struct {
		Link string `json:"link"`
		Name string `json:"name"`
		Note string `json:"note"`
	} `json:"provider"`
}

// ResponseJson はAPIからのレスポンスを表す構造体
type ResponseJson struct {
	PublicTime          string `json:"publicTime"`
	PublicTimeFormatted string `json:"publicTimeFormatted"`
	PublishingOffice    string `json:"publishingOffice"`
	Title               string `json:"title"`
	Link                string `json:"link"`
	Description         struct {
		PublicTimeFormatted string `json:"publicTimeFormatted"`
		HeadlineText        string `json:"headlineText"`
		BodyText            string `json:"bodyText"`
		Text                string `json:"text"`
	} `json:"description"`
	Forecasts []Forecast `json:"forecasts"`
	Location  Location   `json:"location"`
	Copyright Copyright  `json:"copyright"`
}

// --- DBにデータを格納するための構造体 ---
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

// ForecastData は API レスポンスから変換したデータを格納する構造体
type ForecastData struct {
	Reagion      Reagion
	Temperatures []Temperature
	Details      []Detail
	Chances      []ChanceOfRain
}

const (
	apiEndpoint = "https://weather.tsukumijima.net/api/forecast/city"
	cityId      = "011000"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var greeting string
	sourceIP := request.RequestContext.Identity.SourceIP

	if sourceIP == "" {
		greeting = "Hello, world!\n"
	} else {
		greeting = fmt.Sprintf("Hello, %s!\n", sourceIP)
	}

	// handler関数内でrequest関数を呼び出す
	response, err := getWeatherInfo(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       greeting,
			StatusCode: 500,
		}, err
	}

	putDataDB(response, request)

	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func getWeatherInfo(request events.APIGatewayProxyRequest) (*ResponseJson, error) {
	var r ResponseJson
	// URLの組み立て
	url := fmt.Sprintf("%s/%s", apiEndpoint, cityId)
	// HTTP GETリクエストの送信
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("リクエストエラー:", err)
		return nil, err
	}
	defer response.Body.Close()

	// レスポンスボディの読み取り
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("レスポンスボディ読み取りエラー:", err)
		return nil, err
	}

	// json.UnmarshalでJSONデータをGoのオブジェクトに変換する
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("レスポンスボディJSONに変換エラー:", err)
		return nil, err
	}

	return &r, nil
}

func putDataDB(response *ResponseJson, request events.APIGatewayProxyRequest) {
	fmt.Println("😄 1")
	todayForecasts := response.Forecasts[1]

	// 新しい構造体にデータを詰め替える
	forecastData := &ForecastData{
		Reagion: Reagion{
			Id:   cityId,
			Name: response.Location.City,
			Pref: response.Location.Prefecture,
		},
	}
	fmt.Println("😄 2")
	fmt.Println(forecastData)

	temperature := Temperature{
		Id:   cityId,
		Date: todayForecasts.Date,
	}
	// MinTemperatureがnilでない場合、Celsiusをセット
	if todayForecasts.MinTemperature != nil && todayForecasts.MinTemperature.Celsius != nil {
		temperature.Min = *todayForecasts.MinTemperature.Celsius
	}
	// MaxTemperatureがnilでない場合、Celsiusをセット
	if todayForecasts.MaxTemperature != nil && todayForecasts.MaxTemperature.Celsius != nil {
		temperature.Max = *todayForecasts.MaxTemperature.Celsius
	}
	forecastData.Temperatures = append(forecastData.Temperatures, temperature)

	fmt.Println("😄 4")
	fmt.Println(temperature)
	// Details
	detail := Detail{
		Id:      cityId,
		Date:    todayForecasts.Date,
		Weather: todayForecasts.Weather,
		Wind:    todayForecasts.Wind,
		Wave:    todayForecasts.Wave,
		Telop:   todayForecasts.Telop,
	}
	forecastData.Details = append(forecastData.Details, detail)

	fmt.Println("😄 5")
	fmt.Println(detail)
	// Chances
	chance := ChanceOfRain{
		Id:     cityId,
		Date:   todayForecasts.Date,
		T00_06: todayForecasts.ChanceOfRain.T00_06,
		T06_12: todayForecasts.ChanceOfRain.T06_12,
		T12_18: todayForecasts.ChanceOfRain.T12_18,
		T18_24: todayForecasts.ChanceOfRain.T18_24,
	}
	forecastData.Chances = append(forecastData.Chances, chance)

	fmt.Println("😄 6")
	fmt.Println(chance)
	// レスポンス表示

	fmt.Println("😄 7")
	fmt.Println("posgreに接続する")

	var dbName string = "rds-for-postgeresql-weather-record"
	var dbPassword string = "weather_password_2024"
	var dbUser string = "weather_reporter"
	var dbHost string = "rds-for-postgeresql-weather-record.c1owwq2mqjfe.ap-northeast-1.rds.amazonaws.com"
	var dbPort int = 5432

	// PostgreSQL 接続文字列
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", dbHost, dbPort, dbUser, dbPassword, dbName)

	// PostgreSQL に接続
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	fmt.Println("😄 8")
	fmt.Println("posgreに接続後")
}
