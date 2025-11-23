package types

type DiffItem struct {
	Key      string
	Value    []interface{}
	Result   string
	Children []DiffItem
	Path     string
}
