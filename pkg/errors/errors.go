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
)

func ReturnErrorResponse(errorData map[string]string) string {
	errBytes, err := json.Marshal(errorData)
	if err == nil {
		return string(errBytes)
	}
	return "unable to complete this request"
}
