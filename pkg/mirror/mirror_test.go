// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/check"
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

	t.Run("pointer to pointer shares struct cache", func(t *testing.T) {
		// --- Given ---
		s := &struct{ F string }{}
		p := &s

		// --- When ---
		have := Reflect(p)

		// --- Then ---
		assert.Same(t, Reflect(s), have)
		assert.Equal(t, reflect.Struct, have.Kind())
	})
}

func Test_ReflectValue(t *testing.T) {
	t.Run("func", func(t *testing.T) {
		// --- Given ---
		val := reflect.ValueOf(check.After)

		// --- When ---
		have := ReflectValue(val)

		// --- Then ---
		assert.Equal(t, reflect.Func, have.kind)
		assert.Nil(t, have.fields)
		assert.Equal(t, "github.com/ctx42/testing/pkg/check", have.pkg)
		assert.Equal(t, "After", have.name)
	})

	t.Run("uses cache", func(t *testing.T) {
		// --- Given ---
		dm := ReflectValue(reflect.ValueOf(check.After))

		// --- When ---
		have := ReflectValue(reflect.ValueOf(check.After))

		// --- Then ---
		assert.Same(t, dm, have)
	})

	t.Run("func pointer", func(t *testing.T) {
		// --- Given ---
		tst := &TwoStr{}
		val := reflect.ValueOf(tst)

		// --- When ---
		have := ReflectValue(val)

		// --- Then ---
		assert.Equal(t, reflect.Struct, have.kind)
		assert.NotNil(t, have.fields)
		assert.Equal(t, "github.com/ctx42/mirror/pkg/mirror", have.pkg)
		assert.Equal(t, "TwoStr", have.name)
	})

	t.Run("same signature keeps its own name", func(t *testing.T) {
		// --- Given ---
		first := reflect.ValueOf(fxA)
		second := reflect.ValueOf(fxB)
		ReflectValue(first)

		// --- When ---
		have := ReflectValue(second)

		// --- Then ---
		assert.Equal(t, "fxB", have.Name())
		assert.Equal(t, "github.com/ctx42/mirror/pkg/mirror", have.Package())
	})

	t.Run("type reflect keeps the value name", func(t *testing.T) {
		// --- Given ---
		Reflect(fy)
		val := reflect.ValueOf(fy)

		// --- When ---
		have := ReflectValue(val)

		// --- Then ---
		assert.Equal(t, "fy", have.Name())
	})

	t.Run("pointer to func", func(t *testing.T) {
		// --- Given ---
		fn := fy
		p := &fn
		val := reflect.ValueOf(p)

		// --- When ---
		have := ReflectValue(val)

		// --- Then ---
		assert.Equal(t, "fy", have.Name())
		assert.Equal(t, reflect.Func, have.Kind())
	})
}

func fxA(int) int { return 1 }

func fxB(int) int { return 2 }

func fy(string) string { return "" }
