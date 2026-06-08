// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

import (
	"reflect"
	"runtime"
)

// Metadata represents struct metadata.
type Metadata struct {
	typ      reflect.Type   // Struct type (after indirect).
	kind     reflect.Kind   // Struct kind.
	fields   []*Field       // Struct fields. Nil when the struct has no fields.
	fieldIdx map[string]int // Field name to fields slice index.
	pkg      string         // Type import string may be empty.
	name     string         // Type name when, may be empty.
}

// NewMetadata extracts [Metadata] about type of "v". Panics for nil value.
func NewMetadata(v any) *Metadata {
	return NewTypeMetadata(reflect.TypeOf(v))
}

// NewTypeMetadata extracts [Metadata] for the type. Panics when the type
// represents nil value.
func NewTypeMetadata(typ reflect.Type) *Metadata {
	typ = indirect(typ)
	md := &Metadata{
		typ:  typ,
		kind: typ.Kind(),
		pkg:  typ.PkgPath(),
		name: typ.Name(),
	}
	if md.IsStruct() {
		md.getFields()
	}
	return md
}

// NewValueMetadata extracts [Metadata] about type of "v".
func NewValueMetadata(val reflect.Value) *Metadata {
	typ := val.Type()
	md := NewTypeMetadata(typ)
	if md.kind == reflect.Func && md.name == "" {
		if val.IsValid() && val.Pointer() != 0 {
			if fn := runtime.FuncForPC(val.Pointer()); fn != nil {
				name := fn.Name()
				md.pkg, md.name = funcPkg(name)
			}
		}
	}
	return md
}

// Type returns struct type.
func (md *Metadata) Type() reflect.Type { return md.typ }

// Kind returns the reflect.Kind of the underlying type.
func (md *Metadata) Kind() reflect.Kind { return md.kind }

// Package returns import string for the type. May return an empty string.
func (md *Metadata) Package() string { return md.pkg }

// Name returns type name. May return an empty string.
func (md *Metadata) Name() string { return md.name }

// PackageAndName returns full type name including package.
func (md *Metadata) PackageAndName() string { return md.pkg + "." + md.name }

// IsStruct returns true if the type is a struct, otherwise false.
func (md *Metadata) IsStruct() bool {
	return indirect(md.typ).Kind() == reflect.Struct
}

// Fields returns structure fields. The slice must be considered as read-only.
func (md *Metadata) Fields() []*Field { return md.fields }

// FieldByName returns a struct field by name or nil if the field doesn't exist.
func (md *Metadata) FieldByName(name string) *Field {
	if i, ok := md.fieldIdx[name]; ok {
		return md.fields[i]
	}
	return nil
}

// FieldByIndex returns the field at the specified index in the struct. If the
// index is out of range, it returns nil.
func (md *Metadata) FieldByIndex(idx int) *Field {
	if idx >= len(md.fields) {
		return nil
	}
	return md.fields[idx]
}

// getFields gets all struct fields.
func (md *Metadata) getFields() {
	nf := md.typ.NumField()
	if nf == 0 {
		return
	}
	md.fields = make([]*Field, nf)
	md.fieldIdx = make(map[string]int, nf)
	for i := 0; i < nf; i++ {
		md.fields[i] = NewField(md.typ.Field(i))
		md.fieldIdx[md.fields[i].Name()] = i
	}
}
