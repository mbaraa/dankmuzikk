package entities

type List[T any] struct {
	Items  []T    `json:"items"`
	Size   int    `json:"size"`
	Cursor string `json:"cursor"`
}

func NewList[T any](items []T, cursor string) List[T] {
	return List[T]{
		Items:  items,
		Size:   len(items),
		Cursor: cursor,
	}
}
