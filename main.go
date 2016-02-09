package main

import (
	"fmt"
	"os"
	"github.com/artpar/gisio/reader"
)

func Send(msg... string) {
	fmt.Printf(msg)
}

func main() {
	filename := os.Args[0]
	reader, err := reader.NewCsvReader(filename)
	if err != nil {
		panic(err)
	}
	data, err := reader.ReadTable()
	if err != nil {
		panic(err)
	}
	if (len(data) < 4) {
		Send("Data is too less in [" + filename + "]. Need atleast 4 rows")
		return
	}
	numberOfRows := len(data)
	numberOfColumns := len(data[0])
}
