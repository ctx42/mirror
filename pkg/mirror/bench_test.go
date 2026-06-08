// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

import (
	"testing"
)

// wide20 is a struct with 20 fields — exercises the O(n) FieldByName path.
type wide20 struct {
	F00, F01, F02, F03, F04 string
	F05, F06, F07, F08, F09 string
	F10, F11, F12, F13, F14 string
	F15, F16, F17, F18, F19 string
}

// BenchmarkReflect measures the cost of a cache-hit Reflect call.
func BenchmarkReflect(b *testing.B) {
	s := &TStruct{}
	Reflect(s) // prime the cache
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Reflect(s)
	}
}

// BenchmarkReflect_Concurrent measures cache-read throughput under concurrent
// access from multiple goroutines.
func BenchmarkReflect_Concurrent(b *testing.B) {
	s := &TStruct{}
	Reflect(s) // prime the cache
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Reflect(s)
		}
	})
}

// BenchmarkMetadata_FieldByName measures the linear scan cost of
// Metadata.FieldByName for a small and a large struct.
func BenchmarkMetadata_FieldByName(b *testing.B) {
	b.Run("small_2fields/hit_first", func(b *testing.B) {
		md := Reflect(&TwoStr{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = md.FieldByName("FStr")
		}
	})

	b.Run("small_2fields/hit_last", func(b *testing.B) {
		md := Reflect(&TwoStr{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = md.FieldByName("FStrPtr")
		}
	})

	b.Run("wide_20fields/hit_first", func(b *testing.B) {
		md := Reflect(&wide20{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = md.FieldByName("F00")
		}
	})

	b.Run("wide_20fields/hit_last", func(b *testing.B) {
		md := Reflect(&wide20{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = md.FieldByName("F19")
		}
	})

	b.Run("wide_20fields/miss", func(b *testing.B) {
		md := Reflect(&wide20{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = md.FieldByName("NoSuchField")
		}
	})
}

// BenchmarkStructValue_FieldByName measures the double-linear-search cost of
// StructValue.FieldByName, which calls both reflect.Value.FieldByName (O(n))
// and Metadata.FieldByName (O(n)) on every call.
func BenchmarkStructValue_FieldByName(b *testing.B) {
	b.Run("small_2fields/hit_first", func(b *testing.B) {
		sv := NewStructValue(&TwoStr{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sv.FieldByName("FStr")
		}
	})

	b.Run("small_2fields/hit_last", func(b *testing.B) {
		sv := NewStructValue(&TwoStr{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sv.FieldByName("FStrPtr")
		}
	})

	b.Run("wide_20fields/hit_first", func(b *testing.B) {
		sv := NewStructValue(&wide20{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sv.FieldByName("F00")
		}
	})

	b.Run("wide_20fields/hit_last", func(b *testing.B) {
		sv := NewStructValue(&wide20{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sv.FieldByName("F19")
		}
	})
}

// BenchmarkStructValue_FieldByIndex establishes the O(1) baseline for
// index-based lookup to compare against FieldByName.
func BenchmarkStructValue_FieldByIndex(b *testing.B) {
	b.Run("wide_20fields/first", func(b *testing.B) {
		sv := NewStructValue(&wide20{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sv.FieldByIndex(0)
		}
	})

	b.Run("wide_20fields/last", func(b *testing.B) {
		sv := NewStructValue(&wide20{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sv.FieldByIndex(19)
		}
	})
}
