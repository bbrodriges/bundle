package main

import (
	"bytes"
	"testing"
)

type testCase struct {
	caseName    string
	packageName string
	filesArg    []string
	expecting   string
	hasError    bool
}

var testCases = []testCase{
	{"all", "bundle", []string{"testdata/*.go"}, all_bundle_out, false},
	{"specificFilesPattern", "bundle", []string{"testdata/*s.go"}, specificFilesPattern_bundle_out, false},
	{"specificFilesList", "bundle", []string{"testdata/days.go", "testdata/number.go"}, specificFilesList_bundle_out, false},
	{"notAFile", "badbundle", []string{"testdata/"}, "", true},
	{"singleFile", "badbundle", []string{"testdata/days.go"}, "", true},
}

var all_bundle_out = `// Code generated by bundle generation tool; DO NOT EDIT.

package bundle

import (
	"fmt"
	ftm "fmt"
	"math/rand"
	"strconv"
)

// The code below has been bundled from "testdata/days.go" source file.
// Code generated by "stringer -type=Emitted -output=base_event.str.go -linecomment -withfromstring"; DO NOT EDIT.

type Day int

const (
	Monday	Day	= iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (d Day) Print() string {
	return strconv.Itoa(int(d))
}

// The code below has been bundled from "testdata/enums.go" source file.

type Level int

const (
	Low	Level	= iota << 2
	Medium
	High
)

func (l Level) Print() string {
	return strconv.Itoa(int(l))
}

func (l Level) Printf() (int, error) {
	return ftm.Printf("")
}

func (l Level) String() string {
	switch l {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// The code below has been bundled from "testdata/number.go" source file.

func GetRandomNumber() int {
	return rand.Intn(1000)
}

`

var specificFilesPattern_bundle_out = `// Code generated by bundle generation tool; DO NOT EDIT.

package bundle

import (
	"fmt"
	ftm "fmt"
	"strconv"
)

// The code below has been bundled from "testdata/days.go" source file.
// Code generated by "stringer -type=Emitted -output=base_event.str.go -linecomment -withfromstring"; DO NOT EDIT.

type Day int

const (
	Monday	Day	= iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (d Day) Print() string {
	return strconv.Itoa(int(d))
}

// The code below has been bundled from "testdata/enums.go" source file.

type Level int

const (
	Low	Level	= iota << 2
	Medium
	High
)

func (l Level) Print() string {
	return strconv.Itoa(int(l))
}

func (l Level) Printf() (int, error) {
	return ftm.Printf("")
}

func (l Level) String() string {
	switch l {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

`

var specificFilesList_bundle_out = `// Code generated by bundle generation tool; DO NOT EDIT.

package bundle

import (
	"math/rand"
	"strconv"
)

// The code below has been bundled from "testdata/days.go" source file.
// Code generated by "stringer -type=Emitted -output=base_event.str.go -linecomment -withfromstring"; DO NOT EDIT.

type Day int

const (
	Monday	Day	= iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (d Day) Print() string {
	return strconv.Itoa(int(d))
}

// The code below has been bundled from "testdata/number.go" source file.

func GetRandomNumber() int {
	return rand.Intn(1000)
}

`

func TestGenerator(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			g := generator{
				packageName:   tc.packageName,
				removeSources: false,
				filesArgs:     tc.filesArg,
			}

			w := bytes.NewBuffer(nil)
			err := g.makeBundle(w)

			if err != nil {
				if tc.hasError {
					return
				}
				t.Fatalf("Unexpected error: %s", err)
			}

			if !bytes.Equal(w.Bytes(), []byte(tc.expecting)) {
				t.Fatalf("Outputs mismatch.\n\nExpecting:\n%+v\n\n================\n\nGot:\n%+v", tc.expecting, w.String())
			}
		})
	}
}
