package testdata

import (
	"fmt"
	ftm "fmt"
	"strconv"
)

type Level int

const (
	Low Level = iota << 2
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
