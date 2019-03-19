package dot

import (
	"errors"
	"reflect"
	"strings"
)

type Map struct {
	values map[string]interface{}
}

// Get returns *Value with underlying from values property. If it's not there default underlying or nil will be
// returned
func (c *Map) Get(key string, defValue ...interface{}) *Value {
	return dotGetProp(c.values, key, defValue...)
}

// Raw returns raw data representation inside Map
func (c *Map) Raw() map[string]interface{} {
	return c.values
}

// Hash checks whether config has a underlying with given key in Service property
func (c *Map) Has(key string) bool {
	if !strings.Contains(key, ".") {
		_, ok := c.values[key]

		return ok
	}

	val := dotGetProp(c.values, key).underlying

	return val != nil
}

// NewMap returns new initialized Map struct pointer with given values
func NewMap(values map[string]interface{}) *Map {
	return &Map{values: values}
}

// dotGetProp returns Value with underlying received from obj by given key.
// It splits key by dot "." and look to nested structure of obj.
// WARN: obj must be map
func dotGetProp(obj interface{}, key string, defValue ...interface{}) *Value {
	var (
		v   interface{}
		err error
	)

	hasDefaultValue := len(defValue) > 0

	arr := strings.Split(key, ".")
	last := len(arr) - 1

	for ix, key := range arr {
		obj, err = getProp(obj, key)
		if err != nil || obj == nil {
			if hasDefaultValue {
				v = defValue[0]
				break
			}

			break
		}

		if ix == last {
			v = obj
		}
	}

	if v == nil {
		return &Value{}
	}

	return &Value{Key: key, underlying: v}
}

// getProp returns interface{} from obj property. It takes returns nil for underlying if it's not valid.
// Only map should be passed as obj
func getProp(obj interface{}, prop string) (interface{}, error) {
	if reflect.TypeOf(obj).Kind() != reflect.Map {
		return nil, errors.New("only map should be passed to getProp func")
	}

	val := reflect.ValueOf(obj)

	value := val.MapIndex(reflect.ValueOf(prop))

	if !value.IsValid() {
		return nil, nil
	}

	return value.Interface(), nil
}
