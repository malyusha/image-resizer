package dot

import (
	"fmt"
)

// Info is the structure for additional information of configuration
// It includes Key and underlying and can cast this underlying into requested type
type Value struct {
	Key        string
	underlying interface{}
}

// Bool returns boolean representation of underlying
func (i *Value) Bool() bool {
	v, ok := i.underlying.(bool)

	if !ok {
		return false
	}

	return v
}

// Float32 returns float32 representation of underlying
func (i *Value) Float32() float32 {
	switch i.underlying.(type) {
	case float32:
		return i.underlying.(float32)
	case float64:
		return float32(i.underlying.(float64))
	}

	return 0
}

// Float64 returns float64 representation of underlying
func (i *Value) Float64() float64 {
	switch i.underlying.(type) {
	case float32:
		return float64(i.underlying.(float32))
	case float64:
		return i.underlying.(float64)
	}

	return 0
}

// Int returns integer representation of underlying
func (i *Value) Int() int {
	switch i.underlying.(type) {
	case int:
		return i.underlying.(int)
	case float64:
		return int(i.underlying.(float64))
	case float32:
		return int(i.underlying.(float32))
	}

	return 0
}

// String returns string representation of underlying
func (i *Value) String() string {
	v, ok := i.underlying.(string)

	if !ok {
		return ""
	}

	return v
}

// Map returns dynamic config for underlying
func (i *Value) Map() *Map {
	_, ok := i.underlying.(map[string]interface{})

	if ok {
		return &Map{i.underlying.(map[string]interface{})}
	}

	if v, ok := i.underlying.(map[interface{}]interface{}); ok {
		return &Map{cleanUpInterfaceMap(v)}
	}

	return &Map{}
}

func cleanUpInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case map[interface{}]interface{}:
		return cleanUpInterfaceMap(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
