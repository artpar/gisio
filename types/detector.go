package types

import (
	"github.com/artpar/gisio/mtime"
	"net"
	"strconv"
	"regexp"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type EntityType int

func (t EntityType) String() string {
	switch t {
	case Time:
		return "time"
	case Ipaddress:
		return "ipaddress"
	case Money:
		return "money"
	case Number:
		return "number"
	case None:
		return "none"
	}
	return "failed-to-detect"
}

func (t EntityType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.String() + "\""), nil
}

const (
	Time EntityType = iota
	Ipaddress
	Money
	Number
	None
)

var (
	order = []EntityType{Time, Ipaddress, Number, Money}
)
var detector map[EntityType]func(string) (bool, interface{})

func init() {
	detector = make(map[EntityType]func(string) (bool, interface{}))
	detector[Time] = func(d string) (bool, interface{}) {
		t, _, err := mtime.GetTime(d)
		if err == nil {
			return true, t
		}
		return false, time.Now()
	}
	detector[Ipaddress] = func(d string) (bool, interface{}) {
		s := net.ParseIP(d)
		if s != nil {
			return true, net.IP("")
		}
		return false, s
	}
	detector[Money] = func(d string) (bool, interface{}) {
		r := regexp.MustCompile("^([a-zA-Z]{0,3}\\.?)?[0-9]+\\.[0-9]{0,2}([a-zA-Z]{0,3})?")
		return r.MatchString(d), d
	}
	detector[Number] = func(d string) (bool, interface{}) {
		v, err := strconv.ParseFloat(d, 64)
		if err == nil {
			return true, v
		}
		log.Printf("Parse %v as float failed - %v", d, err)
		v1, err := strconv.ParseInt(d, 10, 64)
		if err == nil {
			return true, v1
		}
		log.Printf("Parse %v as int failed - %v", d, err)
		return false, 0
	}
}

func ConvertValues(d []string, typ EntityType) ([]interface{}, error) {
	converted := make([]interface{}, len(d))
	converter, ok := detector[typ]
	if !ok {
		log.Printf("Converter not found for %v", typ)
		return converted, errors.New("Converter not found for " + typ.String())
	}
	for i, v := range d {
		ok, val := converter(v)
		if !ok {
			log.Printf("Conversion of %s as %v failed", v, typ)
			continue
		}
		converted[i] = val
	}
	return converted, nil
}

func DetectType(d []string) (EntityType, bool, error) {
	unidentified := make([]string, 0)
	thisHeaders := false
	for _, typeInfo := range order {
		detect := detector[typeInfo]
		log.Printf("Try 1 %s as %v", d, typeInfo)
		ok := true
		for _, s := range d {
			thisOk, _ := detect(s)
			if !thisOk {
				unidentified = append(unidentified, s)
				ok = false
				break
			}
		}
		if ok {
			return typeInfo, thisHeaders, nil
		}
	}

	foundType := None
	thisHeaders = true
	for _, typeInfo := range order {
		detect := detector[typeInfo]
		log.Printf("Try 2 %s as %v", d[1:], typeInfo)
		ok := true
		for _, s := range d[1:] {
			thisOk, _ := detect(s)
			if !thisOk {
				unidentified = append(unidentified, s)
				ok = false
				break
			}
		}
		if ok {
			foundType = typeInfo
			break
		}
	}

	if thisHeaders {
		columnName := d[0]
		typeByColumnName := columnTypeFromName(columnName)
		if typeByColumnName != None {
			foundType = typeByColumnName
		}
	}

	if foundType != None {
		return foundType, thisHeaders, nil
	}

	return None, thisHeaders, errors.New(fmt.Sprintf("Failed to identify - %v", unidentified))
}

var nameMap = map[string]EntityType{
	"price": Money,
	"income": Money,
	"amount": Money,
	"wage": Money,
	"cost":Money,
}

func columnTypeFromName(name string) EntityType {
	for n, typ := range nameMap {
		if strings.HasSuffix(name, n) {
			log.Printf("Selecting type %s because of Suffix %s in %s", typ.String(), n, name)
			return typ
		}
		if strings.HasPrefix(name, n) {
			log.Printf("Selecting type %s because of Prefix %s in %s", typ.String(), n, name)
			return typ
		}
	}
	return None
}
