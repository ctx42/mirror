// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package mirror

// ptr returns a pointer to the given value.
//
// Example:
//
//	s := "Hello"
//	p := ptr(s) // p is a pointer to s
//
//	i := 10
//	q := ptr(i) // q is a pointer to i
//
//	f := 3.14
//	r := ptr(f) // r is a pointer to f
//
//	b := true
//	t := ptr(b) // t is a pointer to b
func ptr[T any](v T) *T { return &v }

// TwoStr is a struct with two string fields.
type TwoStr struct {
	FStr    string
	FStrPtr *string
}

// TStruct is a struct with multiple fields used for tests.
type TStruct struct {
	FStr  string `json:"f_json"`
	fStr  string
	FpStr *string `json:"-"`
	FsStr []string
	FaStr [4]string
	FmStr map[int]string
	SPtr  *TwoStr
	SVal  TwoStr
	SNil  *TwoStr
}

// NewTStruct returns TStruct with default values.
func NewTStruct() TStruct {
	FpStr := "TStruct.FpStr"
	PtrTwoStrFStrPtr := "ptr.TwoStr.FpStr"
	ValTwoStrFStrPtr := "val.TwoStr.FpStr"

	return TStruct{
		FStr:  "FStr",
		FpStr: &FpStr,
		FsStr: []string{"0", "1", "2"},
		FaStr: [4]string{"0", "1", "2", "3"},
		FmStr: map[int]string{1: "v1", 3: "vs"},
		fStr:  "fStr",
		SPtr: &TwoStr{
			FStr:    "ptr.TwoStr.FStr",
			FStrPtr: &PtrTwoStrFStrPtr,
		},
		SVal: TwoStr{
			FStr:    "val.TwoStr.FStr",
			FStrPtr: &ValTwoStrFStrPtr,
		},
		SNil: nil,
	}
}
