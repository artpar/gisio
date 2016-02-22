package grossfilter

import (
	"github.com/artpar/gisio/table"
	"github.com/artpar/gisio/types"
	"log"
)

type GrossFilterInterface interface {
	// add a new record
	Add([]string)
	// remove all records which match the current filter
	RemoveData()
	// add a new dimension with a specified accessor function
	Dimension()
	// convienence method for group all on a dummy dimention
	GroupAll()
	// returns number of records in this grossfilter, irrespective of filters
	Size()
}

type GrossFilter struct {
	columnData [][]interface{} `json:"-"`
	table.LoadedFile
	n          int
	// a bit mask representing which dimension are in use
	m          int
	// number of dimensions which can fit in a dimension
	M          int
	filters    []int
}

func (g GrossFilter) Add(newData [][]string) GrossFilter {
	n1 := len(newData)
	if n1 > 0 {
		g.LoadedFile.AddRows(newData)
		g.n = g.n + n1
		g.filters = grossfilter_arrayLengthen(g.filters, g.n)
	}
	return g
}

func grossfilter_index(n, m int) []interface{} {
	return grossfilter_array8(n)
}

func (g GrossFilter) RemoveData() {
	//newIndex := grossfilter_index(g.n, g.n)
	//var removed [][]int
	//j := 0
	//for i := 0; i < g.n; i ++ {
	//	if (g.filters[i] != 0) {
	//		newIndex[i] = j
	//		j += 1
	//	} else {
	//		removed = append(removed, i)
	//	}
	//}
	//
	//j = 0
	//for i := 0; i < g.n; i++ {
	//	k := g.filters[i]
	//	if k {
	//		if i != j {
	//			g.filters[j] = k
	//			g.Data[j] = g.Data[i]
	//		}
	//		j += 1
	//	}
	//}
	//g.Data = g.Data[0:j]
	//for ; g.n > j; {
	//	g.filters[g.n] = 0
	//	g.n = g.n - 1
	//}
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
	g := GrossFilter{columnData:columnData}
	g.LoadedFile = loadedfile
	return g
}