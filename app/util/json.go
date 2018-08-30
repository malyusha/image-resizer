package util

import (
	"bytes"
	"encoding/json"
)

// jsonPretty returns given interface as json string
func JsonPretty(v interface{}) string {
	var out bytes.Buffer
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	if err := json.Indent(&out, b, "  ", "  "); err != nil {
		panic(err)
	}

	return out.String()
}
