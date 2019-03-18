package errors

import (
	"strings"
)

const (
	NotFoundMessage        = "Not found"
	InternalErrorMessage   = "Internal error"
)

func messOrDefault(defaultMsg string, customMsg ...string) string {
	if customMsg == nil {
		return defaultMsg
	}

	return strings.Join(customMsg, ".")
}
