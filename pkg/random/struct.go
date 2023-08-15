package random

import (
	"math/rand"
	"reflect"
)

type Min interface {
	Min() int
}

type Max interface {
	Max() int
}

func StructTyped[T any]() T {
	var result T
	Struct(&result)
	return result
}

func Struct(v interface{}) {

	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return
	}

	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		rng := Range{Max: 100}
		if min, ok := field.Addr().Interface().(Min); ok {
			rng.Min = int64(min.Min())
		}
		if max, ok := field.Addr().Interface().(Max); ok {
			rng.Max = int64(max.Max())
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(rng.String())
		case reflect.Int:
			field.SetInt(rng.Int64())
		case reflect.Int8:
			field.SetInt(rng.Int64())
		case reflect.Int16:
			field.SetInt(rng.Int64())
		case reflect.Int32:
			field.SetInt(rng.Int64())
		case reflect.Int64:
			field.SetInt(rng.Int64())
		case reflect.Uint:
			field.SetUint(rng.Uint64())
		case reflect.Uint8:
			field.SetUint(rng.Uint64())
		case reflect.Uint16:
			field.SetUint(rng.Uint64())
		case reflect.Uint32:
			field.SetUint(rng.Uint64())
		case reflect.Uint64:
			field.SetUint(rng.Uint64())
		case reflect.Float32:
			field.SetFloat(rng.Float64())
		case reflect.Float64:
			field.SetFloat(rng.Float64())
		case reflect.Bool:
			field.SetBool(rand.Intn(2) == 0)
		}
	}
}
