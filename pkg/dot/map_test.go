package dot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_Get(t *testing.T) {
	testMap := map[string]interface{}{
		"first": map[string]string{
			"nested": "underlying",
		},
		"second": map[string]interface{}{
			"nested": map[string]interface{}{
				"more":   "levels",
				"number": 150,
				"float":  1.5,
			},
		},
		"third": true,
		"fourth": map[string]bool{
			"enabled": false,
		},
	}

	dynamicCfg := Map{values: testMap}

	assert.Equal(t, "underlying", dynamicCfg.Get("first.nested").String())
	assert.Equal(t, "levels", dynamicCfg.Get("second.nested.more").String())
	assert.Equal(t, 150, dynamicCfg.Get("second.nested.number").Int())
	assert.Equal(t, float64(1.5), dynamicCfg.Get("second.nested.float").Float64())
	assert.Equal(t, "", dynamicCfg.Get("not_existing.level.of.map").String())
	assert.Equal(t, "default_value", dynamicCfg.Get("not_existing.level.of.map", "default_value").String())
}

func TestMap_Has(t *testing.T) {
	testMap := map[string]interface{}{
		"exists": map[string]interface{}{
			"yes": true,
		},
	}

	dynamicCfg := Map{values: testMap}

	assert.True(t, dynamicCfg.Has("exists"))
	assert.True(t, dynamicCfg.Has("exists.yes"))
	assert.False(t, dynamicCfg.Has("exists.no"))
	assert.False(t, dynamicCfg.Has("not_exist_at_all.some.nested.prop"))
}
