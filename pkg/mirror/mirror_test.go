// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_Reflect(t *testing.T) {
	t.Run("pointer to struct", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}

		// --- When ---
		have := Reflect(s)

		// --- Then ---
		assert.NotNil(t, have)

		typCacheMX.Lock()
		defer typCacheMX.Unlock()
		cached, ok := typCache[reflect.TypeOf(s).Elem()]
		assert.True(t, ok)
		assert.Same(t, have, cached)
	})

	t.Run("struct", func(t *testing.T) {
		// --- Given ---
		s := struct{ F string }{}

		// --- When ---
		have := Reflect(s)

		// --- Then ---
		assert.NotNil(t, have)

		typCacheMX.Lock()
		defer typCacheMX.Unlock()
		cached, ok := typCache[reflect.TypeOf(s)]
		assert.True(t, ok)
		assert.Same(t, have, cached)
	})

	t.Run("int", func(t *testing.T) {
		// --- Given ---
		i := 42

		// --- When ---
		have := Reflect(i)

		// --- Then ---
		assert.NotNil(t, have)

		typCacheMX.Lock()
		defer typCacheMX.Unlock()
		cached, ok := typCache[reflect.TypeOf(i)]
		assert.True(t, ok)
		assert.Same(t, have, cached)
	})
}
