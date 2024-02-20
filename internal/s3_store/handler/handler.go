/*
	@author: Sushant
	@GitHub: https://github.com/sushant102004

	PoF: -
	1. Store images to S3
	2. Allow user to specify whether to compress image before storing or not -> MVP
	3. Support for multiple images. -> Future Feature

	Working of this function: -
	1. Receive and validate request data from API Gateway.
	2. Compress image before storing (if compress : true in request body)
	3. Store images to S3 via multipart upload to handle large files.
*/

/*
	Environment variables for this function:-
	1. region -> Your aws region
	2. access_key -> Your aws account access key.
	3. secret_access_key -> Your aws account secret access key.
*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	customErrors "github.com/sushant102004/Nebula/pkg/errors"
	"github.com/sushant102004/Nebula/pkg/response"
)

func init() {
	// Using zerolog for better logging.
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Configuration to pretty print logs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	lambda.Start(Handler)
}

type RequestEvent struct {
	HTTPMethod string      `json:"http_method"`
	Body       RequestBody `json:"body"`
}

type RequestBody struct {
	// It will be used to save Base64 encoded image before storing into S3.
	ImageContent string `json:"image_content"`

	// The title with which this image will be stored in S3. Basically it is the Object Key.
	// Also make sure to prefix the with some unique value to prevent accidental overwriting. If there is file with same name than
	// S3 will overwrite the existing file
	FileName         string `json:"file_name"`
	BucketName       string `json:"bucket_name"`
	ApplyCompression bool   `json:"apply_compression"`

	// This refers to the quality of the image after compression. Less the quality higher the compression is.
	Quality string `json:"quality"`
}

func Handler(ctx context.Context, event RequestEvent) (*events.APIGatewayProxyResponse, error) {
	// Ignoring all the requests where method is not "POST".
	if event.HTTPMethod != http.MethodPost {
		return returnError(404, customErrors.MethodNotAllowed), nil
	}

	// Checking if any essential data is missing from request body
	if event.Body.FileName == "" || event.Body.ImageContent == "" || event.Body.BucketName == "" {
		return returnError(400, customErrors.InvalidInputBody), nil
	}

	region := os.Getenv("region")
	accessKey := os.Getenv("access_key")
	secretKey := os.Getenv("secret_key")

	if region == "" || accessKey == "" || secretKey == "" {
		return returnError(400, customErrors.UnableToFindEnvVariable), nil
	}

	// TODO: Check what does this session (3rd parameter in NewStaticCredentialsProvider) means.
	credentials := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials),
	)
	if err != nil {
		log.Error().Msgf("Unable to load aws default config variables: %v", err)
		return returnError(500, customErrors.UnableToLoadAWSConfiguration), nil
	}

	s3Client := s3.NewFromConfig(cfg)

	// Checking if user want to compress image before saving.

	if event.Body.ApplyCompression == true && event.Body.Quality != "" {
		compressImageFxUrl := os.Getenv("COMPRESS_IMAGE_FUNCTION_URL")
		if compressImageFxUrl == "" {
			log.Error().Msgf("compress image function url is not specified in .env file")
			return returnError(500, customErrors.UnableToFindEnvVariable), nil
		}

		// TODO: Currently I'm keeping these structs inside function for simplicity. Later I'll place them in dedicated directory.
		type reqBodyStruct struct {
			Image   string `json:"image"`
			Quality string `json:"quality"`
		}

		type responseBodyStruct struct {
			ResizedImage string `json:"resizedImage"`
			Message      string `json:"message"`
		}

		reqBody := reqBodyStruct{
			Image:   event.Body.ImageContent,
			Quality: event.Body.Quality,
		}

		var buff bytes.Buffer
		err = json.NewEncoder(&buff).Encode(reqBody)
		if err != nil {
			log.Error().Msgf("unable to encode request body")
			return returnError(500, customErrors.UnableToEncodeHTTPRequestBody), nil
		}

		httpClient := &http.Client{}

		req, err := http.NewRequest("POST", compressImageFxUrl, &buff)
		if err != nil {
			log.Error().Msgf("unable to create 'compress image' http request")
			return returnError(500, customErrors.UnableToCreateHTTPRequest), nil
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Error().Msgf("compress image http response err: %v", err.Error())
			return returnError(500, customErrors.HTTPResponseError), nil
		}

		var httpResp responseBodyStruct

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Msgf("unable to read data from HTTP response: %v", err.Error())
			return returnError(500, customErrors.UnableToReadDataFromHTTPResponse), nil
		}

		err = json.Unmarshal(respBytes, &httpResp)
		if err != nil {
			log.Error().Msgf("unable to marshal json: %v", err.Error())
			return returnError(500, customErrors.UnableToUnmarshalJSON), nil
		}

		if err := saveToS3(ctx, s3Client, event.Body.BucketName, event.Body.FileName, httpResp.ResizedImage); err != nil {
			log.Error().Msgf("Unable to store image to S3: %v", err)
			return returnError(500, customErrors.UnableToStoreImageToS3), nil
		}

		return &events.APIGatewayProxyResponse{StatusCode: 200, Body: response.SuccessfulResponse}, nil
	}

	if err := saveToS3(ctx, s3Client, event.Body.BucketName, event.Body.FileName, event.Body.ImageContent); err != nil {
		log.Error().Msgf("Unable to store image to S3: %v", err)
		return returnError(500, customErrors.UnableToStoreImageToS3), nil
	}

	return &events.APIGatewayProxyResponse{StatusCode: 200, Body: response.SuccessfulResponse}, nil
}

func saveToS3(ctx context.Context, s3Client *s3.Client, bucketName, objectKey, base64EncodedImage string) error {
	reader := strings.NewReader(base64EncodedImage)

	input := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	output, err := s3Client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return err
	}

	var completedParts []types.CompletedPart

	partNumber := 1
	buffer := make([]byte, 5*1024*1024)

	for {
		// Some shitty logic to break file in parts and then store those parts to S3.
		// This is done because S3 was giving request too big error. Now image will be uploaded in parts and then merged.
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		partInput := &s3.UploadPartInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(objectKey),
			UploadId:   output.UploadId,
			Body:       bytes.NewReader(buffer[:n]),
			PartNumber: aws.Int32(int32(partNumber)),
		}

		partOutput, err := s3Client.UploadPart(ctx, partInput)
		if err != nil {
			return err
		}

		// An array to hold all the completed parts.
		completedParts = append(completedParts, types.CompletedPart{
			PartNumber: aws.Int32(int32(partNumber)),
			ETag:       partOutput.ETag,
		})

		partNumber++

	}

	_, err = s3Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(bucketName),
		Key:             aws.String(objectKey),
		UploadId:        output.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{Parts: completedParts},
	})
	if err != nil {
		return err
	}
	return nil
}

// We were handling error several times in our code. So we created this function to make our process easy.
func returnError(statusCode int, body string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}
