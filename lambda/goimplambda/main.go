package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kacperjurak/goimp"
	"github.com/kacperjurak/goimpcore"
	"net/http"
)

func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := goimp.Request{}
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		panic(err)
	}

	result := solve(request.Task)

	response := goimp.Response{}

	response.Index = request.Index
	response.Result = result
	response.Result.InitValues = request.Task.InitValues

	r, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(r),
	}, nil

}

func solve(task goimp.Task) goimp.Result {
	var (
		freqs   = task.Freqs
		impData = task.ImpData
	)

	s := goimpcore.NewSolver(task.Code, freqs, impData)
	s.InitValues = task.InitValues
	s.SmartMode = "eis"

	result := s.Solve(10, 1000)

	return result
}

func main() {
	lambda.Start(HandleRequest)
}
