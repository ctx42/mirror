// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NewStructValue(t *testing.T) {
	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}

		// --- When ---
		have := NewStructValue(s)

		// --- Then ---
		assert.NotNil(t, have)
		assert.True(t, have.IsValid())
	})

	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F int }{}

		// --- When ---
		have := NewStructValue(s)

		// --- Then ---
		assert.Nil(t, have)
		assert.False(t, have.IsValid())
	})

	t.Run("not a struct", func(t *testing.T) {
		// --- Given ---
		var s *int

		// --- When ---
		have := NewStructValue(s)

		// --- Then ---
		assert.Nil(t, have)
		assert.False(t, have.IsValid())
	})

	t.Run("pointer to struct variable", func(t *testing.T) {
		// --- Given ---
		var s *TStruct

		// --- When ---
		have := NewStructValue(s)

		// --- Then ---
		assert.Nil(t, have)
		assert.False(t, have.IsValid())
	})

	t.Run("nil", func(t *testing.T) {
		// --- When ---
		sv := NewStructValue(nil)

		// --- Then ---
		assert.Nil(t, sv)
	})
}

func Test_NewStructValue_IsPtr(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.IsPtr()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("struct pointer", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.IsPtr()

		// --- Then ---
		assert.True(t, have)
	})
}

func Test_StructValue_IsValid(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.IsValid()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("invalid", func(t *testing.T) {
		// --- Given ---
		var s *TStruct
		sv := NewStructValue(s)

		// --- When ---
		have := sv.IsValid()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_StructValue_Metadata(t *testing.T) {
	// --- Given ---
	s := &struct{ F int }{}
	sv := NewStructValue(s)

	// --- When ---
	have := sv.Metadata()

	// --- Then ---
	assert.Same(t, sv.metadata, have)
}

func Test_StructValue_NumFields(t *testing.T) {
	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &TStruct{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.NumField()

		// --- Then ---
		assert.Equal(t, 9, have)
	})
}

func Test_StructValue_FieldByName(t *testing.T) {
	t.Run("field of a struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByName("F")

		// --- Then ---
		assert.NotNil(t, have)
	})

	t.Run("field of a pointer to a struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByName("F")

		// --- Then ---
		assert.NotNil(t, have)
	})

	t.Run("not existing field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByName("NotExisting")

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("not existing field of a pointer to a struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F int }{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByName("NotExisting")

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("unexported field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ f int }{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByName("f")

		// --- Then ---
		assert.NotNil(t, have)
		assert.False(t, have.IsExported())
	})
}

func Test_StructValue_FieldByIndex(t *testing.T) {
	t.Run("field of a struct", func(t *testing.T) {
		// --- Given ---
		s := &struct {
			F0 int
			F1 int
		}{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByIndex(1)

		// --- Then ---
		assert.NotNil(t, have)
		assert.Equal(t, "F1", have.Name())
	})

	t.Run("field of a pointer to a struct", func(t *testing.T) {
		// --- Given ---
		s := &struct {
			F0 int
			F1 int
		}{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByIndex(1)

		// --- Then ---
		assert.NotNil(t, have)
		assert.Equal(t, "F1", have.Name())
	})

	t.Run("not existing field", func(t *testing.T) {
		// --- Given ---
		s := &struct {
			F0 int
			F1 int
		}{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByIndex(42)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("not existing field of a pointer to a struct", func(t *testing.T) {
		// --- Given ---
		s := &struct {
			F0 int
			F1 int
		}{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByIndex(42)

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("unexported field", func(t *testing.T) {
		// --- Given ---
		s := &struct {
			f int
			F int
		}{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.FieldByIndex(0)

		// --- Then ---
		assert.NotNil(t, have)
		assert.Equal(t, "f", have.Name())
	})
}

func Test_StructValue_NewIfNil(t *testing.T) {
	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &TStruct{}
		sv := NewStructValue(s)

		// --- When ---
		have := sv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, sv)
		sv.FieldByName("FStr").Value().SetString("abc")
		assert.Equal(t, "abc", s.FStr)
		assert.Same(t, sv, have)
	})

	t.Run("struct field which is struct pointer", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F *TStruct }{}
		sv := NewStructValue(s).FieldByName("F").StructValue()

		// --- When ---
		have := sv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, sv)
		assert.NotNil(t, s.F)
		assert.Same(t, sv, have)
	})

	t.Run("time.Time struct field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F time.Time }{}
		sv := NewStructValue(s).FieldByName("F").StructValue()

		// --- When ---
		have := sv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, sv)
		assert.Zero(t, s.F)
		assert.Same(t, sv, have)
	})

	t.Run("pointer to time.Time struct field", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F *time.Time }{}
		sv := NewStructValue(s).FieldByName("F").StructValue()

		// --- When ---
		have := sv.NewIfNil()

		// --- Then ---
		assert.NotNil(t, sv)
		assert.NotNil(t, s.F)
		assert.Zero(t, s.F)
		assert.Same(t, sv, have)
	})
}
