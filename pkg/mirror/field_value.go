// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

import (
	"fmt"
	"reflect"
)

type field = Field // Do not expose the embedded struct.

// FieldValue is a wrapper for a struct field.
type FieldValue struct {
	*field
	value reflect.Value
}

// NewFieldValue returns a new instance of [FieldValue].
func NewFieldValue(fld *Field, value reflect.Value) *FieldValue {
	return &FieldValue{field: fld, value: value}
}

// StructValue returns the field as [StructValue].
func (fv *FieldValue) StructValue() *StructValue {
	return &StructValue{
		metadata: ReflectType(fv.Type()),
		value:    fv.value,
		kind:     fv.kind,
	}
}

func (fv *FieldValue) Field() *Field { return fv.field }

func (fv *FieldValue) Value() reflect.Value { return fv.value }

// NewIfNil allocates a nil pointer, map, or slice.
// An unsettable field is left unchanged. A pointer becomes a pointer to
// the zero element. A map or slice becomes empty. Other kinds stay nil.
func (fv *FieldValue) NewIfNil() *FieldValue {
	v := fv.value
	if !v.IsValid() || !v.CanSet() {
		return fv
	}
	switch fv.kind {
	case reflect.Pointer:
		if v.IsNil() {
			v.Set(reflect.New(fv.Type().Elem()))
		}

	case reflect.Map:
		if v.IsNil() {
			kt := fv.Type().Key()
			vt := fv.Type().Elem()
			mt := reflect.MapOf(kt, vt)
			m := reflect.MakeMapWithSize(mt, 1)
			v.Set(m)
		}

	case reflect.Slice:
		if v.IsNil() {
			v.Set(reflect.MakeSlice(fv.Type(), 0, 0))
		}

	default:
		// Not a pointer, map, or slice.
	}
	return fv
}

// Get gets the field value or error if the field is invalid.
func (fv *FieldValue) Get() (any, error) {
	if !fv.IsValid() {
		return nil, ErrInvField
	}
	if !fv.IsExported() {
		return nil, fmt.Errorf("%w: %s", ErrUnexportedField, fv.Name())
	}
	return fv.value.Interface(), nil
}
