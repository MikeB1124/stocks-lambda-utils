package stockslambdautils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func CreateResponse(response Response) (events.APIGatewayProxyResponse, error) {
	responseBody, err := json.Marshal(response)
	if err != nil {
		responseBody, _ = json.Marshal(Response{Message: err.Error()})
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string(responseBody),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Body:       string(responseBody),
	}, nil
}
