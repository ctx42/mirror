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

// MetadataFor is a convenience function calling [TypeMetadata], see it for
// details and restrictions on argument v.
func MetadataFor(v any) *Metadata {
	return TypeMetadata(reflect.TypeOf(v))
}

// TypeMetadata returns struct metadata from the global cache or crates
// new Metadata and caches it before returning. Expects typ to be of kind
// [reflect.Struct] or [reflect.Ptr] with element of kind [reflect.Struct]
// - it will panic otherwise.
func TypeMetadata(typ reflect.Type) *Metadata {
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
