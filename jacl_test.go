package jacl

import (
	"fmt"
	"testing"
)

// ------------------------------------------------------------
// TEST-SINGLE-CMP

func TestSingleCmp(t *testing.T) {
	cases := []struct {
		A        interface{}
		B        interface{}
		WantResp bool
		WantErr  error
	}{
		{"a", "a", true, nil},
		{"a", "b", false, nil},
		{AT{A: "a"}, AT{A: "a"}, true, nil},
		{AT{A: "a"}, AT{A: "b"}, false, nil},
		{AT{A: "a"}, BT{A: "a"}, true, nil},
		{AT{A: "a"}, BT{A: "a", B: "b"}, true, nil},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := Cmp(tc.A)
			haveResp, haveErr := c.Cmp(tc.B)
			if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			} else if tc.WantResp != haveResp {
				fmt.Printf("have resp %v want match %v\n", haveResp, tc.WantResp)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-SLICE-CMP

func TestSliceCmp(t *testing.T) {
	cases := []struct {
		Keys     interface{}
		A        []interface{}
		B        []interface{}
		WantResp bool
		WantErr  error
	}{
		{"", []interface{}{"a"}, []interface{}{"a"}, true, nil},
		{"", []interface{}{"a", "b"}, []interface{}{"a", "b"}, true, nil},
		{"", []interface{}{"a"}, []interface{}{"b"}, false, nil},
		{"", []interface{}{AT{A: "a"}}, []interface{}{AT{A: "a"}}, true, nil},
		{"", []interface{}{AT{A: "a"}}, []interface{}{AT{A: "b"}}, false, nil},
		{"", []interface{}{AT{A: "a"}}, []interface{}{BT{A: "a", B: "b"}}, true, nil},
		{"a", []interface{}{AT{A: "a"}}, []interface{}{BT{A: "a", B: "b"}}, true, nil},
		{"a", []interface{}{AT{A: "a"}}, []interface{}{AT{A: "-"}, BT{A: "a", B: "b"}}, true, nil},
		{[]string{"a", "b"}, []interface{}{BT{A: "a", B: "b"}}, []interface{}{BT{A: "a", B: "b"}}, true, nil},
		{[]string{"a", "b"}, []interface{}{BT{A: "a", B: "b"}}, []interface{}{BT{A: "a", B: "c"}}, false, nil},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := Cmps(tc.Keys, tc.A...)
			haveResp, haveErr := c.Cmp(tc.B)
			if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			} else if tc.WantResp != haveResp {
				fmt.Printf("have resp %v want match %v\n", haveResp, tc.WantResp)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-EQUAL

func TestEqual(t *testing.T) {
	cases := []struct {
		Keys     interface{}
		A        []interface{}
		B        []interface{}
		WantResp bool
		WantErr  bool
	}{
		{"", []interface{}{SizeIs(1), "a"}, []interface{}{"a"}, true, false},
		{"", []interface{}{SizeIs(2), "a"}, []interface{}{"a"}, false, true},
		{"", []interface{}{SizeIs(2), "a", "b"}, []interface{}{"a", "b"}, true, false},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := Cmps(tc.Keys, tc.A...)
			haveResp, haveErr := c.Cmp(tc.B)
			if checkErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			} else if tc.WantResp != haveResp {
				fmt.Printf("have resp %v want match %v\n", haveResp, tc.WantResp)
				t.Fatal()
			}
		})
	}
}

func checkErr(have error, want bool) bool {
	switch want {
	case true:
		return have == nil
	default:
		return have != nil
	}
}

// ------------------------------------------------------------
// TEST-CMPER-FACTORY

func TestSingleCmperFactory(t *testing.T) {
	cases := []struct {
		Req     singleCmp
		WantErr error
	}{
		{singleCmp{}, nil},
		{singleCmp{A: "a"}, nil},
		{singleCmp{A: BT{A: "a", B: "b"}}, nil},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			input := CmperFactory{Cmper: tc.Req}
			output := CmperFactory{}
			err := toFromJson(input, &output)
			if err != nil {
				panic(err)
			}
			haveResp, haveErr := output.Cmper.Cmp(tc.Req.A)
			if tc.Req.FactoryKey() != output.Cmper.FactoryKey() {
				fmt.Printf("wrong key have %v want %v\n", output.Cmper.FactoryKey(), tc.Req.FactoryKey())
				t.Fatal()
			} else if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			} else if haveResp != true {
				fmt.Printf("failed want %v\n", tc.Req.A)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// COMPARISON TYPES

type AT struct {
	A interface{} `json:"a,omitempty"`
}

type BT struct {
	A interface{} `json:"a,omitempty"`
	B interface{} `json:"b,omitempty"`
}

type CT struct {
	A interface{} `json:"a,omitempty"`
	B interface{} `json:"b,omitempty"`
	C interface{} `json:"c,omitempty"`
}
