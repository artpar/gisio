package main

import (
  "fmt"
  "os"
  "github.com/artpar/gisio/reader"
  "net/http"
  "github.com/gorilla/mux"
  "log"
  "io/ioutil"
  "runtime"
  "text/template"
  _ "net/http/pprof"

  "io"
  "github.com/howeyc/fsnotify"
  "errors"
  "encoding/json"
  "github.com/artpar/gisio/table"
  "strings"
  "github.com/artpar/gisio/types"
  "strconv"
  "github.com/artpar/gisio/grossfilter"
)

const (
  resourceDir = "resources"
  htmlTemplatesDir = resourceDir + "/html"
)

var templates = template.Must(template.ParseGlob(htmlTemplatesDir + "/*.html"))

func init() {
  dataMap = make(map[string]grossfilter.GrossFilter)
  watcher, err := fsnotify.NewWatcher()
  CheckErr(err, "Failed to create new watcher")

  go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
  }()
  go func() {
    for {
      select {
      case ev := <-watcher.Event:
        if ev.IsDelete() {
          templates = template.Must(template.ParseGlob(htmlTemplatesDir + "/*"))
        }
      case err := <-watcher.Error:
        log.Println("error:", err)
      }
    }
  }()
  err = watcher.Watch(htmlTemplatesDir)
  CheckErr(err, "Start watch failed")
}

func Render(templateName string, w io.Writer, d interface{}) {
  templates.ExecuteTemplate(w, templateName, d)
}

func Send(msg... interface{}) {
  fmt.Printf(msg[0].(string), msg[1:]...)
}

var dirName string

func main() {

  if len(os.Args) < 2 {
    log.Println("Usage: <exec> <dirName> # dirName is path to csv files directory")
    return
  }

  dirName = os.Args[1]

  rtr := mux.NewRouter()
  rtr.HandleFunc("/data/{filename:.+}/index.html", index).Methods("GET")
  rtr.HandleFunc("/data/{filename:.+}/info", info).Methods("GET")
  rtr.HandleFunc("/data/{filename:.+}/operation", operation).Methods("GET")
  rtr.HandleFunc("/data/list", list).Methods("GET")
  rtr.HandleFunc("/meta/templates", templateList).Methods("GET")
  rtr.HandleFunc("/", list).Methods("GET")

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./resources/static"))))
  http.Handle("/", rtr)
  log.Println("Listening... on http://localhost:2299/")
  http.ListenAndServe(":2299", nil)
}

func templateList(w http.ResponseWriter, r *http.Request) {
  Render("templateList", w, templates.Templates())
}

func CheckErr(err error, msg... interface{}) {
  if err != nil {
    const size = 4096
    stack := make([]byte, size)
    stack = stack[:runtime.Stack(stack, false)]

    log.Panic(fmt.Sprintf(msg[0].(string), msg[1:]...) + " \n" + err.Error() + "\n" + string(stack))
  }
}

func list(w http.ResponseWriter, r *http.Request) {
  list, err := ioutil.ReadDir(dirName)
  CheckErr(err, "Failed to read data dir")
  Render("dataList", w, list)
}

var dataMap map[string]grossfilter.GrossFilter

func LoadData(filename string) (error) {
  location := dirName + "/" + filename
  _, ok := dataMap[filename]
  if ok {
    log.Printf("Filename %s is already loaded\n", filename)
    return nil
  }

  reader, err := reader.NewCsvReader(location)
  if err != nil {
    return err
  }
  data, err := reader.ReadTable()
  if err != nil {
    return err
  }
  if (len(data) < 0) {
    return errors.New(fmt.Sprintf("Data is too less in [%s]. Need atleast 4 rows", dirName))
  }
  dataMap[filename] = grossfilter.NewGrossFilter(table.NewLoadedFile(filename, data))
  return nil
}

func SendJson(w http.ResponseWriter, d interface{}) {
  by, err := json.Marshal(d)
  CheckErr(err, fmt.Sprintf("Failed to write value as json - %v", err))
  w.Header().Set("Content-type", "application/json")
  w.Write(by)
}

func info(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  filename, _ := vars["filename"]
  log.Printf("Get info for file - %s", filename)
  f, ok := dataMap[filename]
  if !ok {
    log.Printf("Load file info - %s\n", filename)
    err := LoadData(filename)
    CheckErr(err, "Load failed")
  }
  f, _ = dataMap[filename]
  b, _ := json.Marshal(f)
  log.Printf("Data - %s", string(b))
  SendJson(w, &f)
}

type Query struct {
  Operation  string
  Function   string
  ColumnName string
  Data       []Query
}

func operation(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  vars := mux.Vars(r)
  filename, _ := vars["filename"]
  queryString := r.Form["q"][0]
  var query Query
  var result [][]string
  err := json.Unmarshal([]byte(queryString), &query)
  CheckErr(err, "Failed to understand input - " + queryString)
  log.Printf("Get info for file - %s", filename)
  f, ok := dataMap[filename]
  if !ok {
    log.Printf("Load file info - %s\n", filename)
    err := LoadData(filename)
    CheckErr(err, "Load failed")
  }
  f, _ = dataMap[filename]
  log.Printf("Query - %v", query)
  switch strings.ToLower(query.Operation) {
  case "groupby":
    log.Printf("Group by ")
    switch strings.ToLower(query.Function) {
    case "sum":
      log.Printf("Doing sum of %v over %v", query.Data[0].ColumnName, query.Data[1].ColumnName)
      resultMap := make(map[string]interface{})
      col1Index := GetColumnIndexByName(f.ColumnInfo, query.Data[0].ColumnName)
      col2Index := GetColumnIndexByName(f.ColumnInfo, query.Data[1].ColumnName)

      if col2Index == -1 || col1Index == -1 {
        log.Panicf("Col 1 or Col 2 not found %d,%d", col1Index, col2Index)
      }
      if f.ColumnInfo[col2Index].TypeInfo != types.Number {
        log.Panicf("Col 2 is not a numeric column - " + f.ColumnInfo[col2Index].TypeInfo.String())
      }
      for i := f.FirstRowIndex; i < f.RowCount; i++ {
        key := f.GetData(i, col1Index)
        value := f.GetData(i, col2Index)
        if value == "NA" {
          value = "0"
        }
        numValue, err := strconv.ParseFloat(value, 64)
        CheckErr(err, "Failed to parse %v as number", numValue)
        _, ok := resultMap[key]
        if ok {
          resultMap[key] = resultMap[key].(float64) + numValue
        } else {
          resultMap[key] = numValue
        }
      }
      result = Map2Array(resultMap)

    case "count":
      log.Printf("Doing count of %v over %v", query.Data[0].ColumnName, query.Data[1].ColumnName)
      //resultMap := make(map[string]interface{})
      countMap := make(map[string]map[string]int)

      col1Index := GetColumnIndexByName(f.ColumnInfo, query.Data[0].ColumnName)
      col2Index := GetColumnIndexByName(f.ColumnInfo, query.Data[1].ColumnName)

      if col2Index == -1 || col1Index == -1 {
        log.Panicf("Col 1 or Col 2 not found %d,%d", col1Index, col2Index)
      }
      for i := f.FirstRowIndex; i < f.RowCount; i++ {
        key := f.GetData(i, col1Index)
        _, ok := countMap[key]
        if !ok {
          m := make(map[string]int)
          countMap[key] = m
        }
        value := f.GetData(i, col2Index)
        if value == "NA" {
          value = "0"
        }
        //log.Printf("Increase count for [%v] [%v]", key, value)
        m1 := countMap[key]
        _, ok = m1[value]
        if ok {
          countMap[key][value] = countMap[key][value] + 1
        } else {
          countMap[key][value] = 1
        }
      }

      resultMap := map[string]interface{}{}

      for key, val := range countMap {
        //_, ok := resultMap[key]

        resultMap[key] = 0

        for range val {
          resultMap[key] = resultMap[key].(int) + 1
        }
      }

      result = Map2Array(resultMap)


    }
  }
  SendJson(w, &result)
}

func GetColumnIndexByName(columns []table.ColumnInfo, columnName string) int {
  for i, col := range columns {
    if col.ColumnName == columnName {
      return i
    }
  }
  return -1
}

func Map2Array(m map[string]interface{}) [][]string {
  result := make([][]string, 0)
  for k, v := range m {
    if len(k) < 1 {
      continue
    }
    result = append(result, []string{k, fmt.Sprintf("%v", v)})
  }

  return result
}

func index(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  filename, _ := vars["filename"]
  err := LoadData(filename)
  CheckErr(err, "Load failed")
  fmt.Printf("Rows: %d, Cols: %d", dataMap[filename].RowCount, dataMap[filename].ColumnCount)
  // Print2D(data)

  Render("data", w, dataMap[filename].FileInfo)
}

func Print2D(data [][]string) {
  for _, x := range data {
    for _, y := range x {
      fmt.Printf("%s\t", y)
    }
    fmt.Println()
  }
}
