// Code generated by "stringer -type=Emitted -output=base_event.str.go -linecomment -withfromstring"; DO NOT EDIT.

package testdata

import (
	"strconv"
)

type Day int

const (
	Monday Day = iota
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