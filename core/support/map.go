package support

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"reflect"
)

type Map map[string]Value

func NewMapE(itemsRange ...interface{}) (Map, error) {
	var err error
	result := Map{}

	for _, rawItems := range itemsRange {
		v := reflect.ValueOf(rawItems)
		if v.Kind() != reflect.Map {
			err = errors.WithStack(errors.Wrap(CanNotCreateMapError, fmt.Sprintf("type %s", v.Kind().String())))
			continue
		}

		for _, key := range v.MapKeys() {
			value := v.MapIndex(key).Interface()
			key, err := cast.ToStringE(key.Interface())
			if err != nil {
				return nil, errors.WithStack(errors.Wrap(err, "invalid key in map"))
			}
			result[key] = NewValue(value)
		}
	}

	return result, err
}

func (m Map) Merge(maps ...Map) Map {
	for _, bag := range maps {
		for key, item := range bag {
			m.Push(key, item)
		}
	}

	return m
}

func (m Map) Push(key string, input interface{}) Map {
	if rawValue, found := m[key]; found {
		source := rawValue.Source()
		switch source.(type) {
		case Collection:
			collection := source.(Collection)
			m[key] = NewValue(collection.Push(input))
		default:
			m[key] = NewValue(input)
		}
	} else {
		m[key] = NewValue(input)
	}

	return m
}

func (v Value) Source() interface{} {
	return v.source
}

func (c Collection) Push(item interface{}) Collection {
	return append(c, NewValue(item))
}

func (m Map) GetE(key string) (Value, error) {
	if key == "" {
		return NewValue(m), nil
	}

	currentKey, rest := splitKey(key)

	// when you request something with an asterisk, you always develop a collection
	if currentKey == "*" {
		collection := Collection{}
		for _, value := range m {
			nestedValueRaw, err := value.GetE(joinRest(rest))
			if err != nil {
				return nestedValueRaw, err
			}
			switch nestedValues := nestedValueRaw.source.(type) {
			case Collection:
				for _, nestedValue := range nestedValues {
					collection = collection.Push(nestedValue)
				}
			case Map:
				for _, nestedValue := range nestedValues {
					collection = collection.Push(nestedValue)
				}
			default:
				// If there are no keys to search further, the nested value is the final value
				collection = collection.Push(nestedValueRaw)
			}
		}

		return NewValue(collection), nil
	}

	value, found := m[key]
	if found {
		return value, nil
	}
	value, found = m[currentKey]
	if !found {
		return Value{}, errors.Wrap(CanNotFoundValueError, fmt.Sprintf("key '%s'%s", currentKey, getKeyInfo(key, currentKey)))
	}

	switch value.Source().(type) {
	case Collection:
		return value.Collection().GetE(joinRest(rest))
	case Map:
		return value.Map().GetE(joinRest(rest))
	default:
		return value.GetE(joinRest(rest))
	}
}

func getKeyInfo(key string, currentKey string) string {
	info := ""
	if currentKey != key {
		info = " ('" + key + "')"
	}
	return info
}

func (m Map) Collection() Collection {
	collection := Collection{}
	for _, value := range m {
		collection = collection.Push(value)
	}

	return collection
}

func (m Map) Raw() interface{} {
	result := map[string]interface{}{}

	for key, value := range m {
		// Handle value
		result[key] = value.Raw()
	}

	return result
}
