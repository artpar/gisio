package mtime

import (
	"testing"
	//"fmt"
)

func TestTimeParse(t *testing.T) {
	times := []string{
		"26 January 2016",
		"00",
	}
	for _, time := range times {
		theTime, format, _ := GetDate(time)
		{
			t.Logf("%s parsed from [%s] to %v", time, format, theTime)
		}
	}
}
