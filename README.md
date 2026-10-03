[![Tests](https://github.com/ctx42/mirror/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/ctx42/mirror/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ctx42/mirror.svg)](https://pkg.go.dev/github.com/ctx42/mirror)
[![Go Version](https://img.shields.io/badge/go-1.26-00ADD8)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE.md)

# Mirror

Cached struct reflection for Go programs that read and write struct fields.

It parses a struct once and keeps the metadata, so later reads of the same
type skip that work.

<!-- TOC -->
* [Mirror](#mirror)
  * [Features](#features)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
  * [Usage](#usage)
  * [Accessing Cached Struct](#accessing-cached-struct)
  * [Accessing Cached Field](#accessing-cached-field)
  * [Accessing Cached Field Tags](#accessing-cached-field-tags)
  * [Setting Struct Fields](#setting-struct-fields)
  * [Getting Struct Field Value](#getting-struct-field-value)
  * [Functions and Methods](#functions-and-methods)
<!-- TOC -->

## Features

- **Cached reflection**: Parse struct metadata once and reuse it.
- **Field values**: Read and set struct fields through the cached metadata.
- **Tag inspection**: Read and parse struct field tags.

## Prerequisites

Go 1.26 or later.

## Installation

```bash
go get github.com/ctx42/mirror
```

Import the package:

```go
import "github.com/ctx42/mirror/pkg/mirror"
```

## Usage

`Reflect`, `ReflectType`, and `ReflectValue` read metadata and cache it.

```go
func Reflect(v any) *Metadata
func ReflectType(typ reflect.Type) *Metadata
func ReflectValue(val reflect.Value) *Metadata
```

## Accessing Cached Struct

`Reflect` accepts any non-nil value. It follows every pointer and caches the
metadata for that type. A later call with the same type returns the cached
metadata. An untyped nil panics.

<!-- gmmce:pkg/mirror/ExampleReflect -->
```go
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
```

## Accessing Cached Field

Field metadata is cached with the struct, and each field reports its type,
name, and kind.

<!-- gmmce:pkg/mirror/ExampleReflect_field -->
```go
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
```

## Accessing Cached Field Tags

`Tag` returns the tag for a key. A missing key is a zero tag and a nil
error. A tag string that failed to parse returns a zero tag and
`ErrTagSyntax`.

<!-- gmmce:pkg/mirror/ExampleReflect_tag -->
```go
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
```

## Setting Struct Fields

`NewStructValue` wraps a struct pointer. `NewIfNil` initializes a nil pointer
field before the value is set.

<!-- gmmce:pkg/mirror/ExampleStructValue_set -->
```go
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
```

## Getting Struct Field Value

`Get` returns the field's current value.

<!-- gmmce:pkg/mirror/ExampleFieldValue_Get -->
```go
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
```

## Functions and Methods

`ReflectValue` keeps the runtime name of a function or method, including an
anonymous closure.

<!-- gmmce:pkg/mirror/ExampleReflectValue -->
```go
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
```
