package jacl

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type WantTestOpts struct {
	Filename string // The settings file name.
}

// WantTest answers true if the current test should run.
// If no tests are configured, all tests run. Tests are configured via
// a settings file, by default located at testdata/settings.json.
//
// This is a utility function ancillary to the main package. It's included
// because it is being used as part of the testing, and I exposed it because
// it's handy for other clients, also.
func WantTest() bool {
	return WantTestCase(-1)
}

// WantTestCase answers true if the current test case should run.
// If no tests are configured, all tests run. Tests are configured via
// a settings file, by default located at testdata/settings.json.
//
// This is a utility function ancillary to the main package. It's included
// because it is being used as part of the testing, and I exposed it because
// it's handy for other clients, also.
func WantTestCase(caseIndex int) bool {
	return WantTestCaseWith(caseIndex, WantTestOpts{})
}

// WantTestCaseWith answers true if the current test case should run.
// If no tests are configured, all tests run. Tests are configured via
// a settings file, by default located at testdata/settings.json.
//
// This is a utility function ancillary to the main package. It's included
// because it is being used as part of the testing, and I exposed it because
// it's handy for other clients, also.
func WantTestCaseWith(caseIndex int, opts WantTestOpts) bool {
	// Load desired tests from settings
	tests := loadTests(opts)
	if len(tests) < 1 {
		return true
	}

	// Match desired test and index against something in my
	// callstack and requested index.
	fn := func(a string) bool {
		idx, ok := tests[a]
		if !ok {
			return false
		}
		return caseIndex < 0 || idx < 0 || caseIndex == idx
	}
	return hasFunctionName(fn)
}

// loadTests() loads my test list from the settings, answering
// a map of test names and the desired index for each.
func loadTests(opts WantTestOpts) map[string]int {
	fn := opts.Filename
	if fn == "" {
		fn = filepath.Join("./testdata", "settings.json")
	}
	contents, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil
	}

	type Settings struct {
		WantTests []string `json:"want_tests,omitempty"`
	}
	settings := Settings{}
	err = json.Unmarshal(contents, &settings)
	if err != nil {
		return nil
	}
	tests := make(map[string]int)
	for _, t := range settings.WantTests {
		if t != "" {
			n, idx := getTestName(t)
			tests[n] = idx
		}
	}
	return tests
}

type hasStringFunc func(string) bool

// hasFunctionName() walks my call stack and compares
// the name of each frame with the function.
func hasFunctionName(fn hasStringFunc) bool {
	pc := make([]uintptr, 8)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, more := frames.Next()
	for more {
		name := getFrameName(frame)
		if fn(name) {
			return true
		}
		frame, more = frames.Next()
	}
	return false
}

func getFrameName(f runtime.Frame) string {
	all := strings.Split(f.Function, "/")
	if len(all) < 1 {
		return ""
	}
	last := strings.Split(all[len(all)-1], ".")
	if len(last) < 1 {
		return ""
	}
	return last[len(last)-1]
}

// getTestName() takes a string that might have an index
// and answers the string and index (or -1).
func getTestName(s string) (name string, idx int) {
	all := strings.Split(s, ":")
	if len(all) > 0 {
		name = all[0]
	}
	idx = -1
	if len(all) > 1 {
		i2, err := strconv.ParseInt(all[1], 10, 32)
		if err == nil {
			idx = int(i2)
		}
	}
	return
}
