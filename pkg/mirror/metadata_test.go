// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NewAnyMetadata_NewTypeMetadata(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}

		// --- When ---
		md := NewAnyMetadata(s)

		// --- Then ---
		assert.NotNil(t, md.typ)
		assert.Len(t, 1, md.fields)
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}

		// --- When ---
		md := NewAnyMetadata(s)

		// --- Then ---
		assert.NotNil(t, md.typ)
		assert.Len(t, 1, md.fields)
	})

	t.Run("type of nil panics", func(t *testing.T) {
		assert.Panic(t, func() { NewAnyMetadata(reflect.TypeOf(nil)) })
	})

	t.Run("nil panics", func(t *testing.T) {
		assert.Panic(t, func() { NewAnyMetadata(nil) })
	})

	t.Run("not a struct", func(t *testing.T) {
		// --- Given ---
		var i int

		// --- When ---
		have := NewAnyMetadata(i)

		// --- Then ---
		assert.Equal(t, reflect.TypeOf(i), have.typ)
		assert.Equal(t, reflect.TypeOf(i).Kind(), have.kind)
		assert.Nil(t, have.fields)
	})

	t.Run("pointer to not a struct", func(t *testing.T) {
		// --- Given ---
		var i *int

		// --- When ---
		have := NewAnyMetadata(i)

		// --- Then ---
		assert.Equal(t, reflect.TypeOf(i).Elem(), have.typ)
		assert.Equal(t, reflect.TypeOf(i).Elem().Kind(), have.kind)
		assert.Nil(t, have.fields)
	})
}

func Test_Metadata_Type(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}
		md := NewAnyMetadata(s)

		// --- When ---
		have := md.Type()

		// --- Then ---
		assert.Equal(t, reflect.TypeOf(s), have)
		assert.Equal(t, reflect.Struct, have.Kind())
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		md := NewAnyMetadata(s)

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
		md := NewAnyMetadata(s)

		// --- When ---
		have := md.Kind()

		// --- Then ---
		assert.Equal(t, reflect.Struct, have)
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		md := NewAnyMetadata(s)

		// --- When ---
		have := md.Kind()

		// --- Then ---
		assert.Equal(t, reflect.Struct, have)
	})
}

func Test_Metadata_IsStruct(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}
		md := NewAnyMetadata(s)

		// --- When ---
		have := md.IsStruct()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}
		md := NewAnyMetadata(s)

		// --- When ---
		have := md.IsStruct()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not struct", func(t *testing.T) {
		// --- Given ---
		var i int
		md := NewAnyMetadata(i)

		// --- When ---
		have := md.IsStruct()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Metadata_Fields(t *testing.T) {
	t.Run("struct with fields", func(t *testing.T) {
		// --- Given ---
		md := NewAnyMetadata(TStruct{})

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
		md := NewAnyMetadata(struct{}{})

		// --- When ---
		fields := md.Fields()

		// --- Then ---
		assert.Nil(t, fields)
	})
}

func Test_Metadata_FieldByName(t *testing.T) {
	t.Run("known exported field", func(t *testing.T) {
		// --- Given ---
		md := NewAnyMetadata(TStruct{})

		// --- When ---
		have := md.FieldByName("FStr")

		// --- Then ---
		assert.Equal(t, "FStr", have.Name())
	})

	t.Run("known unexported field", func(t *testing.T) {
		// --- Given ---
		md := NewAnyMetadata(TStruct{})

		// --- When ---
		have := md.FieldByName("fStr")

		// --- Then ---
		assert.Equal(t, "fStr", have.Name())
	})

	t.Run("unknown field", func(t *testing.T) {
		// --- Given ---
		md := NewAnyMetadata(TStruct{})

		// --- When ---
		have := md.FieldByName("Unknown")

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_Metadata_FieldByIndex(t *testing.T) {
	t.Run("known exported field", func(t *testing.T) {
		// --- Given ---
		md := NewAnyMetadata(TStruct{})

		// --- When ---
		have := md.FieldByIndex(2)

		// --- Then ---
		assert.Equal(t, "FpStr", have.Name())
	})

	t.Run("known unexported field", func(t *testing.T) {
		// --- Given ---
		md := NewAnyMetadata(TStruct{})

		// --- When ---
		have := md.FieldByIndex(1)

		// --- Then ---
		assert.Equal(t, "fStr", have.Name())
	})

	t.Run("unknown field", func(t *testing.T) {
		// --- Given ---
		md := NewAnyMetadata(TStruct{})

		// --- When ---
		have := md.FieldByIndex(42)

		// --- Then ---
		assert.Nil(t, have)
	})
}
