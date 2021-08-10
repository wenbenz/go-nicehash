package client

import "fmt"

const UNEXPECTED_RESPONSE_FORMAT = "unexpected response code: %d"

func NewUnexpectedResponseError(responseCode int, expectedCodes ...int) error {
	for _, code := range expectedCodes {
		if responseCode == code {
			return nil
		}
	}
	return fmt.Errorf(UNEXPECTED_RESPONSE_FORMAT, responseCode)
}
