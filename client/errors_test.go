package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUnexpectedResponseError(t *testing.T) {
	assert.Error(t, NewUnexpectedResponseError(400), UNEXPECTED_RESPONSE_FORMAT, 400)
	assert.Error(t, NewUnexpectedResponseError(400, 200), UNEXPECTED_RESPONSE_FORMAT, 400)
	assert.Error(t, NewUnexpectedResponseError(400, 200, 300), UNEXPECTED_RESPONSE_FORMAT, 400)
	assert.Nil(t, NewUnexpectedResponseError(400, 200, 300, 400))
}
