// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/kit/reflectkit"
)

func Test_NewField(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")

		// --- When ---
		fld := NewField(sf)

		// --- Then ---
		assert.Equal(t, sf, fld.sf)
		assert.Nil(t, fld.tags)
	})
}

func Test_Field_StructField(t *testing.T) {
	// --- Given ---
	s := &struct{ F string }{}
	sf := reflectkit.GetField(t, s, "F")
	fld := NewField(sf)

	// --- When ---
	have := fld.StructField()

	// --- Then ---
	assert.Equal(t, sf, have)
}

func Test_Field_Type(t *testing.T) {
	// --- Given ---
	s := &struct{ F string }{}
	sf := reflectkit.GetField(t, s, "F")
	fld := NewField(sf)

	// --- When ---
	have := fld.Type()

	// --- Then ---
	assert.Equal(t, sf.Type, have)
}

func Test_Field_Kind(t *testing.T) {
	// --- Given ---
	s := &struct{ F string }{}
	sf := reflectkit.GetField(t, s, "F")
	fld := NewField(sf)

	// --- When ---
	have := fld.Kind()

	// --- Then ---
	assert.Equal(t, sf.Type.Kind(), have)
}

func Test_Field_Index(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		// --- Given ---
		s := &struct {
			F0 string
			F1 int
		}{}
		sf := reflectkit.GetField(t, s, "F1")
		fld := NewField(sf)

		// --- When ---
		have := fld.Index()

		// --- Then ---
		assert.Equal(t, []int{1}, have)
	})
}

func Test_Field_Name(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.Name()

		// --- Then ---
		assert.Equal(t, "F", have)
	})
}

func Test_Field_Tag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		s := &struct {
			F string `tag:"t1,t2, t3"`
		}{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.Tag("tag")

		// --- Then ---
		want := Tag{
			field:   "F",
			key:     "tag",
			name:    "t1",
			options: []string{"t2", "t3"},
		}
		assert.Equal(t, want, have)
	})

	t.Run("tag not existing", func(t *testing.T) {
		// --- Given ---
		s := &struct {
			F string `tag:"t1,t2, t3"`
		}{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.Tag("other")

		// --- Then ---
		assert.True(t, have.IsZero())
	})
}

func Test_Field_IsValid(t *testing.T) {
	t.Run("is valid", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsValid()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("embedded interface", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F io.Reader }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsValid()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("embedded field", func(t *testing.T) {
		// --- Given ---
		s := struct{ bytes.Buffer }{}
		sf := reflect.TypeOf(s).Field(0)
		fld := NewField(sf)

		// --- When ---
		have := fld.IsValid()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Field_IsSliceOfPtr(t *testing.T) {
	t.Run("slice of pointers", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F []*string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsSliceOfPtr()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("array of pointers", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F [1]*string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsSliceOfPtr()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not slice pointer", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F *string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsSliceOfPtr()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Field_IndirectType(t *testing.T) {
	t.Run("pointer type", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F *string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IndirectType()

		// --- Then ---
		assert.Equal(t, have, reflect.TypeOf(""))
	})

	t.Run("not pointer type", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IndirectType()

		// --- Then ---
		assert.Equal(t, have, reflect.TypeOf(""))
	})
}

func Test_Field_IsSlice(t *testing.T) {
	t.Run("slice", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F []string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsSlice()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("array", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F [2]string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsSlice()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("not slice or array", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsSlice()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Field_IsArray(t *testing.T) {
	t.Run("array", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F [2]string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsArray()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("slice", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F []string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsArray()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("not slice or array", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsArray()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Field_IsSliceOrArray(t *testing.T) {
	t.Run("array", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F [2]string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsSliceOrArray()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("slice", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F []string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsSliceOrArray()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not slice or array", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsSliceOrArray()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Field_IsMap(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F map[int]string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsMap()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("slice", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F []string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsMap()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("array", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F [2]string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsMap()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("not slice or array", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsMap()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Field_IsInterface(t *testing.T) {
	t.Run("interface", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F io.Reader }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsInterface()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not interface", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsInterface()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Field_IsStruct(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F bytes.Buffer }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsStruct()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("interface", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F io.Reader }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsStruct()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("string", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsStruct()

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Field_IsExported(t *testing.T) {
	s := struct {
		FStr string
		fStr string
	}{}

	typ := reflect.TypeOf(s)

	t.Run("exported", func(t *testing.T) {
		// --- Given ---
		f := typ.Field(0)

		// --- When ---
		sf := NewField(f)

		// --- Then ---
		assert.True(t, sf.exported)
	})

	t.Run("not-exported", func(t *testing.T) {
		// --- Given ---
		f := typ.Field(1)

		// --- When ---
		sf := NewField(f)

		// --- Then ---
		assert.False(t, sf.exported)
	})
}

func Test_Field_TypeMetadata(t *testing.T) {
	t.Run("metadata", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.TypeMetadata()

		// --- Then ---
		assert.False(t, have.IsStruct())
	})
}

func Test_Field_IsAnonymous(t *testing.T) {
	t.Run("not anonymous", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsAnonymous()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("embedded interface", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F io.Reader }{}
		sf := reflectkit.GetField(t, s, "F")
		fld := NewField(sf)

		// --- When ---
		have := fld.IsAnonymous()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("embedded field", func(t *testing.T) {
		// --- Given ---
		s := struct{ bytes.Buffer }{}
		sf := reflect.TypeOf(s).Field(0)
		fld := NewField(sf)

		// --- When ---
		have := fld.IsAnonymous()

		// --- Then ---
		assert.True(t, have)
	})
}
