[![Go Report Card](https://goreportcard.com/badge/github.com/ctx42/mirror)](https://goreportcard.com/report/github.com/ctx42/mirror)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/ctx42/mirror)
![Tests](https://github.com/ctx42/mirror/actions/workflows/go.yml/badge.svg?branch=master)

<!-- TOC -->
* [Mirror: Cached Struct Reflection for Go](#mirror-cached-struct-reflection-for-go)
  * [Features](#features)
  * [Installation](#installation)
  * [Usage](#usage)
  * [Accessing Cached Struct](#accessing-cached-struct-)
  * [Accessing Cached Field](#accessing-cached-field)
  * [Accessing Cached Field Tags](#accessing-cached-field-tags)
  * [Setting Struct Fields](#setting-struct-fields)
  * [Getting Struct Field Value](#getting-struct-field-value)
  * [Accessing a Function or Method Value](#accessing-a-function-or-method-value)
<!-- TOC -->

# Mirror: Cached Struct Reflection for Go

`mirror` is a lightweight library that provides a simplified interface for 
reflecting structs. The metadata about the struct and its fields is cached to 
improve performance. 

## Features

- **Cached Reflection**: Parses struct metadata once and caches it for later fast access.
- **Simple Interface**: Provides an intuitive API to access struct and field metadata.
- **Field Manipulation**: Supports setting struct field values using cached metadata.
- **Tag Inspection**: Access and parse struct field tags.
- **Lightweight**: Minimal overhead with a focus on performance and simplicity.

## Installation

To use `mirror` in your Go project, install it using:

```bash
go get github.com/ctx42/mirror
```

## Usage

The `mirror` library provides two primary functions to access struct metadata:

```go
func Reflect(v any) *Metadata
func ReflectType(typ reflect.Type) *Metadata
```

These functions parse and cache struct metadata for faster later access.
Below are examples demonstrating common use cases.

## Accessing Cached Struct 

The `Reflect` function accepts a struct pointer and caches its metadata.
Subsequent calls with the same struct type retrieve the cached metadata,
improving performance.

<!-- gmdoceg:ExampleReflect -->
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

Field metadata is cached alongside the struct metadata, providing a simple
interface to inspect field properties.

<!-- gmdoceg:ExampleReflect_field -->
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

You can access and inspect struct field tags using the cached metadata.

<!-- gmdoceg:ExampleReflect_tag -->
```go
s := &struct {
    F1 int `my:"t1,t2, t3"`
}{}

smd := mirror.Reflect(s)
field := smd.FieldByName("F1")
tag := field.Tag("my")

fmt.Printf("F1 tag `my` key: %s\n", tag.Key())
fmt.Printf("F1 tag `my` name: %s\n", tag.Name())
fmt.Printf("F1 tag `my` ignored: %v\n", tag.IsIgnored())

// Output:
// F1 tag `my` key: my
// F1 tag `my` name: t1
// F1 tag `my` ignored: false
```

## Setting Struct Fields

The `mirror` library allows you to set struct field values fast by using cached
metadata, with support for initializing nil pointer fields.

<!-- gmdoceg:ExampleStructValue_set -->
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

<!-- gmdoceg:ExampleFieldValue_Get -->
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

## Accessing a Function or Method Value

`ReflectValue` resolves package path and name for any function or method value,
including anonymous closures.

<!-- gmdoceg:ExampleReflectValue -->
```go
val := reflect.ValueOf(check.After)
md := mirror.ReflectValue(val)

fmt.Printf("type   : %s\n", md.Type().String())
fmt.Printf("kind   : %s\n", md.Kind().String())
fmt.Printf("name   : %s\n", md.Name())
fmt.Printf("package: %s\n", md.Package())
// Output:
// type   : func(interface {}, interface {}, ...interface {}) error
// kind   : func
// name   : After
// package: github.com/ctx42/testing/pkg/check
```