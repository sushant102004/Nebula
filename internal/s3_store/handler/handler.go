/*
	@author: Sushant
	@last-modified: 24 January 2024
	@GitHub: https://github.com/sushant102004
*/

package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	customErrors "github.com/sushant102004/Nebula/pkg/errors"
	"github.com/sushant102004/Nebula/pkg/response"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Configuration to pretty print logs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Ignoring all the requests where method is not "POST".
	if req.HTTPMethod != http.MethodPost {
		return &events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       customErrors.MethodNotAllowed,
		}, nil
	}

	region := os.Getenv("region")
	accessKey := os.Getenv("access_key")
	secretKey := os.Getenv("secret_key")

	if region == "" || accessKey == "" || secretKey == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       customErrors.UnableToFindEnvVariable,
		}, nil
	}

	credentials := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),

		config.WithCredentialsProvider(credentials),
	)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       customErrors.UnableToLoadAWSConfiguration,
		}, nil
	}

	_ = s3.NewFromConfig(cfg)

	log.Info().Msg("Connected to S3")

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       response.SuccessfulResponse,
	}, nil
}
