package reader

type TableReader interface {
	ReadRow() ([]string, error)
	ReadTable() ([][]string, error)
}