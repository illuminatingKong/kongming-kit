package support

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

type Collection []Value

func NewCollection(items ...interface{}) Collection {
	collection := Collection{}

	for _, item := range items {
		if inputCollection, ok := item.(Collection); ok {
			collection = append(collection, inputCollection...)
			continue
		}

		switch Kind(item) {
		case reflect.Array, reflect.Slice:
			s := reflect.ValueOf(item)
			for i := 0; i < s.Len(); i++ {
				value := s.Index(i).Interface()
				collection = append(collection, NewValue(value))
			}
		default:
			collection = append(collection, NewValue(item))
		}
	}

	return collection
}

func (c Collection) GetE(key string) (Value, error) {
	if key == "" {
		return NewValue(c), nil
	}

	currentKey, rest := splitKey(key)

	// when you request something with an Asterisk, you always develop a collection
	if currentKey == "*" {

		flattenCollection := Collection{}
		flattenMap := Map{}

		for _, value := range c {
			switch Kind(value.Source()) {
			case reflect.Slice, reflect.Array:
				flattenCollection = append(flattenCollection, value.Source().(Collection)...)
			case reflect.Map:
				flattenMap = value.Source().(Map).Merge(flattenMap)
			default:
				return NewValue(c), nil
			}
		}

		if len(flattenMap) > 0 {
			return flattenMap.GetE(joinRest(rest))
		}
		return flattenCollection.GetE(joinRest(rest))
	}

	index, err := strconv.Atoi(currentKey)
	if err != nil {
		return Value{}, errors.Wrap(InvalidCollectionKeyError, fmt.Sprintf("'%s' can only be a number or *", key))
	}

	if len(c) < (index + 1) {
		return Value{}, errors.Wrap(CanNotFoundValueError, fmt.Sprintf("'%s' not found", key))
	}

	return c[index].GetE(joinRest(rest))
}

func (c Collection) Raw() interface{} {
	var result []interface{}
	var raw interface{}

	for _, value := range c {
		raw = value.Raw()
		result = append(result, raw)
	}

	return result
}
