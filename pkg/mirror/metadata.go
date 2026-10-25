// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

import (
	"reflect"
	"runtime"
	"slices"
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
//
// Within one call every occurrence of a struct type shares one [Metadata], so
// a recursive type yields a cyclic graph: [Field.TypeMetadata] of the Next
// field in "type Node struct{ Next *Node }" returns the metadata of Node.
func NewTypeMetadata(typ reflect.Type) *Metadata {
	return newTypeMetadata(typ, nil)
}

// newTypeMetadata extracts [Metadata] for the type. The seen map holds the
// struct metadata built so far in the call, keyed by type; a nil map is
// allocated when the first struct is found.
func newTypeMetadata(
	typ reflect.Type,
	seen map[reflect.Type]*Metadata,
) *Metadata {

	typ = indirect(typ)
	if md, ok := seen[typ]; ok {
		return md
	}
	md := &Metadata{
		typ:  typ,
		kind: typ.Kind(),
		pkg:  typ.PkgPath(),
		name: typ.Name(),
	}
	if md.kind == reflect.Struct {
		if seen == nil {
			seen = make(map[reflect.Type]*Metadata)
		}
		seen[typ] = md // Before the fields, so a cycle finds it.
		md.getFields(seen)
	}
	return md
}

// NewValueMetadata extracts [Metadata] about type of "v".
// Panics when val is the zero [reflect.Value].
func NewValueMetadata(val reflect.Value) *Metadata {
	if !val.IsValid() {
		panic("mirror: NewValueMetadata of invalid value")
	}
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

func (md *Metadata) Type() reflect.Type { return md.typ }

func (md *Metadata) Kind() reflect.Kind { return md.kind }

// Package returns import string for the type. May return an empty string.
func (md *Metadata) Package() string { return md.pkg }

// Name returns type name. May return an empty string.
func (md *Metadata) Name() string { return md.name }

// PackageAndName returns the package path and the type name joined by a
// dot. An empty part is omitted.
func (md *Metadata) PackageAndName() string {
	if md.pkg == "" {
		return md.name
	}
	if md.name == "" {
		return md.pkg
	}
	return md.pkg + "." + md.name
}

func (md *Metadata) IsStruct() bool {
	return md.kind == reflect.Struct
}

// Fields returns a copy of the structure fields.
func (md *Metadata) Fields() []*Field { return slices.Clone(md.fields) }

// NumField returns the number of struct fields, zero for a type that is not a
// struct. Unlike [Metadata.Fields] it copies nothing, so a loop over
// [Metadata.FieldByIndex] reads the fields without allocating.
func (md *Metadata) NumField() int { return len(md.fields) }

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
	if idx < 0 || idx >= len(md.fields) {
		return nil
	}
	return md.fields[idx]
}

// getFields gets all struct fields. The seen map is passed to
// newTypeMetadata for each field type.
func (md *Metadata) getFields(seen map[reflect.Type]*Metadata) {
	nf := md.typ.NumField()
	if nf == 0 {
		return
	}
	md.fields = make([]*Field, nf)
	md.fieldIdx = make(map[string]int, nf)
	for i := 0; i < nf; i++ {
		md.fields[i] = newField(md.typ.Field(i), seen)
		md.fieldIdx[md.fields[i].Name()] = i
	}
}
