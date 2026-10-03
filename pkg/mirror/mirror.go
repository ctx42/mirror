// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

// Package mirror is a lightweight library that provides a simplified interface
// for reflecting structs. The metadata about the struct and its fields is
// cached to improve performance.
package mirror

import (
	"errors"
	"reflect"
	"sync"
)

// Sentinel errors.
var (
	// ErrInvField is returned by [FieldValue.Get] for an interface field or
	// an embedded field. A missing field name returns nil, not this error.
	ErrInvField = errors.New("invalid field")

	// ErrUnexportedField represents an error when accessing an unexported
	// field.
	ErrUnexportedField = errors.New("unexported field")

	// ErrTagSyntax represents the error when parsing a struct field tag.
	ErrTagSyntax = errors.New("struct field tag syntax error")
)

var (
	typCache   = map[reflect.Type]*Metadata{} // Type metadata cache.
	typCacheMX sync.RWMutex                   // Guards typCache.

	fnCache   = map[uintptr]*Metadata{} // Function value metadata cache.
	fnCacheMX sync.RWMutex              // Guards fnCache.
)

// Reflect extracts [Metadata] about type of "v". Panics if v is an untyped
// nil.
func Reflect(v any) *Metadata {
	return ReflectType(reflect.TypeOf(v))
}

// ReflectType extracts [Metadata] about the type. Panics if typ is nil.
func ReflectType(typ reflect.Type) *Metadata {
	typ = indirect(typ)
	return cachedType(typ, func() *Metadata {
		return NewTypeMetadata(typ)
	})
}

// ReflectValue extracts [Metadata] about the value.
// A function value keeps the runtime name of that function.
// Panics when val is the zero [reflect.Value].
func ReflectValue(val reflect.Value) *Metadata {
	if !val.IsValid() {
		panic("mirror: ReflectValue of invalid value")
	}
	if fn, ok := funcValue(val); ok {
		return cachedFunc(fn)
	}

	typ := indirect(val.Type())
	return cachedType(typ, func() *Metadata {
		return NewValueMetadata(val)
	})
}

// cachedType returns the metadata cached for typ, building it once.
// Concurrent callers for one type share the stored pointer.
func cachedType(typ reflect.Type, build func() *Metadata) *Metadata {
	typCacheMX.RLock()
	md, found := typCache[typ]
	typCacheMX.RUnlock()
	if found {
		return md
	}

	md = build()
	typCacheMX.Lock()
	if existing, ok := typCache[typ]; ok {
		typCacheMX.Unlock()
		return existing
	}
	typCache[typ] = md
	typCacheMX.Unlock()
	return md
}

// funcValue returns val when it is a function, following pointers.
// The boolean is false for a nil pointer or any other kind.
func funcValue(val reflect.Value) (reflect.Value, bool) {
	for val.IsValid() && val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return reflect.Value{}, false
		}
		val = val.Elem()
	}
	if !val.IsValid() || val.Kind() != reflect.Func || val.Pointer() == 0 {
		return reflect.Value{}, false
	}
	return val, true
}

// cachedFunc returns metadata for one function value.
// Repeated calls with the same code pointer return the same metadata.
func cachedFunc(fn reflect.Value) *Metadata {
	pc := fn.Pointer()

	fnCacheMX.RLock()
	md, found := fnCache[pc]
	fnCacheMX.RUnlock()
	if found {
		return md
	}

	md = NewValueMetadata(fn)
	fnCacheMX.Lock()
	if existing, ok := fnCache[pc]; ok {
		fnCacheMX.Unlock()
		return existing
	}
	fnCache[pc] = md
	fnCacheMX.Unlock()
	return md
}
