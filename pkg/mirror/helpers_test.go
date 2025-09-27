// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_ParseTags(t *testing.T) {
	tt := []struct {
		testN string

		fieldName string
		fieldTag  string
		want      []Tag
	}{
		{
			"empty tags",
			"Field",
			"",
			nil,
		},
		{
			"tags with only spaces",
			"Field",
			` `,
			nil,
		},
		{
			"one tag one option",
			"FieldA",
			`tag:"t1"`,
			[]Tag{
				{
					field:   "FieldA",
					key:     "tag",
					name:    "t1",
					options: nil,
				},
			},
		},
		{
			"one tag no options",
			"FieldB",
			`tag:""`,
			[]Tag{
				{
					field:   "FieldB",
					key:     "tag",
					name:    "",
					options: nil,
				},
			},
		},
		{
			"one tag two options",
			"FieldC",
			`tag:"t1,t2"`,
			[]Tag{
				{
					field:   "FieldC",
					key:     "tag",
					name:    "t1",
					options: []string{"t2"},
				},
			},
		},
		{
			"options strings are trimmed",
			"Field",
			`tag:" t1, t2, o3 "`,
			[]Tag{
				{
					field:   "Field",
					key:     "tag",
					name:    "t1",
					options: []string{"t2", "o3"},
				},
			},
		},
		{
			"empty options are removed",
			"Field",
			`tag:"t1,,t2,"`,
			[]Tag{
				{
					field:   "Field",
					key:     "tag",
					name:    "t1",
					options: []string{"t2"},
				},
			},
		},
		{
			"options with whitespace",
			"Field",
			`tag:" "`,
			[]Tag{
				{
					field:   "Field",
					key:     "tag",
					name:    "Field",
					options: nil,
				},
			},
		},
		{
			"tag with leading spaces",
			"Field",
			` tag:"t1,t2"`,
			[]Tag{
				{
					field:   "Field",
					key:     "tag",
					name:    "t1",
					options: []string{"t2"},
				},
			},
		},
		{
			"tag with escaped quote-mark",
			"Field",
			`tag:"t1,t2\""`,
			[]Tag{
				{
					field:   "Field",
					key:     "tag",
					name:    "t1",
					options: []string{"t2\""},
				},
			},
		},
		{
			"multiple tags",
			"Field",
			`tag:"t1,t2" other:"o1,o2"`,
			[]Tag{
				{
					field:   "Field",
					key:     "tag",
					name:    "t1",
					options: []string{"t2"},
				},
				{
					field:   "Field",
					key:     "other",
					name:    "o1",
					options: []string{"o2"},
				},
			},
		},
		{
			"later tag overrides previous one",
			"Field",
			`tag:"t1,t2" other:"o1,o2" tag:"t3,t4"`,
			[]Tag{
				{
					field:   "Field",
					key:     "tag",
					name:    "t3",
					options: []string{"t4"},
				},
				{
					field:   "Field",
					key:     "other",
					name:    "o1",
					options: []string{"o2"},
				},
			},
		},
		{
			"use field name",
			"Field",
			`tag:",t1,t2" other:",o1,o2"`,
			[]Tag{
				{
					field:   "Field",
					key:     "tag",
					name:    "Field",
					options: []string{"t1", "t2"},
				},
				{
					field:   "Field",
					key:     "other",
					name:    "Field",
					options: []string{"o1", "o2"},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have, err := ParseTags(tc.fieldName, tc.fieldTag)

			// --- Then ---
			assert.NoError(t, err)
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_ParseTags_errors_tabular(t *testing.T) {
	tt := []struct {
		testN string

		tags string
	}{
		{"1", `tag`},
		{"2", `tag:`},
		{"3", `tag:\n`},
		{"4", `tag: "\n"`},
		{"5", ` tag`},
		{"6", `tag:"a",`},
		{"7", `tag:"`},
		{"8", `tag:"""`},
		{"9", `tag:"\19"`},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have, err := ParseTags("Field", tc.tags)

			// --- Then ---
			assert.ErrorIs(t, ErrTagSyntax, err)
			assert.Empty(t, have)
		})
	}
}

func Test_indirect(t *testing.T) {
	tt := []struct {
		testN string

		v any
	}{
		{"pointer to struct", &bytes.Buffer{}},
		{"slice", []string{"abc"}},
		{"map", map[string]string{"key": "abc"}},
		{"string", "abc"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := indirect(reflect.TypeOf(tc.v))

			// --- Then ---
			assert.NotEqual(t, reflect.Ptr, have.Kind())
		})
	}
}
