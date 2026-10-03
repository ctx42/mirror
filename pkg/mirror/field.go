// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

import (
	"reflect"
	"slices"
)

// Field represents a struct field.
type Field struct {
	*metadata
	sf         reflect.StructField // Corresponding struct field.
	typ        reflect.Type        // Field type.
	kind       reflect.Kind        // Field type kind.
	anonymous  bool                // Is an embedded field?
	sliceOfPtr bool                // Is a slice of pointers?
	sliceOrArr bool                // Is slice or array?

	index []int // Index sequence for [reflect.Type.FieldByIndex].
	tags  []Tag // Additional tag options.
}

// NewField returns a new instance of the struct field.
func NewField(sf reflect.StructField) *Field {
	kind := sf.Type.Kind()
	sf.Index = slices.Clone(sf.Index)
	fld := &Field{
		metadata:   NewTypeMetadata(sf.Type),
		sf:         sf,
		typ:        sf.Type,
		kind:       kind,
		anonymous:  sf.Anonymous,
		sliceOrArr: kind == reflect.Slice || kind == reflect.Array,
		index:      sf.Index,
	}
	fld.tags, _ = ParseTags(fld.sf.Name, string(fld.sf.Tag))
	if fld.sliceOrArr && sf.Type.Elem().Kind() == reflect.Pointer {
		fld.sliceOfPtr = true
	}
	return fld
}

// StructField returns the underlying [reflect.StructField] for the field.
// The index sequence is a copy.
func (fld *Field) StructField() reflect.StructField {
	sf := fld.sf
	sf.Index = slices.Clone(fld.index)
	return sf
}

func (fld *Field) Type() reflect.Type { return fld.typ }

func (fld *Field) Kind() reflect.Kind { return fld.kind }

// Index returns a copy of the index sequence for the field.
func (fld *Field) Index() []int { return slices.Clone(fld.index) }

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

// IsValid returns false for interface fields and anonymous (embedded) fields;
// true for all others.
func (fld *Field) IsValid() bool { return !fld.IsInterface() && !fld.anonymous }

// IsExported returns true if the name starts with uppercase (i.e. field is
// public).
func (fld *Field) IsExported() bool { return fld.sf.IsExported() }

// IsSliceOfPtr returns true if the field is a slice or array of pointers,
// otherwise false.
func (fld *Field) IsSliceOfPtr() bool { return fld.sliceOfPtr }

// IndirectType returns the pointed-to type when the field is a pointer.
// Otherwise it returns the type set in the constructor.
func (fld *Field) IndirectType() reflect.Type {
	if fld.kind == reflect.Pointer {
		return fld.typ.Elem()
	}
	return fld.typ
}

func (fld *Field) IsSlice() bool { return fld.kind == reflect.Slice }

func (fld *Field) IsArray() bool { return fld.kind == reflect.Array }

func (fld *Field) IsSliceOrArray() bool { return fld.sliceOrArr }

func (fld *Field) IsMap() bool { return fld.kind == reflect.Map }

func (fld *Field) IsInterface() bool { return fld.kind == reflect.Interface }

// TypeMetadata returns [Metadata] for the field type.
func (fld *Field) TypeMetadata() *Metadata { return fld.metadata }

// IsAnonymous returns true for embedded fields, false otherwise.
func (fld *Field) IsAnonymous() bool { return fld.anonymous }
