package main

import (
	"fmt"
	"reflect"
)

type Filler interface {
	RegisterType(defaultValue interface{})
	Fill(dataPtr interface{}) error
}

type filler struct {
	types map[string]reflect.Value
}

func New() Filler {
	return &filler{
		types: make(map[string]reflect.Value),
	}
}

func (f *filler) RegisterType(defaultValue interface{}) {
	value := reflect.ValueOf(defaultValue)
	f.types[value.Type().String()] = value
}

func (f *filler) Fill(dataPtr interface{}) (err error) {
	ptrVal := reflect.ValueOf(dataPtr)
	func() {
		// trying to get value from pointer
		defer func() {
			r := recover()
			if r != nil {
				err = fmt.Errorf("data type %q is not a pointer", ptrVal.Type().String())
			}
		}()
		ptrVal.Elem()
	}()
	if err != nil {
		return err
	}

	return f.fillRecursive(ptrVal)
}

func (f *filler) fillRecursive(value reflect.Value) error {
	// set if it's basic type
	defaultValue, ok := f.types[value.Elem().Type().String()]
	if ok {
		value.Elem().Set(defaultValue)
		return nil
	}

	// fill recursive if it's container type
	switch value.Elem().Type().Kind() {
	case reflect.Struct:
		for i := 0; i < value.Elem().NumField(); i++ {
			if err := f.fillRecursive(value.Elem().Field(i).Addr()); err != nil {
				return err
			}
		}
	case reflect.Slice:
		sliceEl := reflect.New(value.Elem().Type().Elem())
		if err := f.fillRecursive(sliceEl); err != nil {
			return err
		}
		value.Elem().Set(reflect.Append(value.Elem(), sliceEl.Elem()))
	case reflect.Ptr:
		elem := reflect.New(value.Elem().Type().Elem())
		if err := f.fillRecursive(elem); err != nil {
			return err
		}
		value.Elem().Set(elem)
	case reflect.Map:
		mapKey := reflect.New(value.Elem().Type().Key())
		mapElement := reflect.New(value.Elem().Type().Elem())
		if err := f.fillRecursive(mapKey); err != nil {
			return err
		}
		if err := f.fillRecursive(mapElement); err != nil {
			return err
		}
		mapValue := reflect.MakeMap(value.Elem().Type())
		mapValue.SetMapIndex(mapKey.Elem(), mapElement.Elem())
		value.Elem().Set(mapValue)
	default:
		return fmt.Errorf("default value for type %s, %s is not set", value.Elem().Type().String(), value.Elem().Type().Kind().String())
	}

	return nil
}
