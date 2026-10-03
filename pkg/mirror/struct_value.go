// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

import (
	"reflect"
)

type metadata = Metadata // Do not expose the embedded struct.

// StructValue is a wrapper for a concrete struct instance.
type StructValue struct {
	*metadata               // The struct metadata.
	value     reflect.Value // Indirection of the struct s.
	kind      reflect.Kind  // Value kind.
}

// NewStructValue wraps a struct pointer of any type and provides an easy
// reflection interface. The argument `s` must be a pointer to a struct
// otherwise the function returns nil.
func NewStructValue(s any) *StructValue {
	value := reflect.ValueOf(s)
	sv := &StructValue{
		value: value,
		kind:  value.Kind(),
	}

	if sv.kind != reflect.Ptr {
		return nil
	}

	val := reflect.Indirect(sv.value)
	if val.Kind() != reflect.Struct {
		return nil
	}
	sv.metadata = Reflect(s)
	return sv
}

// IsPtr returns true if the struct value is a pointer type.
func (sv *StructValue) IsPtr() bool { return sv.kind == reflect.Ptr }

// IsValid returns true when the value field is valid.
func (sv *StructValue) IsValid() bool {
	if sv == nil {
		return false
	}
	return sv.value.IsValid()
}

// Metadata returns metadata for the struct type.
func (sv *StructValue) Metadata() *Metadata { return sv.metadata }

// NumField returns the number of fields in the structure.
// A value that is not a struct has no fields.
func (sv *StructValue) NumField() int {
	if sv.metadata == nil || !sv.IsStruct() {
		return 0
	}
	return sv.Type().NumField()
}

// FieldByName returns a struct field or nil if the field does not exist.
// A nil struct pointer returns nil.
func (sv *StructValue) FieldByName(name string) *FieldValue {
	if sv.metadata == nil {
		return nil
	}
	fld := sv.metadata.FieldByName(name)
	if fld == nil || len(fld.index) == 0 {
		return nil
	}
	val := sv.deref()
	if !val.IsValid() {
		return nil
	}
	if val = val.Field(fld.index[0]); val.IsValid() {
		return NewFieldValue(fld, val)
	}
	return nil
}

// FieldByIndex returns a struct field or nil if the field does not exist.
// A nil struct pointer returns nil.
func (sv *StructValue) FieldByIndex(idx int) *FieldValue {
	if sv.metadata == nil {
		return nil
	}
	fld := sv.metadata.FieldByIndex(idx)
	if fld == nil {
		return nil
	}
	val := sv.deref()
	if !val.IsValid() {
		return nil
	}
	if val = val.Field(idx); val.IsValid() {
		return NewFieldValue(fld, val)
	}
	return nil
}

// deref returns the struct sv holds.
// A nil pointer or a non-struct produces the zero Value.
func (sv *StructValue) deref() reflect.Value {
	val := sv.value
	if sv.IsPtr() {
		if !val.IsValid() || val.IsNil() {
			return reflect.Value{}
		}
		val = val.Elem()
	}
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return reflect.Value{}
	}
	return val
}

// NewIfNil initializes the field value with its zero value if it is nil.
func (sv *StructValue) NewIfNil() *StructValue {
	val := sv.value
	if sv.IsPtr() && val.IsNil() {
		val.Set(reflect.New(sv.Type()))
	}
	return sv
}
