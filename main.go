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
	"html/template"
	"io"
	"github.com/howeyc/fsnotify"
	"errors"
	"github.com/artpar/gisio/types"
	"encoding/json"
	"strconv"
)

const (
	resourceDir = "resources"
	htmlTemplatesDir = resourceDir + "/html"
)

var templates = template.Must(template.ParseGlob(htmlTemplatesDir + "/*.html"))

func init() {
	dataMap = make(map[string]LoadedFile)
	watcher, err := fsnotify.NewWatcher()
	CheckErr(err, "Failed to create new watcher")

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
	dirName = os.Args[1]

	rtr := mux.NewRouter()
	rtr.HandleFunc("/data/{filename:.+}/index.html", index).Methods("GET")
	rtr.HandleFunc("/data/{filename:.+}/info", info).Methods("GET")
	rtr.HandleFunc("/data/list", list).Methods("GET")
	rtr.HandleFunc("/meta/templates", templateList).Methods("GET")
	rtr.HandleFunc("/", list).Methods("GET")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./resources/static"))))
	http.Handle("/", rtr)
	log.Println("Listening...")
	http.ListenAndServe(":2299", nil)
}

func templateList(w http.ResponseWriter, r *http.Request) {
	Render("templateList", w, templates.Templates())
}

func CheckErr(err error, msg string) {
	if err != nil {
		const size = 4096
		stack := make([]byte, size)
		stack = stack[:runtime.Stack(stack, false)]

		log.Panic(msg + " \n" + err.Error() + "\n" + string(stack))
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	list, err := ioutil.ReadDir(dirName)
	CheckErr(err, "Failed to read data dir")
	Render("dataList", w, list)
}

type LoadedFile struct {
	data [][]string
	*FileInfo
}

type FileInfo struct {
	Filename    string  `json:"Filename"`
	ColumnCount int  `json:"ColumnCount"`
	RowCount    int  `json:"RowCount"`
	HasHeaders  bool
	ColumnInfo  []ColumnInfo
}

type ColumnInfo struct {
	TypeInfo           types.EntityType
	IsEnum             bool
	DistinctValueCount int
	ValueCounts        map[string]int
	Percent            int
	ColumnName         string
}

func (file LoadedFile) DetectColumnTypes() {
	file.RowCount = len(file.data)
	if file.RowCount == 0 {
		log.Printf("Row count is zero")
		return
	}
	log.Printf("Number of rows : %d\n", file.RowCount)
	file.ColumnCount = len(file.data[0])
	file.FileInfo.ColumnInfo = make([]ColumnInfo, file.ColumnCount)
	enumThreshHold := (file.RowCount * 15) / 100

	hasHeaders := false
	for i := 0; i < file.ColumnCount; i++ {
		thisColumnHeaders := false
		colValues := make([]string, 0)
		for j := 0; j < file.RowCount && j < 10; j++ {
			colValues = append(colValues, file.data[j][i])
		}
		log.Printf("Values for detection 1 - %s", colValues)
		var err error
		temp1, thisColumnHeaders, err := types.DetectType(colValues)
		if err != nil {
			log.Printf("Could not deduce type 1 - %v - %v", colValues, err)
		}
		if thisColumnHeaders {
			hasHeaders = true
		}

		distinctCount := 0
		counted := make(map[string]int, 0)
		isEnum := true
		startAt := 0
		if thisColumnHeaders {
			startAt = 1
		}
		for j := startAt; j < file.RowCount; j++ {
			_, ok := counted[file.data[j][i]]
			if ok {
				counted[file.data[j][i]] = counted[file.data[j][i]] + 1
			} else {
				distinctCount = distinctCount + 1
				counted[file.data[j][i]] = 1
			}

			if distinctCount > enumThreshHold && isEnum {
				isEnum = false
			}

		}

		if !isEnum {
			counted = make(map[string]int, 0)
		}
		columnName := "column_" + strconv.Itoa(i)
		if thisColumnHeaders {
			columnName = file.data[0][i]
		}

		file.FileInfo.ColumnInfo[i] = ColumnInfo{
			TypeInfo:temp1,
			IsEnum: isEnum,
			DistinctValueCount: distinctCount,
			ValueCounts: counted,
			Percent: (distinctCount * 100) / file.RowCount,
			ColumnName: columnName,
		}
	}

	file.FileInfo.HasHeaders = hasHeaders
	log.Printf("FileInfo: %v", file.FileInfo)
}

func (file LoadedFile) ProcessLoadedFile() {
	file.DetectColumnTypes()
}

var dataMap map[string]LoadedFile

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
	dataMap[filename] = NewLoadedFile(filename, data)
	return nil
}

func NewLoadedFile(filename string, data [][]string) LoadedFile {
	t := make([]ColumnInfo, 0)
	loadedFile := LoadedFile{data: data,
		FileInfo: &FileInfo{
			Filename: filename,
			ColumnInfo: t,
		},
	}
	log.Printf("Process loaded file - %s", filename)
	loadedFile.ProcessLoadedFile()
	return loadedFile
}

func SendJson(w http.ResponseWriter, d interface{}) {
	by, err := json.Marshal(d)
	CheckErr(err, fmt.Sprintf("Failed to write value as json - %v", err))
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
