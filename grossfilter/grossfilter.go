package grossfilter

import (
	"github.com/artpar/gisio/table"
	"github.com/artpar/gisio/types"
	"log"
)

type GrossFilter struct {
	columnData [][]interface{}
	table.LoadedFile
}

func NewGrossFilter(loadedfile table.LoadedFile) GrossFilter {
	colCount := loadedfile.ColumnCount
	rowCount := loadedfile.RowCount

	columnData := make([][]interface{}, colCount)
	start := 0
	if loadedfile.HasHeaders {
		start = 1
	}

	for i := 0; i < colCount; i++ {
		initialValues := make([]string, rowCount - start)
		for j := start; j < rowCount; j++ {
			initialValues[j - start] = loadedfile.GetData(j, i)
		}
		colData, err := types.ConvertValues(initialValues, loadedfile.ColumnInfo[i].TypeInfo)
		if err != nil {
			log.Printf("Converion of types failed - %s", err)
		}
		columnData[i] = colData
	}
	return GrossFilter{columnData, loadedfile}
}