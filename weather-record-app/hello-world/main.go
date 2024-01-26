// 取得できるやつ

package main

import (
	"fmt"
	"hello-world/apis"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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
	response, err := apis.PutDataDB(request)
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
