package automapper

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type AutoMapper[F any, T any] interface {
	MapTo(from *F, to *T) error
	MapFrom(to *T, from *F) error
}

func NewAutoMapper[F any, T any](from F, to T) AutoMapper[F, T] {
	return &autoMapper[F, T]{}
}

type autoMapper[F any, T any] struct {
}

func (m *autoMapper[F, T]) MapTo(src *F, dst *T) error {
	return m.doMap(src, dst)
}

func (m *autoMapper[F, T]) MapFrom(src *T, dst *F) error {
	return m.doMap(src, dst)
}

func (m *autoMapper[F, T]) doMap(src any, dst any) error {
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	if srcValue.Kind() == reflect.Pointer {
		srcValue = srcValue.Elem()
	}
	if dstValue.Kind() == reflect.Pointer {
		dstValue = dstValue.Elem()
	}

	if srcValue.Kind() == reflect.Struct && dstValue.Kind() == reflect.Struct {
		return m.mapStruct(src, dst)
	} else if srcValue.CanConvert(dstValue.Type()) {
		srcJson, _ := json.Marshal(src)
		_ = json.Unmarshal(srcJson, dst)
	}
	return nil
}

func (m *autoMapper[F, T]) mapStruct(src any, dst any) error {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	var err error

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcFieldType := srcVal.Type().Field(i)

		if srcField.Kind() == reflect.Struct {
			srcFieldVal := reflect.ValueOf(srcField.Interface())
			for j := 0; j < srcFieldVal.NumField(); j++ {
				nestedSrcField := srcFieldVal.Field(j)
				nestedSrcFieldType := srcFieldVal.Type().Field(j)
				err = m.mapField(nestedSrcField, dstVal, nestedSrcFieldType.Name)
				if err != nil {
					break
				}
			}
			if err == nil {
				return nil
			}
		}

		err := m.mapField(srcField, dstVal, srcFieldType.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *autoMapper[F, T]) mapField(srcField reflect.Value, dstVal reflect.Value, field string) error {
	dstField := dstVal.FieldByName(field)
	if dstField.Kind() == reflect.Invalid {
		// Ignore fields that are not on dst
		return fmt.Errorf("%v does not exist on dst", field)
	}

	if dstField.IsValid() && dstField.CanSet() {
		// Convert to the correct type, i.e. if src is int and dst is int64
		if srcField.CanConvert(dstField.Type()) {
			srcField = srcField.Convert(dstField.Type())
		} else if srcField.Type() != dstField.Type() {
			return fmt.Errorf(
				"cannot convert src type %v to dst type %v for field %v",
				srcField.Type(),
				dstField.Type(),
				field,
			)
		}

		dstField.Set(srcField)
	}
	return nil
}
