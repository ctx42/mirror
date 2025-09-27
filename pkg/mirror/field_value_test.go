// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"io"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/kit/reflectkit"
)

func Test_NewFieldValue(t *testing.T) {
	t.Run("field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")

		// --- When ---
		fv := NewFieldValue(fld, val)

		// --- Then ---
		assert.Equal(t, "F", fv.Name())
	})

	t.Run("pointer field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F *string }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")

		// --- When ---
		fv := NewFieldValue(fld, val)

		// --- Then ---
		assert.Equal(t, "F", fv.Name())
	})

	t.Run("not exported field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ f string }{}
		fld := NewField(reflectkit.GetField(t, s, "f"))
		val := reflectkit.GetValue(t, s, "f")

		// --- When ---
		fv := NewFieldValue(fld, val)

		// --- Then ---
		assert.Equal(t, "f", fv.Name())
	})
}

func Test_FieldValue_Field(t *testing.T) {
	// --- Given ---
	s := &struct{ F string }{}
	fld := NewField(reflectkit.GetField(t, s, "F"))
	val := reflectkit.GetValue(t, s, "F")
	fv := NewFieldValue(fld, val)

	// --- When ---
	have := fv.Field()

	// --- Then ---
	assert.Same(t, fld, have)
}

func Test_FieldValue_Value(t *testing.T) {
	// --- Given ---
	s := &struct{ F string }{}
	fld := NewField(reflectkit.GetField(t, s, "F"))
	val := reflectkit.GetValue(t, s, "F")
	fv := NewFieldValue(fld, val)

	// --- When ---
	have := fv.Value()

	// --- Then ---
	assert.Equal(t, val, have)
}

func Test_FieldValue_StructValue(t *testing.T) {
	t.Run("field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F TStruct }{F: TStruct{FStr: "a"}}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have := fv.StructValue()

		// --- Then ---
		assert.True(t, have.IsValid())
		have.FieldByName("FStr").Value().SetString("b")
		assert.Equal(t, "b", s.F.FStr)
	})

	t.Run("pointer field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F *TStruct }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have := fv.StructValue()

		// --- Then ---
		assert.True(t, have.IsValid())
		have.NewIfNil().FieldByName("FStr").Value().SetString("b")
		assert.Equal(t, "b", s.F.FStr)
	})
}

func Test_FieldValue_NewIfNil(t *testing.T) {
	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F *TStruct }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have := fv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, s.F)
		assert.Same(t, fv, have)
	})

	t.Run("map", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F map[int]string }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have := fv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, s.F)
		assert.Same(t, fv, have)
	})

	t.Run("slice", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F []string }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have := fv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, s.F)
		assert.Same(t, fv, have)
	})

	t.Run("slice of pointers", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F []*string }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have := fv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, s.F)
		assert.Same(t, fv, have)
	})

	t.Run("slice of pointers to structs", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F []*TStruct }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have := fv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, s.F)
		assert.Same(t, fv, have)
	})

	t.Run("pointer to build in type", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F *float64 }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have := fv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, s.F)
		assert.Same(t, fv, have)
	})

	t.Run("not a pointer", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have := fv.NewIfNil()

		// --- Then ---
		assert.Same(t, fv, have)
	})
}

func Test_FieldValue_Get(t *testing.T) {
	t.Run("field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{F: "abc"}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have, err := fv.Get()

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "abc", have)
	})

	t.Run("pointer field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F *string }{F: ptr("abc")}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have, err := fv.Get()

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "abc", *have.(*string))
	})

	t.Run("unexported field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ f string }{f: "abc"}
		fld := NewField(reflectkit.GetField(t, s, "f"))
		val := reflectkit.GetValue(t, s, "f")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have, err := fv.Get()

		// --- Then ---
		assert.ErrorIs(t, ErrUnexportedField, err)
		assert.ErrorEqual(t, "unexported field: f", err)
		assert.Nil(t, have)
	})

	t.Run("invalid field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F io.Reader }{}
		fld := NewField(reflectkit.GetField(t, s, "F"))
		val := reflectkit.GetValue(t, s, "F")
		fv := NewFieldValue(fld, val)

		// --- When ---
		have, err := fv.Get()

		// --- Then ---
		assert.ErrorIs(t, ErrInvField, err)
		assert.Nil(t, have)
	})
}
