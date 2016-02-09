package mtime

import "testing"

func TestTimeParse(t *testing.T) {
	times := []string{
		"26 January 2016",
	}
	for _, time := range times {
		theTime, format, err := GetTime(time)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%s parsed from [%s] to %v", time, format, theTime)
	}
}
