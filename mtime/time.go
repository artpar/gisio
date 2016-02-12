package mtime

import (
	"time"
	"errors"
	"sort"
	//"fmt"
)

var timeFormat []string

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func init() {
	timeFormat = []string{
		"Mon Jan _2 15:04:05 2006",
		"Mon Jan _2 15:04:05 MST 2006",
		"Mon Jan 02 15:04:05 -0700 2006",
		"02 Jan 06 15:04 MST",
		"02 Jan 2006",
		"Jan 02, 2006",
		"02 January 2006",
		"January 02, 2006",
		"January 02",
		"Jan 02",
		"06",
		"02 Jan 06",
		"02 Jan 06 15:04 -0700",
		"Monday, 02-Jan-06 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"01021504",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.999999999Z07:00",
		"3:04PM",
		"Jan _2 15:04:05",
		"Jan _2 15:04:05.000",
		"Jan _2 15:04:05.000000",
		"Jan _2 15:04:05.000000000",
	}
	sort.Sort(ByLength(timeFormat))
}

func GetTime(t string) (time.Time, string, error) {
	for _, format := range timeFormat {
		// fmt.Printf("Testing %s with %s\n", t, format)
		t, err := time.Parse(format, t)
		if err == nil {
			return t, format, nil
		}
	}
	return time.Now(), "", errors.New("Unrecognised time format - " + t)
}

func GetTimeByFormat(t string, f string) (time.Time, error) {
	return time.Parse(f, t)
}