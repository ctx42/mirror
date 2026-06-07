This file provides guidance to AI agents (and human contributors) when working
with code in this repository.

## Commands

```bash
# Run all tests with race detector
go test -v -race ./...

# Run a single test
go test -run TestName ./pkg/mirror/

# Run a single test with race detector
go test -v -race -run TestName ./pkg/mirror/

# Lint (config lives in tmp/.golangci.yml, not the project root)
golangci-lint run --config tmp/.golangci.yml ./...
```

## Architecture

`mirror` is a single-package library (`pkg/mirror`) that provides cached struct
reflection. The key design split is between **type-level metadata** (cached,
immutable, shareable) and **value-level wrappers** (bound to a concrete
instance).

### Type-level layer (cached)

- **`Metadata`** — struct-level type info: `reflect.Type`, kind, package path,
  name, and slice of `*Field`. Created by `Reflect`, `ReflectType`, or
  `ReflectValue`. The global `typCache` (guarded by `sync.RWMutex`) stores one
  `*Metadata` per `reflect.Type`.
- **`Field`** — field-level type info, embedded inside `Metadata.fields`. Embeds
  `*Metadata` (via private alias `metadata`) because a field's type is itself a
  type. Holds `reflect.StructField`, kind, tag slice, and predicated booleans (
  slice, map, etc.). Constructed once during `Metadata.getFields()`.
- **`Tag`** — parsed representation of a single struct tag key (
  `key:"name,opt1,opt2"`). Carried in `Field.tags`, accessed via
  `Field.Tag(key)`.

### Value-level layer (per-instance)

- **`StructValue`** — wraps a `*Metadata` plus a `reflect.Value` for a concrete
  pointer-to-struct. Obtained via `NewStructValue(s)`. Returns `*FieldValue`
  from `FieldByName`/`FieldByIndex`.
- **`FieldValue`** — wraps a `*Field` plus a `reflect.Value` for a concrete
  field. Supports `Get()`, `Value()`, `NewIfNil()` (initialises nil
  pointers/maps/slices), and `StructValue()` for descending into nested structs.

Both `metadata` and `field` are private type aliases for `Metadata` and `Field`
respectively, used to embed without exposing the embedded type name in the
public API.

### Helpers

`helpers.go` contains `ParseTags` (the tag parser), `Ptr[T]` (generic pointer
helper), and internal utilities (`indirect`, `funcPkg`, `splitOnLastPeriod`).

### Testing

Tests use `github.com/ctx42/testing` (the org's own assertion library). Shared
test fixtures (`TStruct`, `TwoStr`, `NewTStruct`, `ptr`) live in `all_test.go`
and are available to all tests in the package. Example functions are in
`examples_test.go` and serve as both documentation and runnable tests.
