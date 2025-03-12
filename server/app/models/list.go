package models

import "iter"

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

func (l List[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range l.Items {
			if !yield(item) {
				return
			}
		}
	}
}

func (l List[T]) Seq2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, item := range l.Items {
			if !yield(i, item) {
				return
			}
		}
	}
}
