package table

import (
  "github.com/artpar/gisio/types"
  "strconv"
  "log"
  "math"
  "math/rand"
  "strings"
)

type LoadedFile struct {
  Data          [][]string `json:"-"`
  convertedData [][]interface{}
  *FileInfo
}

func (l LoadedFile) GetData(i, j int) string {
  return l.Data[i][j]
}

func (l LoadedFile) AddRows(newData [][]string) {
  l.Data = append(l.Data, newData...)
}

type FileInfo struct {
  Filename      string  `json:"Filename"`
  ColumnCount   int  `json:"ColumnCount"`
  RowCount      int  `json:"RowCount"`
  HasHeaders    bool
  FirstRowIndex int
  ColumnInfo    []ColumnInfo
}

type ColumnStats struct {
  DistinctValueCount int
  ValueCounts        map[string]int
  Percent            int
}

type ColumnInfo struct {
  TypeInfo   types.EntityType
  IsEnum     bool
  IsUnique   bool
  ColumnName string
  ColumnStats
}

func (file LoadedFile) ProcessLoadedFile() error {
  file.DetectColumnTypes()
  return file.LoadData()
}

func NewLoadedFile(filename string, data [][]string) LoadedFile {
  t := make([]ColumnInfo, 0)
  loadedFile := LoadedFile{Data: data,
    FileInfo: &FileInfo{
      Filename: filename,
      ColumnInfo: t,
    },
  }
  log.Printf("Process loaded file - %s", filename)
  loadedFile.ProcessLoadedFile()
  return loadedFile
}

func (file LoadedFile) LoadData() error {
  start := 0
  if file.HasHeaders {
    start = 1
  }
  file.convertedData = make([][]interface{}, file.RowCount)
  for i := 0; i < file.RowCount; i ++ {
    file.convertedData[i] = make([]interface{}, file.ColumnCount)
  }
  for i := 0; i < file.ColumnCount; i++ {
    column := ColumnFrom2dArray(file.Data, i, start)
    convertedValues, err := types.ConvertValues(column, file.ColumnInfo[i].TypeInfo)
    if err != nil {
      return err
    }
    for j := start; j < file.RowCount; j++ {
      file.convertedData[j][i] = convertedValues[j - start]
    }
  }
  return nil
}

func ColumnFrom2dArray(data [][]string, colIndex int, start int) []string {
  res := make([]string, len(data) - start)
  for i := start; i < len(data); i++ {
    res[i - start] = strings.TrimSpace(data[i][colIndex])
  }
  return res
}

func (file LoadedFile) DetectColumnTypes() {
  file.RowCount = len(file.Data)
  if file.RowCount == 0 {
    log.Printf("Row count is zero")
    return
  }
  log.Printf("Number of rows : %d\n", file.RowCount)
  file.ColumnCount = len(file.Data[0])
  file.FileInfo.ColumnInfo = make([]ColumnInfo, file.ColumnCount)
  enumThreshHold := int(math.Min(float64((file.RowCount * 15) / 100), 70))

  hasHeaders := false
  for i := 0; i < file.ColumnCount; i++ {
    thisColumnHeaders := false
    colValues := make([]string, 0)
    colValues = append(colValues, file.Data[0][i])
    for j := 1; j < file.RowCount && j < 20; j++ {
      colValues = append(colValues, file.Data[rand.Intn(file.RowCount)][i])
    }
    log.Printf("Values for detection 1 - %s", colValues)
    var err error
    temp1, thisColumnHeaders, err := types.DetectType(colValues)
    log.Printf("Type deduction for [%v]: %v", colValues, temp1)
    if err != nil {
      log.Printf("Could not deduce type 1 - %v - %v", colValues, err)
    }
    if thisColumnHeaders {
      hasHeaders = true
    }

    distinctCount := 0
    counted := make(map[string]int, 0)
    isEnum := true
    isUnique := true
    startAt := 0
    if thisColumnHeaders {
      startAt = 1
    }
    for j := startAt; j < file.RowCount; j++ {
      _, ok := counted[file.Data[j][i]]
      if ok {
        counted[file.Data[j][i]] = counted[file.Data[j][i]] + 1
      } else {
        isUnique = false
        distinctCount = distinctCount + 1
        counted[file.Data[j][i]] = 1
      }

      if distinctCount > enumThreshHold && isEnum {
        isEnum = false
      }

    }
    if temp1 == types.Timestamp || temp1 == types.Rating5 || temp1 == types.Rating10 || temp1 == types.Rating100 {
      isEnum = false
    }
    if !isEnum {
      counted = make(map[string]int, 0)
    }
    columnName := "column_" + strconv.Itoa(i)
    if thisColumnHeaders {
      columnName = file.Data[0][i]
    }

    file.FileInfo.ColumnInfo[i] = ColumnInfo{
      TypeInfo:temp1,
      IsEnum: isEnum,
      IsUnique: isUnique,
      ColumnName: columnName,
      ColumnStats: ColumnStats{
        DistinctValueCount: distinctCount,
        ValueCounts: counted,
        Percent: (distinctCount * 100) / file.RowCount,
      },
    }
  }

  file.FileInfo.HasHeaders = hasHeaders
  file.FileInfo.FirstRowIndex = 0
  if hasHeaders {
    file.FileInfo.FirstRowIndex = 1
    file.RowCount = file.RowCount - 1
  }
  log.Printf("FileInfo: %v", file.FileInfo)
}


