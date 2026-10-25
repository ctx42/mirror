// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

import (
	"iter"
	"reflect"
	"slices"
)

// Names maps the names a struct tag key gives the fields of a struct, and of
// the structs it embeds, to the fields Go promotion selects, as
// [encoding/json] does: the shallowest field with the name or, among equally
// shallow ones, the only one whose tag names it; when there is no such field,
// the name selects none. A field is named by its tag name, else by its Go
// name. Ignored fields (tag name "-") and unexported fields name nothing,
// embedded interfaces are skipped, and embedded structs, or pointers to them,
// are walked, exported or not. The names only ignored fields carry, which no
// other field takes or shares, are recorded as ignored. A struct type met at
// a shallower depth is not walked again, so a recursive type ends.
//
// Build it with [Metadata.Names]; it is read-only and safe for concurrent
// use.
type Names struct {
	fields  map[string][]int // Index path of the field each name selects.
	ignored map[string]bool  // Names carried only by ignored fields.
}

// nameNode is a struct walked by newNames.
type nameNode struct {
	md  *Metadata // The struct.
	idx []int     // Its index path from the struct names are built for.
}

// nameCand is a field that may take a name.
type nameCand struct {
	idx    []int // Index path of the field.
	tagged bool  // The struct tag names the field.
}

// newNames returns the names the key struct tag gives the fields of the md
// struct.
func newNames(md *Metadata, key string) *Names {
	nms := &Names{fields: map[string][]int{}, ignored: map[string]bool{}}
	done := map[string]bool{}       // Names decided at a shallower depth.
	seen := map[reflect.Type]bool{} // Struct types of shallower depths.
	for next := []nameNode{{md: md}}; len(next) > 0; {
		var cands map[string][]nameCand
		cands, next = nms.depth(next, key, seen)
		nms.pick(cands, done)
	}
	for name := range done {
		delete(nms.ignored, name) // A name some field takes or several share.
	}
	return nms
}

// depth returns the fields that may take the names the key struct tag gives
// the fields of the cur structs, which are of one depth, and the structs they
// embed, of the next depth. The struct types in seen, met at a shallower
// depth, are skipped, and the types of cur are added to it.
func (nms *Names) depth(
	cur []nameNode,
	key string,
	seen map[reflect.Type]bool,
) (map[string][]nameCand, []nameNode) {

	cands := map[string][]nameCand{}
	walked := map[reflect.Type]bool{}
	var next []nameNode
	for _, nod := range cur {
		if seen[nod.md.typ] {
			continue
		}
		walked[nod.md.typ] = true
		next = nms.fieldsOf(nod, key, cands, next)
	}
	for typ := range walked {
		seen[typ] = true
	}
	return cands, next
}

// fieldsOf adds the fields of the nod struct that may take a name to cands,
// records the names only ignored fields carry, and returns next with the
// structs nod embeds appended.
func (nms *Names) fieldsOf(
	nod nameNode,
	key string,
	cands map[string][]nameCand,
	next []nameNode,
) []nameNode {

	for i, fld := range nod.md.fields {
		tag, _ := fld.Tag(key)
		isStruct := fld.metadata.kind == reflect.Struct
		if tag.IsIgnored() {
			// An ignored embedded struct names nothing of its own.
			if !fld.anonymous || !isStruct {
				nms.ignored[tag.NameOrField()] = true
			}
			continue
		}
		idx := append(slices.Clip(nod.idx), i)
		switch {
		case fld.anonymous && fld.IsInterface():
			continue

		case fld.anonymous && isStruct:
			next = append(next, nameNode{md: fld.metadata, idx: idx})
			continue

		case !fld.IsExported():
			continue
		}
		name := tag.NameOrField()
		cands[name] = append(cands[name], nameCand{idx, tag.name != ""})
	}
	return next
}

// pick records the field each name of cands, all of one depth, selects: the
// only field taking it, or the only tagged one among several. A name in done,
// decided at a shallower depth, is skipped; the names of cands are added to
// done.
func (nms *Names) pick(cands map[string][]nameCand, done map[string]bool) {
	for name, cds := range cands {
		if done[name] {
			continue
		}
		done[name] = true
		if len(cds) == 1 {
			nms.fields[name] = cds[0].idx
			continue
		}
		var tagged []nameCand
		for _, cnd := range cds {
			if cnd.tagged {
				tagged = append(tagged, cnd)
			}
		}
		if len(tagged) == 1 {
			nms.fields[name] = tagged[0].idx
		}
	}
}

// Field returns the index path, for [reflect.Value.FieldByIndex], of the field
// the name selects, and false when it selects none. The path is shared:
// the caller must not modify it.
func (nms *Names) Field(name string) ([]int, bool) {
	idx, ok := nms.fields[name]
	return idx, ok
}

// Selects returns true if the name selects the field at index i of the struct
// at the index path at, as given by [Names.Field] for an embedded struct; an
// empty at is the struct itself.
func (nms *Names) Selects(name string, at []int, i int) bool {
	idx, ok := nms.fields[name]
	if !ok || len(idx) != len(at)+1 || idx[len(at)] != i {
		return false
	}
	return slices.Equal(idx[:len(at)], at)
}

// All returns an iterator over the names that select a field and the index
// paths, as [Names.Field] returns them, of the fields they select, in no
// particular order. The paths are shared: the caller must not modify them.
func (nms *Names) All() iter.Seq2[string, []int] {
	return func(yield func(string, []int) bool) {
		for name, idx := range nms.fields {
			if !yield(name, idx) {
				return
			}
		}
	}
}

// IsIgnored returns true if only ignored fields carry the name.
func (nms *Names) IsIgnored(name string) bool { return nms.ignored[name] }

// Len returns the number of names that select a field.
func (nms *Names) Len() int { return len(nms.fields) }
