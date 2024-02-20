package customErrors

import "encoding/json"

var (
	MethodNotAllowed = ReturnErrorResponse(map[string]string{
		"error": "method not allowed",
	})

	UnableToLoadAWSConfiguration = ReturnErrorResponse(map[string]string{
		"error": "unable to load AWS configuration",
	})

	UnableToFindEnvVariable = ReturnErrorResponse(map[string]string{
		"error": "Unable to find environment variable",
	})

	UnableToStoreImageToS3 = ReturnErrorResponse(map[string]string{
		"error": "Unable to store image to s3",
	})

	InvalidInputBody = ReturnErrorResponse(map[string]string{
		"error": "Invalid input body. There must be some data missing.",
	})

	UnableToEncodeHTTPRequestBody = ReturnErrorResponse(map[string]string{
		"error": "unable to encode http request body",
	})

	UnableToCreateHTTPRequest = ReturnErrorResponse(map[string]string{
		"error": "unable to create HTTP request",
	})

	HTTPResponseError = ReturnErrorResponse(map[string]string{
		"error": "HTTP Request failed. Check logs for more details",
	})

	UnableToReadDataFromHTTPResponse = ReturnErrorResponse(map[string]string{
		"error": "unable to read data from HTTP response",
	})

	UnableToUnmarshalJSON = ReturnErrorResponse(map[string]string{
		"error": "unable to unmarshal json",
	})
)

func ReturnErrorResponse(errorData map[string]string) string {
	errBytes, err := json.Marshal(errorData)
	if err == nil {
		return string(errBytes)
	}
	return "unable to complete this request"
}
