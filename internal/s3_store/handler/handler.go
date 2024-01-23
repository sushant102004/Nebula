/*
	@author: Sushant
	@last-modified: 23 January 2024
	@GitHub: https://github.com/sushant102004
*/

package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sushant102004/Nebula/pkg/utils"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Ignoring all the requests where method is not "POST".
	if req.HTTPMethod != http.MethodPost {
		err := utils.ReturnErrorResponse(map[string]string{
			"error": "Method not allowed",
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       err,
		}, nil
	}

	resp := utils.ReturnResponse(map[string]string{
		"message": "Request Successful",
	})

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       resp,
	}, nil
}
