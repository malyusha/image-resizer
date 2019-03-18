package config

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo_Bool(t *testing.T) {
	correctInfo := info{Key: "bool_test", value: true}
	invalidInfo := info{Key: "bool_test_invalid", value: "not boolean"}

	assert.True(t, correctInfo.Bool())
	assert.False(t, invalidInfo.Bool())
}

func TestInfo_String(t *testing.T) {
	correctInfo := info{Key: "string_test", value: "hello world"}
	invalidInfo := info{Key: "string_test_invalid", value: 15}

	assert.Equal(t, "hello world", correctInfo.String())
	assert.Equal(t, "", invalidInfo.String())
}

func TestInfo_Int(t *testing.T) {
	correctInfo := info{Key: "int_test", value: 100500}
	invalidInfo := info{Key: "int_test_invalid", value: "not integer"}

	assert.Equal(t, 100500, correctInfo.Int())
	assert.Equal(t, 0, invalidInfo.Int())
}

func TestInfo_Float32(t *testing.T) {
	correctInfo := info{Key: "float32_test", value: float32(math.MaxFloat32)}
	invalidInfo := info{Key: "float32_test_invalid", value: "not integer"}

	assert.Equal(t, float32(math.MaxFloat32), correctInfo.Float32())
	assert.Equal(t, float32(0), invalidInfo.Float32())
}

func TestInfo_Float64(t *testing.T) {
	correctInfo := info{Key: "float64_test", value: float64(math.MaxFloat64)}
	invalidInfo := info{Key: "float64_test_invalid", value: int(50)}

	assert.Equal(t, float64(math.MaxFloat64), correctInfo.Float64())
	assert.Equal(t, float64(0), invalidInfo.Float64())
}
