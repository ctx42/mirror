// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mirror

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
