// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// ErrTagSyntax represents error when parsing struct field tag.
var ErrTagSyntax = errors.New("struct field tag syntax error")

// ParseTags parses a single struct field tag and returns a map of tags where
// keys are tag names. It returns the [ErrTagSyntax] error when a string
// describing tags is invalid.
//
// Example tags:
//
//	`url:"f1"`
//	`url:"f2,required"`
//	`custom:"f3,required"`
//	`custom:"f4,required, other"`
//	`custom:",required,other"`
//
// It will remove whitespace from options.
//
// If the field has multiple tags with the same key, the later one overrides
// the previous one.
//
// This is a modification of the Parse function from the repository:
// https://github.com/fatih/structtag
//
// nolint: gocognit, cyclop
func ParseTags(fieldName, stag string) ([]Tag, error) {
	hasTag := stag != ""
	if !hasTag {
		return nil, nil
	}

	var tags []Tag

	// NOTE(arslan) following code is from the reflect and vet package with
	// some modifications to collect all necessary information and extend it
	// with usable methods
	for stag != "" {
		// Skip leading space.
		i := 0
		for i < len(stag) && stag[i] == ' ' {
			i++
		}
		stag = stag[i:]
		if stag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax
		// error. Strictly speaking, control chars include the range [0x7f,
		// 0x9f], not just [0x00, 0x1f], but in practice, we ignore the
		// multibyte control characters as it is simpler to inspect the tag's
		// bytes than the tag's runes.
		i = 0
		for i < len(stag) && stag[i] > ' ' && stag[i] != ':' &&
			stag[i] != '"' && stag[i] != 0x7f {

			i++
		}

		if i == 0 {
			return nil, ErrTagSyntax
		}
		if i+1 >= len(stag) || stag[i] != ':' {
			return nil, ErrTagSyntax
		}
		if stag[i+1] != '"' {
			return nil, ErrTagSyntax
		}

		key := stag[:i]
		stag = stag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(stag) && stag[i] != '"' {
			if stag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(stag) {
			return nil, ErrTagSyntax
		}

		qvalue := stag[:i+1]
		stag = stag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return nil, ErrTagSyntax
		}

		var name string
		var options []string
		if value != "" {
			for i, opt := range strings.Split(value, ",") {
				opt = strings.TrimSpace(opt)
				if i == 0 && opt == "" {
					opt = fieldName
				}
				if opt == "" {
					continue
				}
				options = append(options, opt)
			}

			switch len(options) {
			case 0:
				options = nil
			case 1:
				name = options[0]
				options = nil

			default:
				name = options[0]
				options = options[1:]
			}
		}

		tag := Tag{
			field:   fieldName,
			key:     key,
			name:    name,
			options: options,
		}

		overwritten := false
		for i := 0; i < len(tags); i++ {
			if tags[i].key == tag.key {
				tags[i] = tag
				overwritten = true
			}
		}
		if !overwritten {
			if tags == nil {
				tags = make([]Tag, 0, 2)
			}
			tags = append(tags, tag)
		}
	}

	if len(tags) == 0 {
		return nil, nil
	}
	return tags, nil
}

// indirect returns the value that typ points to.
// The original typ is returned if typ is not a pointer.
func indirect(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}
