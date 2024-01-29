package customErrors

import "encoding/json"

var (
	MethodNotAllowed = ReturnErrorResponse(map[string]string{
		"error": "method not allowed",
	})

	UnableToLoadAWSConfiguration = ReturnErrorResponse(map[string]string{
		"error": "unable to load AWS configuration",
	})

	UnableToConnectToS3 = ReturnErrorResponse(map[string]string{
		"error": "Unable to connect to S3",
	})

	UnableToLoadEnv = ReturnErrorResponse(map[string]string{
		"error": "Unable to load enviroment variables",
	})
	UnableToFindEnvVariable = ReturnErrorResponse(map[string]string{
		"error": "Unable to find environment variable",
	})

	UnableToUnmarshalJSON = ReturnErrorResponse(map[string]string{
		"error": "Unable to unmarshal JSON",
	})

	UnableToConvertBase64StringToImage = ReturnErrorResponse(map[string]string{
		"error": "Unable to convert base64 string to image",
	})

	UnableToStoreImageToS3 = ReturnErrorResponse(map[string]string{
		"error": "Unable to store image to s3",
	})

	UnableToFindBucketNameAnywhere = ReturnErrorResponse(map[string]string{
		"error": "Bucket name must be defined into request body of either enviroment variables",
	})

	InvalidInputBody = ReturnErrorResponse(map[string]string{
		"error": "Invalid input body. There must be some data missing.",
	})
)

func ReturnErrorResponse(errorData map[string]string) string {
	errBytes, err := json.Marshal(errorData)
	if err == nil {
		return string(errBytes)
	}
	return "unable to complete this request"
}
