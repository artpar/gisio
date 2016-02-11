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
	data         [][]string
	typeInfo     map[string]types.EntityType
	columnsCount int
	rowCount     int
}

func DetectColumnTypes(file LoadedFile) {
	file.typeInfo = make(map[string]types.EntityType)
	file.rowCount = len(file.data)
	if file.rowCount == 0 {
		return
	}
	file.columnsCount = len(file.data[0])
	for i := 0; i < file.columnsCount; i++ {
		colValues := make([]string, 10)
		for j := 0; j < file.rowCount && j < 10; j++ {
			colValues = append(colValues, file.data[j][i])
		}
		typeInfo, err := types.DetectType(colValues)
		if err != nil {
			log.Printf("Could not deduce type - %v", colValues)
		}
		file.typeInfo[i] = typeInfo
	}
}

func ProcessLoadedFile(file LoadedFile) {
	DetectColumnTypes(file)
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
	dataMap[filename] = NewLoadedFile(data)
	return nil
}

func NewLoadedFile(data [][]string) LoadedFile {
	loadedFile := LoadedFile{data: data}
	go ProcessLoadedFile(loadedFile)
	return loadedFile
}

func info(w http.ResponseWriter, r *http.Request)

func index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename, _ := vars["filename"]
	LoadData(filename)
	data := dataMap[filename].data
	numberOfRows := len(data)
	numberOfColumns := len(data[0])
	fmt.Printf("Rows: %d, Cols: %d", numberOfRows, numberOfColumns)
	// Print2D(data)

	con := make(map[string]interface{})
	con["Filename"] = filename
	con["Rows"] = numberOfRows
	con["Cols"] = numberOfColumns
	Render("data", w, con)
}

func Print2D(data [][]string) {
	for _, x := range data {
		for _, y := range x {
			fmt.Printf("%s\t", y)
		}
		fmt.Println()
	}
}
