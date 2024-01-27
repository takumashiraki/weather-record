// 取得できるやつ

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	// レスポンス表示
	fmt.Println("レスポンス:")
	fmt.Println(response)

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
