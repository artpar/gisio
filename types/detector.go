package types

import (
	"github.com/artpar/gisio/mtime"
	"net"
	"strconv"
	"regexp"
	"errors"
	"fmt"
)

type EntityType int

func (t EntityType) String() string {
	switch t {
	case time:
		return "time"
	case ipaddress:
		return "ipaddress"
	case money:
		return "money"
	case number:
		return "number"
	case none:
		return "none"
	}
	return "failed-to-detect"
}

func (t EntityType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.String() + "\""), nil
}

const (
	time EntityType = iota
	ipaddress
	money
	number
	none
)

var detector map[EntityType]func(string) bool

func init() {
	detector = make(map[EntityType]func(string) bool)
	detector[time] = func(d string) bool {
		_, _, err := mtime.GetTime(d)
		if err == nil {
			return true
		}
		return false
	}
	detector[ipaddress] = func(d string) bool {
		s := net.ParseIP(d)
		if s != nil {
			return true
		}
		return false
	}
	detector[money] = func(d string) bool {
		r := regexp.MustCompile("^([a-zA-Z]{0,3}\\.?)?[0-9]+\\.[0-9]{0,2}([a-zA-Z]{0,3})?")
		return r.MatchString(d)
	}
	detector[number] = func(d string) bool {
		_, err := strconv.ParseFloat(d, 64)
		if err != nil {
			return true
		}
		_, err = strconv.ParseInt(d, 10, 64)
		if err != nil {
			return true
		}
		return false
	}
}

func DetectType(d []string) (EntityType, error) {
	unidentified := make([]string, 0)
	for typeInfo, detect := range detector {
		ok := true
		for _, s := range d {
			thisOk := detect(s)
			if !thisOk {
				unidentified = append(unidentified, s)
				ok = false
				break
			}
		}
		if ok {
			return typeInfo, nil
		}
	}
	return none, errors.New(fmt.Sprintf("Failed to identify - %v", unidentified))
}
