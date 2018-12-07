package model

import "encoding/json"

// GenericError base error message
type GenericError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Object  string `json:"object"`
}

// ************************************************************************ Functions

// GenerateError function to create base error message
func GenerateError(m string, c string, o string) GenericError {
	genericError := GenericError{m, c, o}
	return genericError
}

// GenerateErrorString function to create base error message
func GenerateErrorString(m string, c string, o string) ([]byte, error) {
	genericError := GenerateError(m, c, o)
	ge, err := json.Marshal(genericError)
	return ge, err
}
