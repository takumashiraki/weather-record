package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ResponseJson struct {
	Response map[string]interface{}
}

const (
	apiEndpoint = "https://weather.tsukumijima.net/api/forecast/city"
	cityId      = "011000"
)

func PutDataDB(request events.APIGatewayProxyRequest) (*ResponseJson, error) {
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
