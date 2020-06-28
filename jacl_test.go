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
		Key      string
		A        []interface{}
		B        []interface{}
		WantResp bool
		WantErr  error
	}{
		{"", []interface{}{AT{A: "a"}}, []interface{}{AT{A: "a"}}, true, nil},
		{"", []interface{}{AT{A: "a"}}, []interface{}{AT{A: "b"}}, false, nil},
		{"", []interface{}{AT{A: "a"}}, []interface{}{BT{A: "a", B: "b"}}, true, nil},
		{"a", []interface{}{AT{A: "a"}}, []interface{}{BT{A: "a", B: "b"}}, true, nil},
		{"a", []interface{}{AT{A: "a"}}, []interface{}{AT{A: "-"}, BT{A: "a", B: "b"}}, true, nil},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := Cmps(tc.Key, tc.A...)
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
