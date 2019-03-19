package dot

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_Bool(t *testing.T) {
	correctInfo := Value{Key: "bool_test", underlying: true}
	invalidInfo := Value{Key: "bool_test_invalid", underlying: "not boolean"}

	assert.True(t, correctInfo.Bool())
	assert.False(t, invalidInfo.Bool())
}

func TestMap_String(t *testing.T) {
	correctInfo := Value{Key: "string_test", underlying: "hello world"}
	invalidInfo := Value{Key: "string_test_invalid", underlying: 15}

	assert.Equal(t, "hello world", correctInfo.String())
	assert.Equal(t, "", invalidInfo.String())
}

func TestMap_Int(t *testing.T) {
	correctInfo := Value{Key: "int_test", underlying: 100500}
	invalidInfo := Value{Key: "int_test_invalid", underlying: "not integer"}

	assert.Equal(t, 100500, correctInfo.Int())
	assert.Equal(t, 0, invalidInfo.Int())
}

func TestMap_Float32(t *testing.T) {
	correctInfo := Value{Key: "float32_test", underlying: float32(math.MaxFloat32)}
	invalidInfo := Value{Key: "float32_test_invalid", underlying: "not integer"}

	assert.Equal(t, float32(math.MaxFloat32), correctInfo.Float32())
	assert.Equal(t, float32(0), invalidInfo.Float32())
}

func TestMap_Float64(t *testing.T) {
	correctInfo := Value{Key: "float64_test", underlying: float64(math.MaxFloat64)}
	invalidInfo := Value{Key: "float64_test_invalid", underlying: int(50)}

	assert.Equal(t, float64(math.MaxFloat64), correctInfo.Float64())
	assert.Equal(t, float64(0), invalidInfo.Float64())
}

func TestMap_Map(t *testing.T) {
	testData := map[string]interface{}{
		"first": map[string]string{
			"nested": "underlying",
		},
		"second": map[string]interface{}{
			"level": map[string]interface{}{
				"enabled": false,
			},
		},
	}

	cfg := Map{values: testData}

	assert.Equal(t, testData["second"], cfg.Get("second").Map().Raw())
	assert.Equal(t, false, cfg.Get("second").Map().Get("level.enabled").Bool())
}
