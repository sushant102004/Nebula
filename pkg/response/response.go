package response

import (
	"encoding/json"
)

var (
	SuccessfulResponse = ReturnResponse(map[string]string{
		"message": "request completed successfully",
	})
)

func ReturnResponse(responseData map[string]string) string {
	respBytes, err := json.Marshal(responseData)
	if err == nil {
		return string(respBytes)
	}
	return "invalid response data"
}
