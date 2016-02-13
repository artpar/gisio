package table

import (
	"github.com/artpar/gisio/types"
	"strconv"
	"log"
)

type LoadedFile struct {
	data [][]string
	*FileInfo
}

func (l LoadedFile) GetData(i, j int) string {
	return l.data[i][j]
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

func (file LoadedFile) ProcessLoadedFile() {
	file.DetectColumnTypes()
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


