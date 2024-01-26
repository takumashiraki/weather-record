package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	apis "hello-world/apis"
)

type ResponseJson struct {
	Response map[string]interface{}
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

	// 天気予報APIから天気の情報を取得する
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

	// 天気予報APIから天気の情報を取得する
	apis.PutDataDB(response)
	// response, err := operation.PutDataDB()
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		Body:       greeting,
	// 		StatusCode: 500,
	// 	}, err
	// }

	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
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

	var responseBody map[string]interface{}
	// json.UnmarshalでJSONデータをGoのオブジェクトに変換する
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		fmt.Println("レスポンスボディJSONに変換エラー:", err)
		return nil, err
	}
	r.Response = responseBody
	return &r, nil
}

func main() {
	lambda.Start(handler)
}
