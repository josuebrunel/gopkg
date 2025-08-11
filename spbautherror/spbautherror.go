package spbautherror

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Response defines the structure for the JSON part of the error message.
type Response struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Msg       string `json:"msg"`
}

// UnmarshalAuthError parses the specific error string format and extracts the Response.
func Unmarshal(errorString string) (*Response, error) {
	// Find the start of the JSON part
	jsonStart := strings.Index(errorString, "{")
	if jsonStart == -1 {
		return nil, fmt.Errorf("could not find start of JSON object in error string: %s", errorString)
	}

	// Extract the JSON substring
	jsonString := errorString[jsonStart:]

	// Define a variable to hold the unmarshalled data
	var authError Response

	// Unmarshal the JSON string into the struct
	err := json.Unmarshal([]byte(jsonString), &authError)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from error string: %w", err)
	}

	return &authError, nil
}
