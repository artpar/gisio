package reader

import (
	"encoding/csv"
	"os"
)

type csvReader struct {
	filename string
	reader   *csv.Reader
}

func NewCsvReader(filename string) (csvReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return csvReader{}, err
	}
	return csvReader{filename: filename, reader: csv.NewReader(file)}
}

func (c csvReader) ReadRow() ([]string, error) {
	return c.reader.Read()
}

func (c csvReader) ReadTable() ([][]string, error) {
	return c.reader.ReadAll()
}