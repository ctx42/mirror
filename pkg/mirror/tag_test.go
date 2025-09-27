// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_Tag_Key(t *testing.T) {
	// --- Given ---
	tag := Tag{key: "abc"}

	// --- When ---
	have := tag.Key()

	// --- Then ---
	assert.Equal(t, "abc", have)
}

func Test_Tag_Name(t *testing.T) {
	// --- Given ---
	tag := Tag{name: "abc"}

	// --- When ---
	have := tag.Name()

	// --- Then ---
	assert.Equal(t, "abc", have)
}

func Test_Tag_Contains(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		// --- Given ---
		tag := Tag{options: []string{"a", "b", "c"}}

		// --- When ---
		have := tag.Contains("b")

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not contains", func(t *testing.T) {
		// --- Given ---
		tag := Tag{options: []string{"a", "b", "c"}}

		// --- When ---
		have := tag.Contains("x")

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("zero value tag", func(t *testing.T) {
		// --- Given ---
		tag := Tag{}

		// --- When ---
		have := tag.Contains("x")

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Tag_NameOrField(t *testing.T) {
	t.Run("zero value tag", func(t *testing.T) {
		// --- Given ---
		tag := Tag{}

		// --- When ---
		have := tag.NameOrField()

		// --- Then ---
		assert.Equal(t, "", have)
	})

	t.Run("only field name set", func(t *testing.T) {
		// --- Given ---
		tag := Tag{field: "F"}

		// --- When ---
		have := tag.NameOrField()

		// --- Then ---
		assert.Equal(t, "F", have)
	})

	t.Run("name set", func(t *testing.T) {
		// --- Given ---
		tag := Tag{field: "F", name: "name"}

		// --- When ---
		have := tag.NameOrField()

		// --- Then ---
		assert.Equal(t, "name", have)
	})

	t.Run("name set to minus", func(t *testing.T) {
		// --- Given ---
		tag := Tag{field: "F", name: "-"}

		// --- When ---
		have := tag.NameOrField()

		// --- Then ---
		assert.Equal(t, "F", have)
	})
}

func Test_Tag_IsIgnored_tabular(t *testing.T) {
	tt := []struct {
		testN string

		name string
		want bool
	}{
		{"empty", "", false},
		{"set", "abc", false},
		{"ignored", "-", true},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			tag := Tag{name: tc.name}

			// --- When ---
			have := tag.IsIgnored()

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_Tag_IsZero(t *testing.T) {
	t.Run("zero value tag", func(t *testing.T) {
		// --- Given ---
		tag := Tag{}

		// --- When ---
		have := tag.IsZero()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("options not nil but empty", func(t *testing.T) {
		// --- Given ---
		tag := Tag{options: make([]string, 0)}

		// --- When ---
		have := tag.IsZero()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("the key is not empty ", func(t *testing.T) {
		// --- Given ---
		tag := Tag{key: "abc"}

		// --- When ---
		have := tag.IsZero()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("the name is not empty ", func(t *testing.T) {
		// --- Given ---
		tag := Tag{name: "abc"}

		// --- When ---
		have := tag.IsZero()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("the field is not empty ", func(t *testing.T) {
		// --- Given ---
		tag := Tag{field: "abc"}

		// --- When ---
		have := tag.IsZero()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("the options are not empty ", func(t *testing.T) {
		// --- Given ---
		tag := Tag{options: []string{"abc"}}

		// --- When ---
		have := tag.IsZero()

		// --- Then ---
		assert.False(t, have)
	})
}
