// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

import (
	"maps"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

// Types for the [Names] tests.
type (
	tnInner struct {
		A int `k:"a"`
		B int
	}

	tnOther struct {
		A int `k:"a"`
		C int
	}

	tnTagged struct {
		C int `k:"C"`
	}

	tnItf interface{ M() }

	tnRec struct {
		*tnRec //nolint:unused
		R      int
	}

	tnOuter struct {
		tnInner          // Promotes "a" and "B".
		*tnTagged        // Promotes "C", tagged.
		tnOther          // Shares "a" with tnInner, "C" with tnTagged.
		tnItf            // Embedded interfaces are skipped.
		Ign       int    `k:"-"`
		Shared    int    `k:"-"`
		Name      string `k:"Shared"`
		B         int    `k:"top"`
		hidden    int    //nolint:unused
	}
)

func Test_newNames(t *testing.T) {
	t.Run("tag name and Go name", func(t *testing.T) {
		// --- Given ---
		type T struct {
			A int `k:"a"`
			B int
			C int `k:",opt"`
		}

		// --- When ---
		have := newNames(NewMetadata(T{}), "k")

		// --- Then ---
		want := map[string][]int{"a": {0}, "B": {1}, "C": {2}}
		assert.Equal(t, want, have.fields)
		assert.Equal(t, map[string]bool{}, have.ignored)
	})

	t.Run("promotion", func(t *testing.T) {
		// --- When ---
		have := newNames(NewMetadata(tnOuter{}), "k")

		// --- Then ---
		want := map[string][]int{
			"B":      {0, 1},
			"C":      {1, 0},
			"Shared": {6},
			"top":    {7},
		}
		assert.Equal(t, want, have.fields)
		assert.Equal(t, map[string]bool{"Ign": true}, have.ignored)
	})

	t.Run("shallower field wins", func(t *testing.T) {
		// --- Given ---
		type T struct {
			tnInner
			B int
		}

		// --- When ---
		have := newNames(NewMetadata(T{}), "k")

		// --- Then ---
		want := map[string][]int{"a": {0, 0}, "B": {1}}
		assert.Equal(t, want, have.fields)
	})

	t.Run("ignored embedded struct", func(t *testing.T) {
		// --- Given ---
		type T struct {
			tnInner `k:"-"`
			D       int
		}

		// --- When ---
		have := newNames(NewMetadata(T{}), "k")

		// --- Then ---
		assert.Equal(t, map[string][]int{"D": {1}}, have.fields)
		assert.Equal(t, map[string]bool{}, have.ignored)
	})

	t.Run("recursive type", func(t *testing.T) {
		// --- When ---
		have := newNames(NewMetadata(tnRec{}), "k")

		// --- Then ---
		assert.Equal(t, map[string][]int{"R": {1}}, have.fields)
	})

	t.Run("not a struct", func(t *testing.T) {
		// --- When ---
		have := newNames(NewMetadata(42), "k")

		// --- Then ---
		assert.Equal(t, map[string][]int{}, have.fields)
	})
}

func Test_Names_Field(t *testing.T) {
	t.Run("selects a field", func(t *testing.T) {
		// --- Given ---
		nms := NewMetadata(tnOuter{}).Names("k")

		// --- When ---
		have, ok := nms.Field("C")

		// --- Then ---
		assert.True(t, ok)
		assert.Equal(t, []int{1, 0}, have)
	})

	t.Run("selects none", func(t *testing.T) {
		// --- Given ---
		nms := NewMetadata(tnOuter{}).Names("k")

		// --- When ---
		have, ok := nms.Field("a")

		// --- Then ---
		assert.False(t, ok)
		assert.Nil(t, have)
	})
}

func Test_Names_Selects_tabular(t *testing.T) {
	tt := []struct {
		testN string

		name string
		at   []int
		i    int
		want bool
	}{
		{"top level", "top", nil, 7, true},
		{"top level other index", "top", nil, 6, false},
		{"embedded", "C", []int{1}, 0, true},
		{"embedded other struct", "C", []int{2}, 0, false},
		{"embedded at top level", "C", nil, 1, false},
		{"unknown name", "x", nil, 0, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			nms := NewMetadata(tnOuter{}).Names("k")

			// --- When ---
			have := nms.Selects(tc.name, tc.at, tc.i)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_Names_All(t *testing.T) {
	t.Run("all names", func(t *testing.T) {
		// --- Given ---
		nms := NewMetadata(tnOuter{}).Names("k")

		// --- When ---
		have := maps.Collect(nms.All())

		// --- Then ---
		want := map[string][]int{
			"B":      {0, 1},
			"C":      {1, 0},
			"Shared": {6},
			"top":    {7},
		}
		assert.Equal(t, want, have)
	})

	t.Run("stops early", func(t *testing.T) {
		// --- Given ---
		nms := NewMetadata(tnOuter{}).Names("k")
		var have int

		// --- When ---
		for range nms.All() {
			have++
			break
		}

		// --- Then ---
		assert.Equal(t, 1, have)
	})
}

func Test_Names_IsIgnored(t *testing.T) {
	t.Run("ignored", func(t *testing.T) {
		// --- Given ---
		nms := NewMetadata(tnOuter{}).Names("k")

		// --- When ---
		have := nms.IsIgnored("Ign")

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("name another field takes", func(t *testing.T) {
		// --- Given ---
		nms := NewMetadata(tnOuter{}).Names("k")

		// --- When ---
		have := nms.IsIgnored("Shared")

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Names_Len(t *testing.T) {
	// --- Given ---
	nms := NewMetadata(tnOuter{}).Names("k")

	// --- When ---
	have := nms.Len()

	// --- Then ---
	assert.Equal(t, 4, have)
}
