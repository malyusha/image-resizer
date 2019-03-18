package config

// Info is the structure for additional information of configuration
// It includes Key and value and can cast this value into requested type
type info struct {
	Key   string
	value interface{}
}

// Bool returns boolean representation of value
func (i *info) Bool() bool {
	v, ok := i.value.(bool)

	if !ok {
		return false
	}

	return v
}

// Float32 returns float32 representation of value
func (i *info) Float32() float32 {
	switch i.value.(type) {
	case float32:
		return i.value.(float32)
	case float64:
		return float32(i.value.(float64))
	}

	return 0
}

// Float64 returns float64 representation of value
func (i *info) Float64() float64 {
	switch i.value.(type) {
	case float32:
		return float64(i.value.(float32))
	case float64:
		return i.value.(float64)
	}

	return 0
}

// Int returns integer representation of value
func (i *info) Int() int {
	switch i.value.(type) {
	case int:
		return i.value.(int)
	case float64:
		return int(i.value.(float64))
	case float32:
		return int(i.value.(float32))
	}

	return 0
}

// String returns string representation of value
func (i *info) String() string {
	v, ok := i.value.(string)

	if !ok {
		return ""
	}

	return v
}
