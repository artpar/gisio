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
	"encoding/json"
	"github.com/artpar/gisio/table"
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
	rtr.HandleFunc("/data/{filename:.+}/operation", operation).Methods("GET")
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

func operation(w http.ResponseWriter, r *http.Request) {
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
