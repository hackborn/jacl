package jacl

import (
	"fmt"
	"testing"
)

// ------------------------------------------------------------
// TEST-COMPARE

func TestCompare(t *testing.T) {
	cases := []struct {
		A        interface{}
		B        interface{}
		WantResp bool
	}{
		{[]interface{}{"a", "b"}, []interface{}{"a", "b"}, true},
	}
	for i, tc := range cases {
		if !WantTestCase(i) {
			continue
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			haveResp := compare(tc.A, tc.B)
			if haveResp != tc.WantResp {
				fmt.Printf("have %v want %v\n", toJson(tc.B), toJson(tc.A))
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-NIL-CMP

func TestNilCmp(t *testing.T) {
	cases := []struct {
		B       interface{}
		WantErr error
	}{
		{nil, nil},
		{"a", cmpErr},
	}
	for i, tc := range cases {
		if !WantTestCase(i) {
			continue
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := CmpNil()
			haveErr := c.Cmp(tc.B)
			if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-SINGLE-CMP

func TestSingleCmp(t *testing.T) {
	cases := []struct {
		A       interface{}
		B       interface{}
		WantErr error
	}{
		{"a", "a", nil},
		{"a", "b", cmpErr},
		{AT{A: "a"}, AT{A: "a"}, nil},
		{AT{A: "a"}, AT{A: "b"}, cmpErr},
		{AT{A: "a"}, BT{A: "a"}, nil},
		{AT{A: "a"}, BT{A: "a", B: "b"}, nil},
	}
	for i, tc := range cases {
		if !WantTestCase(i) {
			continue
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := Cmp(tc.A)
			haveErr := c.Cmp(tc.B)
			if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-SLICE-CMP

func TestSliceCmp(t *testing.T) {
	cases := []struct {
		A       []interface{}
		B       []interface{}
		WantErr error
	}{
		{[]interface{}{"a"}, []interface{}{"a"}, nil},
		{[]interface{}{"a", "b"}, []interface{}{"a", "b"}, nil},
		{[]interface{}{"a"}, []interface{}{"b"}, cmpErr},
		{[]interface{}{AT{A: "a"}}, []interface{}{AT{A: "a"}}, nil},
		{[]interface{}{AT{A: "a"}}, []interface{}{AT{A: "b"}}, cmpErr},
		{[]interface{}{AT{A: "a"}}, []interface{}{BT{A: "a", B: "b"}}, nil},
		{[]interface{}{BT{A: "d", B: "e"}, BT{A: "a", B: "b"}}, []interface{}{BT{A: "a", B: "c"}}, cmpErr},
	}
	for i, tc := range cases {
		if !WantTestCase(i) {
			continue
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := Cmps(tc.A...)
			haveErr := c.Cmp(tc.B)
			if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-SLICE-KEY

func TestSliceKey(t *testing.T) {
	cases := []struct {
		A       []interface{}
		B       []interface{}
		WantErr error
	}{
		{[]interface{}{Key("a"), AT{A: "a"}}, []interface{}{BT{A: "a", B: "b"}}, nil},
		{[]interface{}{Key("a"), AT{A: "a"}}, []interface{}{AT{A: "-"}, BT{A: "a", B: "b"}}, nil},
		{[]interface{}{Key("a", "b"), BT{A: "a", B: "b"}}, []interface{}{BT{A: "a", B: "b"}}, nil},
		{[]interface{}{Key("a", "b"), BT{A: "a", B: "b"}}, []interface{}{BT{A: "a", B: "c"}}, cmpErr},
		// Everything in A must be in B, even with keys. This is testing unordered comparisons.
		{[]interface{}{Key("a", "b"), BT{A: "d", B: "e"}, BT{A: "a", B: "b"}}, []interface{}{BT{A: "a", B: "b"}}, cmpErr},
		// But B can have more than A. This is testing unordered comparisons.
		{[]interface{}{Key("a", "b"), BT{A: "a", B: "b"}}, []interface{}{BT{A: "d", B: "e"}, BT{A: "a", B: "b"}}, nil},
	}
	for i, tc := range cases {
		if !WantTestCase(i) {
			continue
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := Cmps(tc.A...)
			haveErr := c.Cmp(tc.B)
			if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-SLICE-NOTEXISTS

func TestSliceNotexists(t *testing.T) {
	cases := []struct {
		A       []interface{}
		B       []interface{}
		WantErr error
	}{
		{[]interface{}{NotExists("a")}, []interface{}{"a"}, &ComparisonError{}},
		{[]interface{}{NotExists("b")}, []interface{}{"a"}, nil},
		// Path to a generic object
		{[]interface{}{NotExists("a", "b")}, []interface{}{map[string]string{"a": "b"}}, &ComparisonError{}},
		{[]interface{}{NotExists("a", "b")}, []interface{}{map[string]string{"a": "c"}}, nil},
		{[]interface{}{NotExists("a")}, []interface{}{map[string]string{"a": "c"}}, &ComparisonError{}},
	}
	for i, tc := range cases {
		if !WantTestCase(i) {
			continue
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := Cmps(tc.A...)
			haveErr := c.Cmp(tc.B)
			if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-SLICE-SIZEIS

func TestSliceSizeis(t *testing.T) {
	cases := []struct {
		A       []interface{}
		B       []interface{}
		WantErr error
	}{
		{[]interface{}{SizeIs(1), "a"}, []interface{}{"a"}, nil},
		{[]interface{}{SizeIs(2), "a"}, []interface{}{"a"}, cmpErr},
		{[]interface{}{SizeIs(2), "a", "b"}, []interface{}{"a", "b"}, nil},
	}
	for i, tc := range cases {
		if !WantTestCase(i) {
			continue
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c := Cmps(tc.A...)
			haveErr := c.Cmp(tc.B)
			if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-SINGLE-CMPER-FACTORY

func TestSingleCmperFactory(t *testing.T) {
	cases := []struct {
		Req     singleCmp
		WantErr error
	}{
		{singleCmp{}, nil},
		{singleCmp{A: "a"}, nil},
		{singleCmp{A: BT{A: "a", B: "b"}}, nil},
		{singleCmp{A: []string{"a", "b"}}, nil},
		{singleCmp{A: []interface{}{"a", "b"}}, nil},
	}
	for i, tc := range cases {
		if !WantTestCase(i) {
			continue
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			input := CmperFactory{Cmper: tc.Req}
			output := CmperFactory{}
			err := toFromJson(input, &output)
			if err != nil {
				panic(err)
			}
			haveErr := output.Cmper.Cmp(tc.Req.A)
			outputKey := ""
			if s, ok := output.Cmper.(serializer); ok {
				outputKey = s.SerializeKey()
			}
			if tc.Req.SerializeKey() != outputKey {
				fmt.Printf("wrong key have %v want %v\n", outputKey, tc.Req.SerializeKey())
				t.Fatal()
			} else if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			}
		})
	}
}

// ------------------------------------------------------------
// TEST-SLICE-CMPER-FACTORY

/*
func TestSliceCmperFactory(t *testing.T) {
	cases := []struct {
		Req     sliceCmp
		WantErr error
	}{
		//		{sliceCmp{}, nil},
		{sliceCmp{}.addFn(NotExists("a")), nil},
	}
	for i, tc := range cases {
		if !WantTestCase(i) {
			continue
		}
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			input := CmperFactory{Cmper: tc.Req}
			output := CmperFactory{}
			err := toFromJson(input, &output)
			if err != nil {
				panic(err)
			}
			//			fmt.Println("BEFORE", input, "AFTER", output, "TYPE", output.Keys, "A", output.A, "FN", output.Fn)
			if sc, ok := output.Cmper.(*sliceCmp); ok {
				fmt.Println("BEFORE", input, "AFTER", sc)
			} else {
				fmt.Printf("TYPE IS %t\n", output.Cmper)
			}
			panic("sds")
			haveErr := output.Cmper.Cmp(tc.Req.A)
			outputKey := ""
			if s, ok := output.Cmper.(serializer); ok {
				outputKey = s.SerializeKey()
			}
			if tc.Req.SerializeKey() != outputKey {
				fmt.Printf("wrong key have %v want %v\n", outputKey, tc.Req.SerializeKey())
				t.Fatal()
			} else if !equalErr(haveErr, tc.WantErr) {
				fmt.Printf("have err %v want %v\n", haveErr, tc.WantErr)
				t.Fatal()
			}
		})
	}
}
*/

// ------------------------------------------------------------
// COMPARISON TYPES

type AT struct {
	A interface{} `json:"a,omitempty"`
}

type BT struct {
	A interface{} `json:"a,omitempty"`
	B interface{} `json:"b,omitempty"`
}

// ------------------------------------------------------------
// CONST and VAR

var (
	cmpErr = newComparisonError("")
)
