// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"errors"
	"reflect"
	"sync"
)

// Sentinel errors.
var (
	// ErrInvField represents invalid fields error. Field is invalid if:
	//   - the field does not exist,
	//   - the field type is nil,
	//   - reflect.Value.IsValid returns false.
	ErrInvField = errors.New("invalid field")

	// ErrUnexportedField represents error when accessing unexported field.
	ErrUnexportedField = errors.New("unexported field")
)

var (
	typCache   map[reflect.Type]*Metadata // Type metadata cache.
	typCacheMX sync.RWMutex               // Guards typCache.
)

func init() { typCache = map[reflect.Type]*Metadata{} }

// Reflect extracts [Metadata] about type of "v".
func Reflect(v any) *Metadata {
	return ReflectType(reflect.TypeOf(v))
}

// ReflectType extracts [Metadata] about the type.
func ReflectType(typ reflect.Type) *Metadata {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	typCacheMX.RLock()
	md, found := typCache[typ]
	typCacheMX.RUnlock()
	if found {
		return md
	}

	md = NewTypeMetadata(typ)
	typCacheMX.Lock()
	typCache[typ] = md
	typCacheMX.Unlock()
	return md
}

// ReflectValue extracts [Metadata] about the value.
func ReflectValue(val reflect.Value) *Metadata {
	typ := val.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	typCacheMX.RLock()
	md, found := typCache[typ]
	typCacheMX.RUnlock()
	if found {
		return md
	}

	md = NewValueMetadata(val)
	typCacheMX.Lock()
	typCache[typ] = md
	typCacheMX.Unlock()
	return md
}
