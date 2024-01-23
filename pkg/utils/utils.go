/*
	@author: Sushant
	@last-modified: 23 January 2024
	@GitHub: https://github.com/sushant102004
*/

package utils

import "encoding/json"

/*
	AWS Lambda events.APIGatewayProxyEvent returns string as response body. So here we have two helper functions that will
	convert and return map[string]string as string.
*/

func ReturnErrorResponse(errorData map[string]string) string {
	errBytes, err := json.Marshal(errorData)
	if err != nil {
		return "unable to complete this request"
	}
	return string(errBytes)
}

func ReturnResponse(responseData map[string]string) string {
	respBytes, err := json.Marshal(responseData)
	if err != nil {
		return "invalid response data"
	}
	return string(respBytes)
}
