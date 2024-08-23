package support

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

type Value struct {
	source interface{}
}

type nonValueError struct {
	error
}

func NewValue(val interface{}, preErr ...error) Value {
	if len(preErr) > 0 {
		val = nonValueError{preErr[0]}
	}

	switch val.(type) {
	case []byte:
		return Value{val}
	case Collection, Map:
		return Value{val}
	case Value:
		return val.(Value)
	}

	switch Kind(val) {
	case reflect.Slice, reflect.Array:
		result := NewCollection(val)
		return Value{result}
	case reflect.Map:
		result, err := NewMapE(val)
		if err != nil {
			val = nonValueError{err}
		}
		val = result
	}
	return Value{val}
}

// A value can contain a collection.
func (v Value) Collection() Collection {
	if err, isPre := v.source.(nonValueError); isPre {
		panic(err)
	}

	switch v.source.(type) {
	case Collection:
		return v.source.(Collection)
	case Map:
		return v.source.(Map).Collection()
	default:
		return NewCollection(v.source)
	}
}

func (v Value) GetE(key string) (Value, error) {
	if key == "" {
		return v, nil
	}

	currentKey, rest := splitKey(key)
	nextKey := joinRest(rest)

	// when you request something with an Asterisk, you always develop a collection
	if currentKey == "*" {

		switch v.source.(type) {
		case Collection:
			return v.source.(Collection).GetE(nextKey)
		case Map:
			return v.source.(Map).GetE(nextKey)
		default:
			return Value{}, errors.New("*: is not a Collection or Map")
		}

	}

	switch source := v.source.(type) {
	case Collection:
		keyInt, err := strconv.Atoi(currentKey)
		if err != nil {
			return Value{}, err
		}
		collection := v.source.(Collection)
		if len(collection) < (keyInt + 1) {
			return Value{}, errors.Wrap(CanNotFoundValueError, fmt.Sprintf("key '%s'%s", currentKey, getKeyInfo(key, currentKey)))
		}
		return collection[keyInt].GetE(nextKey)
	case Map:
		value, ok := v.source.(Map)[currentKey]
		if !ok {
			return value, errors.Wrap(CanNotFoundValueError, fmt.Sprintf("key '%s'%s", currentKey, getKeyInfo(key, currentKey)))
		}
		return value.GetE(nextKey)
	default:
		switch Kind(source) {
		case reflect.Struct:
			val := reflect.ValueOf(source).FieldByName(currentKey)
			if val.IsValid() {
				return NewValue(val.Interface()).GetE(nextKey)
			} else {
				return Value{}, errors.New(currentKey + ": can't find value")
			}

		}
		return Value{}, errors.New(currentKey + ": is not a struct, Collection or Map")
	}
}

func (v Value) Raw() interface{} {
	if result, ok := v.source.(Value); ok {
		return result.Raw()
	}
	if result, ok := v.source.(Collection); ok {
		return result.Raw()
	}
	if result, ok := v.source.(Map); ok {
		return result.Raw()
	}

	return v.source
}

func (v Value) Map() Map {
	result, err := v.MapE()
	if err != nil {
		panic(err)
	}
	return result
}

func (v Value) MapE() (Map, error) {
	source := v.source
	if err, isPre := source.(nonValueError); isPre {
		return Map{}, errors.Unwrap(err)
	}

	switch valueType := source.(type) {
	case Map:
		return source.(Map), nil
	default:
		return nil, errors.New("can't create map from reflect.Kind " + strconv.Itoa(int(Kind(valueType))))
	}
}
