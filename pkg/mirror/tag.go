// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

import (
	"slices"
)

// Tag represents a single structure field tag.
//
// Example:
//
//	`key:"name,option0,option1"`
type Tag struct {
	field   string   // Struct field name the tag is attached to.
	key     string   // Tag key.
	name    string   // Tag name.
	options []string // Tag options.
}

// Key returns the key of the tag.
func (tag Tag) Key() string { return tag.key }

// Name returns the name of the tag.
func (tag Tag) Name() string { return tag.name }

// Contains returns true if the option exists on the list.
func (tag Tag) Contains(option string) bool {
	return slices.Contains(tag.options, option)
}

// NameOrField returns the name of the tag or the field name if the tag is
// empty or set to the "-" value.
func (tag Tag) NameOrField() string {
	if tag.name != "" && tag.name != "-" {
		return tag.name
	}
	return tag.field
}

// IsIgnored returns true if the tag name is set to the "-" value.
func (tag Tag) IsIgnored() bool { return tag.name == "-" }

// IsZero returns true if the Tag is empty (all fields are empty strings).
func (tag Tag) IsZero() bool {
	return tag.key == "" && tag.name == "" && len(tag.options) == 0
}
