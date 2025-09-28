// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import "reflect"

// Metadata represents struct metadata.
type Metadata struct {
	typ    reflect.Type // Struct type (after indirect).
	kind   reflect.Kind // Struct kind.
	fields []*Field     // Struct fields. Nil when the struct has no fields.
}

// NewMetadata retrieves struct metadata by calling [NewTypeMetadata]. Expects
// "v" to be a [reflect.Struct] or [reflect.Ptr] to a [reflect.Struct];
// otherwise, it panics. See [NewTypeMetadata] for details.
func NewMetadata(v any) *Metadata {
	return NewTypeMetadata(reflect.TypeOf(v))
}

// NewTypeMetadata returns a new instance of [Metadata] for the type. Expects
// "typ" to be a [reflect.Struct] or [reflect.Ptr] to a [reflect.Struct];
// otherwise, it panics.
func NewTypeMetadata(typ reflect.Type) *Metadata {
	typ = indirect(typ)
	md := &Metadata{
		typ:  typ,
		kind: typ.Kind(),
	}
	if md.IsStruct() {
		md.getFields()
	}
	return md
}

// Type returns struct type.
func (md *Metadata) Type() reflect.Type { return md.typ }

// Kind returns struct kind.
func (md *Metadata) Kind() reflect.Kind { return md.kind }

// IsStruct returns true if the type is a struct, otherwise false.
func (md *Metadata) IsStruct() bool {
	return indirect(md.typ).Kind() == reflect.Struct
}

// Fields returns structure fields. The slice must be considered as read-only.
func (md *Metadata) Fields() []*Field { return md.fields }

// FieldByName returns a struct field by name or nil if the field doesn't exist.
func (md *Metadata) FieldByName(name string) *Field {
	for _, fld := range md.fields {
		if fld.Name() == name {
			return fld
		}
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
	for i := 0; i < nf; i++ {
		md.fields[i] = NewField(md.typ.Field(i))
	}
}
