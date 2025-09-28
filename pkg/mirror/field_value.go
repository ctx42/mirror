// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
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

// Field returns the [Field] associated with the [FieldValue].
func (fv *FieldValue) Field() *Field { return fv.field }

// Value returns the [reflect.Value] associated with the [FieldValue].
func (fv *FieldValue) Value() reflect.Value { return fv.value }

// NewIfNil initializes the value of a field if it is nil with its zero value.
func (fv *FieldValue) NewIfNil() *FieldValue {
	v := fv.value
	switch fv.kind {
	case reflect.Ptr:
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
		// Not a pointer.
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
