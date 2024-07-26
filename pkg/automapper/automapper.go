package automapper

import (
	"encoding/json"
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
		m.mapStruct(src, dst)
	} else if srcValue.CanConvert(dstValue.Type()) {
		srcJson, _ := json.Marshal(src)
		_ = json.Unmarshal(srcJson, dst)
	}
	return nil
}

func (m *autoMapper[F, T]) mapStruct(src any, dst any) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcFieldType := srcVal.Type().Field(i)

		dstField := dstVal.FieldByName(srcFieldType.Name)
		if dstField.Kind() == reflect.Invalid {
			// Ignore fields that are not on dst
			continue
		}

		// if srcField.Kind() == reflect.Struct && dstField.Kind() == reflect.Struct {
		// 	srcFieldVal := srcField.Interface()
		// 	dstFieldVal := dstField.Interface()
		// 	fieldMapper := CreateAutoMapper(srcFieldVal, dstFieldVal)
		// 	fieldMapper.MapTo(&srcFieldVal, &dstFieldVal)
		// 	continue
		// }

		if dstField.IsValid() && dstField.CanSet() {
			// Convert to the correct type, i.e. if src is int and dst is int64
			if srcField.CanConvert(dstField.Type()) {
				srcField = srcField.Convert(dstField.Type())
			}

			dstField.Set(srcField)
		}
	}
}
