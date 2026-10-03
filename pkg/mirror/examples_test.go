// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror_test

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ctx42/mirror/pkg/mirror"
)

func ExampleReflect() {
	s := &struct {
		F1 int
		F2 bool
		F3 string
		f4 time.Time
	}{}

	smd := mirror.Reflect(s)

	fmt.Printf("type: %v\n", smd.Type().String())
	fmt.Printf("kind: %v\n", smd.Kind().String())
	fmt.Printf("number of fields: %d\n", len(smd.Fields()))
	fmt.Printf("field by index: %s\n", smd.FieldByIndex(1).Name())
	fmt.Printf("field by name: %s\n", smd.FieldByName("f4").Name())

	// Output:
	// type: struct { F1 int; F2 bool; F3 string; f4 time.Time }
	// kind: struct
	// number of fields: 4
	// field by index: F2
	// field by name: f4
}

func ExampleReflectValue() {
	val := reflect.ValueOf(strings.TrimSpace)
	md := mirror.ReflectValue(val)

	fmt.Printf("type   : %s\n", md.Type().String())
	fmt.Printf("kind   : %s\n", md.Kind().String())
	fmt.Printf("name   : %s\n", md.Name())
	fmt.Printf("package: %s\n", md.Package())
	// Output:
	// type   : func(string) string
	// kind   : func
	// name   : TrimSpace
	// package: strings
}

func ExampleReflect_field() {
	s := &struct{ f4 time.Time }{}

	smd := mirror.Reflect(s)
	field := smd.FieldByName("f4")
	fmt.Printf("f4 type: %v\n", field.Type().String())
	fmt.Printf("f4 kind: %v\n", field.Kind().String())
	fmt.Printf("f4 index: %v\n", field.Index())
	fmt.Printf("f4 name: %v\n", field.Name())
	fmt.Printf("f4 valid: %v\n", field.IsValid())
	fmt.Printf("f4 exported: %v\n", field.IsExported())
	fmt.Printf("f4 slice: %v\n", field.IsSlice())
	fmt.Printf("f4 array: %v\n", field.IsArray())
	fmt.Printf("f4 slice or array: %v\n", field.IsSliceOrArray())
	fmt.Printf("f4 map: %v\n", field.IsMap())
	fmt.Printf("f4 interface: %v\n", field.IsInterface())
	fmt.Printf("f4 anonymous: %v\n", field.IsAnonymous())

	// Output:
	// f4 type: time.Time
	// f4 kind: struct
	// f4 index: [0]
	// f4 name: f4
	// f4 valid: true
	// f4 exported: false
	// f4 slice: false
	// f4 array: false
	// f4 slice or array: false
	// f4 map: false
	// f4 interface: false
	// f4 anonymous: false
}

func ExampleReflect_tag() {
	s := &struct {
		F1 int `my:"t1,t2, t3"`
	}{}

	smd := mirror.Reflect(s)
	field := smd.FieldByName("F1")
	tag, _ := field.Tag("my")

	fmt.Printf("F1 tag `my` key: %s\n", tag.Key())
	fmt.Printf("F1 tag `my` name: %s\n", tag.Name())
	fmt.Printf("F1 tag `my` ignored: %v\n", tag.IsIgnored())

	// Output:
	// F1 tag `my` key: my
	// F1 tag `my` name: t1
	// F1 tag `my` ignored: false
}

func ExampleStructValue_set() {
	s := &struct {
		F1 *int
	}{}

	smd := mirror.NewStructValue(s)
	field := smd.FieldByName("F1")
	value := field.NewIfNil().Value()

	value.Set(reflect.ValueOf(mirror.Ptr(42)))

	fmt.Printf("F1 value: %d\n", *s.F1)
	// Output:
	// F1 value: 42
}

func ExampleFieldValue_Get() {
	s := &struct {
		F1 int
	}{
		F1: 42,
	}

	smd := mirror.NewStructValue(s)
	field := smd.FieldByName("F1")
	value, _ := field.Get()

	fmt.Printf("F1 value: %v\n", value)
	// Output:
	// F1 value: 42
}

func ExampleReflectType() {
	typ := reflect.TypeOf(&struct{ F int }{})
	md := mirror.ReflectType(typ)

	fmt.Printf("kind: %s\n", md.Kind())
	fmt.Printf("fields: %d\n", len(md.Fields()))
	// Output:
	// kind: struct
	// fields: 1
}

func ExampleParseTags() {
	tags, _ := mirror.ParseTags("F", `json:"name,omitempty"`)

	fmt.Printf("key: %s\n", tags[0].Key())
	fmt.Printf("name: %s\n", tags[0].Name())
	fmt.Printf("omitempty: %t\n", tags[0].Contains("omitempty"))
	// Output:
	// key: json
	// name: name
	// omitempty: true
}

func ExampleFieldValue_NewIfNil() {
	s := &struct{ F *int }{}
	fv := mirror.NewStructValue(s).FieldByName("F")
	fv.NewIfNil()

	fmt.Printf("nil: %t\n", s.F == nil)
	fmt.Printf("zero: %d\n", *s.F)
	// Output:
	// nil: false
	// zero: 0
}

func ExampleStructValue_NewIfNil() {
	s := &struct{ F *int }{}
	sv := mirror.NewStructValue(s).FieldByName("F").StructValue()
	sv.NewIfNil()

	fmt.Printf("nil: %t\n", s.F == nil)
	fmt.Printf("zero: %d\n", *s.F)
	// Output:
	// nil: false
	// zero: 0
}

func ExampleField_Tag() {
	s := &struct {
		Age int `json:age`
	}{}
	fld := mirror.Reflect(s).FieldByName("Age")
	_, err := fld.Tag("json")

	fmt.Println(err)
	// Output:
	// struct field tag syntax error
}
