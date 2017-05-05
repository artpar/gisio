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
  "sort"
)

type EntityType int

func (t EntityType) String() string {
  switch t {
  case Time:
    return "time"
  case Date:
    return "date"
  case DateTime:
    return "datetime"
  case Ipaddress:
    return "ipaddress"
  case Money:
    return "money"
  case Number:
    return "number"
  case None:
    return "none"
  case Boolean:
    return "boolean"
  case Latitude:
    return "location-latitude"
  case Longitude:
    return "location-longitude"
  case City:
    return "location-city"
  case Country:
    return "location-country"
  case Continent:
    return "location-continent"
  case State:
    return "location-state"
  case Pincode:
    return "location-pincode"
  case Timestamp:
    return "timestamp"
  case Rating5:
    return "rating5"
  case Rating10:
    return "rating10"
  case Rating100:
    return "rating100"
  case Id:
    return "id-col"
  }
  return "name-not-set"
}

func (t EntityType) MarshalJSON() ([]byte, error) {
  return []byte("\"" + t.String() + "\""), nil
}

const (
  DateTime EntityType = iota
  Id
  Time
  Date
  Ipaddress
  Money
  Rating5
  Rating10
  Rating100
  Timestamp
  Number
  Boolean
  Latitude
  Longitude
  City
  Country
  Continent
  State
  Pincode

  None
)

var (
  order = []EntityType{Id, Boolean, DateTime, Date, Time, Rating5, Rating10, Rating100, Timestamp, Ipaddress, Latitude, Latitude, City, Country, Continent, State, Pincode, Number, Money}
)
var detector map[EntityType]func(string) (bool, interface{})

var CurrencyType struct {
  Name  string
  Value float64
}

func init() {
  sort.Sort(&unknownNumbers)
  detector = make(map[EntityType]func(string) (bool, interface{}))
  detector[Time] = func(d string) (bool, interface{}) {
    //fmt.Printf("Try to parse [%v] with mtime\n", d)
    t, _, err := mtime.GetTime(d)
    //fmt.Errorf("Fail to parse [%v] with mtime: %v\n", d, err)
    if err == nil {
      return true, t
    }
    return false, time.Now()
  }

  //detector[Id] = func(d string)(bool, interface{}) {
  //
  //}

  detector[Timestamp] = func(d string) (bool, interface{}) {
    //fmt.Printf("Try to parse [%v] with mtime\n", d)


    i, err := strconv.ParseInt(d, 10, 64)
    if err != nil {
      return false, d
    }

    if i < 100000000 {
      return false, d
    }

    tm := time.Unix(i, 0)

    return true, tm
  }

  detector[Date] = func(d string) (bool, interface{}) {
    //fmt.Printf("Try to parse [%v] with mtime\n", d)
    t, _, err := mtime.GetDate(d)
    //fmt.Errorf("Fail to parse [%v] with mtime: %v\n", d, err)
    if err == nil {
      return true, t
    }
    return false, time.Now()
  }

  detector[DateTime] = func(d string) (bool, interface{}) {
    //fmt.Printf("Try to parse [%v] with mtime\n", d)
    t, _, err := mtime.GetDateTime(d)
    //fmt.Errorf("Fail to parse [%v] with mtime: %v\n", d, err)
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
    r := regexp.MustCompile("^([a-zA-Z]{0,3}\\.? )?[0-9]+\\.[0-9]{0,2}([a-zA-Z]{0,3})?")
    return r.MatchString(d), d
  }
  detector[Boolean] = func(d string) (bool, interface{}) {
    d = strings.ToLower(d)
    switch d {
    case "yes":
    case "1":
      d = "true"
    case "no":
    case "0":
      d = "false"
    }
    r, err := strconv.ParseBool(d)
    if err != nil {
      return false, false
    }
    return true, r
  }

  detector[Rating5] = func(d string) (bool, interface{}) {

    numberOk, nValue := detector[Number](d)

    if !numberOk {
      return false, d
    }

    nInt, ok := nValue.(int)
    if ok {
      if nInt <= 5 {
        return true, nInt
      } else {
        return false, nInt
      }
    }

    nFloat, ok := nValue.(float64)
    if ok {
      if nFloat <= 5.0 {
        return true, nFloat
      } else {
        return false, nFloat
      }
    }
    return false, nValue

  }

  detector[Rating10] = func(d string) (bool, interface{}) {

    numberOk, nValue := detector[Number](d)

    if !numberOk {
      return false, d
    }

    nInt, ok := nValue.(int)
    if ok {
      if nInt <= 10 {
        return true, nInt
      } else {
        return false, nInt
      }
    }

    nFloat, ok := nValue.(float64)
    if ok {
      if nFloat <= 10.0 {
        return true, nFloat
      } else {
        return false, nFloat
      }
    }
    return false, nValue

  }

  detector[Rating100] = func(d string) (bool, interface{}) {

    numberOk, nValue := detector[Number](d)

    if !numberOk {
      return false, d
    }

    nInt, ok := nValue.(int)
    if ok {
      if nInt <= 100 {
        return true, nInt
      } else {
        return false, nInt
      }
    }

    nFloat, ok := nValue.(float64)
    if ok {
      if nFloat <= 100.0 {
        return true, nFloat
      } else {
        return false, nFloat
      }
    }
    return false, nValue

  }

  detector[Number] = func(d string) (bool, interface{}) {
    d = strings.ToLower(d)
    in := sort.SearchStrings(unknownNumbers, d)
    if in < len(unknownNumbers) && unknownNumbers[in] == d {
      log.Printf("One of the unknowns - %v : %d", d, sort.SearchStrings(unknownNumbers, strings.ToLower(d)))
      return true, 0
    }
    v, err := strconv.ParseFloat(d, 64)
    if err == nil {
      return true, v
    }
    //log.Printf("Parse %v as float failed - %v", d, err)
    v1, err := strconv.ParseInt(d, 10, 64)
    if err == nil {
      return true, v1
    }
    //log.Printf("Parse %v as int failed - %v", d, err)
    return false, 0
  }

  detector[Latitude] = func(d string) (bool, interface{}) {

    var realFloatValue float64
    isFloat, floatValue := detector[Number](d)

    intVal, isInt := floatValue.(int)

    if !isInt {
      floatVal, isReallyFloat := floatValue.(float64)

      if !isReallyFloat {
        return false, floatValue
      }
      realFloatValue = floatVal
    } else {
      realFloatValue = float64(intVal)
    }

    if !isFloat || realFloatValue > 180.0 {
      return false, floatValue
    }

    return true, floatValue

  }

  detector[None] = func(d string) (bool, interface{}) {
    return true, d
  }
}

var (
  unknownNumbers = sort.StringSlice([]string{"na", "n/a", "-"})
)

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
      // log.Printf("Conversion of %s as %v failed", v, typ)
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

    if detect == nil {
      log.Printf("No detector for type [%v]\n", typeInfo)
      continue
    }

    ok := true
    for _, s := range d {
      log.Printf("Detect value: %v", s)
      t := strings.TrimSpace(s)
      log.Printf("Detected value: %v", t)
      thisOk, _ := detect(t)
      if !thisOk {
        unidentified = append(unidentified, s)
        ok = false
        break
      }
    }
    //log.Printf("Try 1 %s as %v - %v", d, typeInfo, ok)
    if ok {
      return typeInfo, thisHeaders, nil
    }
  }
  thisHeaders = true
  foundType := None

  columnHeader := d[0]
  typeByColumnName := columnTypeFromName(columnHeader)
  if typeByColumnName != None {
    foundType = typeByColumnName
  }

  if foundType == None {
    for _, typeInfo := range order {
      detect := detector[typeInfo]

      if detect == nil {
        log.Printf("No detector available for type [%v]", typeInfo)
        continue
      }

      ok := true
      for _, s := range d[1:] {
        thisOk, _ := detect(s)
        if !thisOk {
          unidentified = append(unidentified, s)
          ok = false
          break
        }
      }
      //log.Printf("Try 2 %s as %v - %v", d[1:], typeInfo, ok)
      if ok {
        foundType = typeInfo
        break
      }
    }
  }

  if foundType != None {
    return foundType, thisHeaders, nil
  }

  return None, thisHeaders, errors.New(fmt.Sprintf("Failed to identify - %v", unidentified))
}

var nameMap = map[EntityType][]string{
  Id: []string{"id"},
  Money: []string{"price", "income", "amount", "wage", "cost", "sale", "profit", "asset", "marketvalue"},
  Latitude: []string{"lat", "latitude"},
  Longitude: []string{"lon", "long", "longitude"},
  City: []string{"city"},
  Country: []string{"country"},
  State: []string{"state"},
  Continent: []string{"continent"},
  Pincode: []string{"pincode", "zipcode"},
}

func columnTypeFromName(name string) EntityType {
  name = strings.ToLower(name)
  for typ, names := range nameMap {
    for _, n := range names {
      if strings.HasSuffix(name, n) {
        log.Printf("Selecting type %s because of Suffix %s in %s", typ.String(), n, name)
        return typ
      }
      if strings.HasPrefix(name, n) {
        log.Printf("Selecting type %s because of Prefix %s in %s", typ.String(), n, name)
        return typ
      }

      if len(n) > 5 && strings.Index(name, n) > -1 {
        log.Printf("Selecting type %s because of Prefix %s in %s", typ.String(), n, name)
        return typ
      }
    }
  }
  return None
}
