package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	defaultMsg = "Error"
)

func Test_messOrDefaultReturnDefaultError(t *testing.T) {
	assert.Equal(t, defaultMsg, messOrDefault(defaultMsg))
	msg := "An error here"
	assert.Equal(t, msg, messOrDefault(defaultMsg, msg))
}
