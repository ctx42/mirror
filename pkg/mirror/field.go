// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"reflect"
)

// Field represents struct field.
type Field struct {
	*metadata
	sf         reflect.StructField // Corresponding struct field.
	typ        reflect.Type        // Field type.
	kind       reflect.Kind        // Field type kind.
	anonymous  bool                // Is an embedded field?
	exported   bool                // Is a field exported?
	sliceOfPtr bool                // Is a slice of pointers?
	sliceOrArr bool                // Is slice or array?
	index      []int               // Index sequence for [reflect.Type.FieldByIndex].
	tags       []Tag               // Additional tag options.
}

// NewField returns new instance of struct field.
func NewField(sf reflect.StructField) *Field {
	kind := sf.Type.Kind()
	fld := &Field{
		metadata:   NewTypeMetadata(sf.Type),
		sf:         sf,
		typ:        sf.Type,
		kind:       kind,
		anonymous:  sf.Anonymous,
		exported:   sf.IsExported(),
		sliceOrArr: kind == reflect.Slice || kind == reflect.Array,
		index:      sf.Index,
	}
	fld.tags, _ = ParseTags(fld.sf.Name, string(fld.sf.Tag))
	if fld.sliceOrArr && sf.Type.Elem().Kind() == reflect.Ptr {
		fld.sliceOfPtr = true
	}
	return fld
}

// StructField returns the underlying [reflect.StructField] for the field.
func (fld *Field) StructField() reflect.StructField { return fld.sf }

// Type returns [reflect.StructField.Type] for the field.
func (fld *Field) Type() reflect.Type { return fld.typ }

// Kind returns [reflect.StructField.Type.Kind] for the field.
func (fld *Field) Kind() reflect.Kind { return fld.kind }

// Index returns the index sequence for the field type in a struct.
func (fld *Field) Index() []int { return fld.index }

// Name returns the name of the struct field.
func (fld *Field) Name() string { return fld.sf.Name }

// Tag returns tag by name, if the tag doesn't exist, it returns a tag for
// which the [Tag.IsZero] method returns true.
func (fld *Field) Tag(key string) Tag {
	for _, tag := range fld.tags {
		if tag.key == key {
			return tag
		}
	}
	return Tag{field: fld.sf.Name}
}

// IsValid returns true for fields that may appear in names to struct fields.
func (fld *Field) IsValid() bool { return !fld.IsInterface() && !fld.anonymous }

// IsExported returns true if the name starts with uppercase (i.e. field is
// public).
func (fld *Field) IsExported() bool { return fld.sf.IsExported() }

// IsSliceOfPtr returns true if the field is a slice or array of pointers,
// otherwise false.
func (fld *Field) IsSliceOfPtr() bool { return fld.sliceOfPtr }

// IndirectType if the field type is a pointer type it returns underlying type,
// otherwise it returns the type set in the constructor function.
func (fld *Field) IndirectType() reflect.Type {
	if fld.kind == reflect.Ptr {
		return fld.typ.Elem()
	}
	return fld.typ
}

// IsSlice returns true if the field type is a slice, false otherwise.
func (fld *Field) IsSlice() bool { return fld.kind == reflect.Slice }

// IsArray returns true if the field type is an array, false otherwise.
func (fld *Field) IsArray() bool { return fld.kind == reflect.Array }

// IsSliceOrArray returns true if the field is a slice or an array, false
// otherwise.
func (fld *Field) IsSliceOrArray() bool { return fld.sliceOrArr }

// IsMap returns true if the field type is a map, false otherwise.
func (fld *Field) IsMap() bool { return fld.kind == reflect.Map }

// IsInterface returns true if the field is an interface, false otherwise.
func (fld *Field) IsInterface() bool { return fld.kind == reflect.Interface }

// TypeMetadata return [Metadata] for the type.
func (fld *Field) TypeMetadata() *Metadata { return fld.metadata }

// IsAnonymous returns true for embedded fields, false otherwise.
func (fld *Field) IsAnonymous() bool { return fld.anonymous }
