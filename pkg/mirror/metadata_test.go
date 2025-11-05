// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/check"
)

func Test_NewMetadata(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}

		// --- When ---
		have := NewMetadata(s)

		// --- Then ---
		assert.NotNil(t, have.typ)
		assert.Len(t, 1, have.fields)
		assert.Equal(t, "", have.name)
		assert.Equal(t, "", have.pkg)

	})

	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}

		// --- When ---
		have := NewMetadata(s)

		// --- Then ---
		assert.NotNil(t, have.typ)
		assert.Len(t, 1, have.fields)
		assert.Equal(t, "", have.name)
		assert.Equal(t, "", have.pkg)
	})

	t.Run("type of nil panics", func(t *testing.T) {
		assert.Panic(t, func() { NewMetadata(reflect.TypeOf(nil)) })
	})

	t.Run("nil panics", func(t *testing.T) {
		assert.Panic(t, func() { NewMetadata(nil) })
	})

	t.Run("not a struct", func(t *testing.T) {
		// --- Given ---
		var i int

		// --- When ---
		have := NewMetadata(i)

		// --- Then ---
		assert.Equal(t, reflect.TypeOf(i), have.typ)
		assert.Equal(t, reflect.TypeOf(i).Kind(), have.kind)
		assert.Nil(t, have.fields)
		assert.Equal(t, "int", have.name)
		assert.Equal(t, "", have.pkg)
	})

	t.Run("pointer to not a struct", func(t *testing.T) {
		// --- Given ---
		var i *int

		// --- When ---
		have := NewMetadata(i)

		// --- Then ---
		assert.Equal(t, reflect.TypeOf(i).Elem(), have.typ)
		assert.Equal(t, reflect.TypeOf(i).Elem().Kind(), have.kind)
		assert.Nil(t, have.fields)
		assert.Equal(t, "int", have.name)
		assert.Equal(t, "", have.pkg)
	})
}

func Test_NewTypeMetadata(t *testing.T) {
	// Tested by Test_NewMetadata.
}

func Test_NewValueMetadata(t *testing.T) {
	t.Run("type", func(t *testing.T) {
		// --- Given ---
		val := reflect.ValueOf(bytes.Buffer{})

		// --- When ---
		have := NewValueMetadata(val)

		// --- Then ---
		assert.Equal(t, "Buffer", have.name)
		assert.Equal(t, "bytes", have.pkg)
	})

	t.Run("pointer type", func(t *testing.T) {
		// --- Given ---
		val := reflect.ValueOf(&bytes.Buffer{})

		// --- When ---
		have := NewValueMetadata(val)

		// --- Then ---
		assert.Equal(t, "Buffer", have.name)
		assert.Equal(t, "bytes", have.pkg)
	})

	t.Run("sdk func", func(t *testing.T) {
		// --- Given ---
		val := reflect.ValueOf(reflect.DeepEqual)

		// --- When ---
		have := NewValueMetadata(val)

		// --- Then ---
		assert.Equal(t, "DeepEqual", have.name)
		assert.Equal(t, "reflect", have.pkg)
	})

	t.Run("external func", func(t *testing.T) {
		// --- Given ---
		val := reflect.ValueOf(check.After)

		// --- When ---
		have := NewValueMetadata(val)

		// --- Then ---
		assert.Equal(t, "After", have.name)
		assert.Equal(t, "github.com/ctx42/testing/pkg/check", have.pkg)
	})
}

func Test_Metadata_Type(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}
		md := NewMetadata(s)

		// --- When ---
		have := md.Type()

		// --- Then ---
		assert.Equal(t, reflect.TypeOf(s), have)
		assert.Equal(t, reflect.Struct, have.Kind())
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		md := NewMetadata(s)

		// --- When ---
		have := md.Type()

		// --- Then ---
		assert.Equal(t, reflect.TypeOf(s).Elem(), have)
		assert.Equal(t, reflect.Struct, have.Kind())
	})
}

func Test_Metadata_Kind(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}
		md := NewMetadata(s)

		// --- When ---
		have := md.Kind()

		// --- Then ---
		assert.Equal(t, reflect.Struct, have)
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		md := NewMetadata(s)

		// --- When ---
		have := md.Kind()

		// --- Then ---
		assert.Equal(t, reflect.Struct, have)
	})
}

func Test_Metadata_Name(t *testing.T) {
	// --- Given ---
	md := &Metadata{name: "abc"}

	// --- When ---
	have := md.Name()

	// --- Then ---
	assert.Equal(t, "abc", have)
}

func Test_Metadata_Package(t *testing.T) {
	// --- Given ---
	md := &Metadata{pkg: "abc"}

	// --- When ---
	have := md.Package()

	// --- Then ---
	assert.Equal(t, "abc", have)
}

func Test_Metadata_IsStruct(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}
		md := NewMetadata(s)

		// --- When ---
		have := md.IsStruct()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}
		md := NewMetadata(s)

		// --- When ---
		have := md.IsStruct()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not struct", func(t *testing.T) {
		// --- Given ---
		var i int
		md := NewMetadata(i)

		// --- When ---
		have := md.IsStruct()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Metadata_Fields(t *testing.T) {
	t.Run("struct with fields", func(t *testing.T) {
		// --- Given ---
		md := NewMetadata(TStruct{})

		// --- When ---
		fields := md.Fields()

		// --- Then ---
		assert.Len(t, 9, fields)
		assert.Equal(t, "FStr", fields[0].Name())
		assert.Equal(t, "fStr", fields[1].Name())
		assert.Equal(t, "FpStr", fields[2].Name())
		assert.Equal(t, "FsStr", fields[3].Name())
		assert.Equal(t, "FaStr", fields[4].Name())
		assert.Equal(t, "FmStr", fields[5].Name())
		assert.Equal(t, "SPtr", fields[6].Name())
		assert.Equal(t, "SVal", fields[7].Name())
		assert.Equal(t, "SNil", fields[8].Name())
	})

	t.Run("struct without fields", func(t *testing.T) {
		// --- Given ---
		md := NewMetadata(struct{}{})

		// --- When ---
		fields := md.Fields()

		// --- Then ---
		assert.Nil(t, fields)
	})
}

func Test_Metadata_FieldByName(t *testing.T) {
	t.Run("known exported field", func(t *testing.T) {
		// --- Given ---
		md := NewMetadata(TStruct{})

		// --- When ---
		have := md.FieldByName("FStr")

		// --- Then ---
		assert.Equal(t, "FStr", have.Name())
	})

	t.Run("known unexported field", func(t *testing.T) {
		// --- Given ---
		md := NewMetadata(TStruct{})

		// --- When ---
		have := md.FieldByName("fStr")

		// --- Then ---
		assert.Equal(t, "fStr", have.Name())
	})

	t.Run("unknown field", func(t *testing.T) {
		// --- Given ---
		md := NewMetadata(TStruct{})

		// --- When ---
		have := md.FieldByName("Unknown")

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_Metadata_FieldByIndex(t *testing.T) {
	t.Run("known exported field", func(t *testing.T) {
		// --- Given ---
		md := NewMetadata(TStruct{})

		// --- When ---
		have := md.FieldByIndex(2)

		// --- Then ---
		assert.Equal(t, "FpStr", have.Name())
	})

	t.Run("known unexported field", func(t *testing.T) {
		// --- Given ---
		md := NewMetadata(TStruct{})

		// --- When ---
		have := md.FieldByIndex(1)

		// --- Then ---
		assert.Equal(t, "fStr", have.Name())
	})

	t.Run("unknown field", func(t *testing.T) {
		// --- Given ---
		md := NewMetadata(TStruct{})

		// --- When ---
		have := md.FieldByIndex(42)

		// --- Then ---
		assert.Nil(t, have)
	})
}
