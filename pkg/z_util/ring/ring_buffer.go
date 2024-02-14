package ring

import "container/list"

type Buffer[T any] struct {
	list *list.List
}

func NewBuffer[T any]() Buffer[T] {
	return Buffer[T]{
		list: list.New(),
	}
}

func (b *Buffer[T]) Get(pos int) T {
	if pos < 0 || pos >= b.list.Len() {
		var zero T
		return zero
	}

	element := b.list.Front()
	for i := 0; i < pos; i++ {
		element = element.Next()
	}
	return element.Value.(T)
}

func (b *Buffer[T]) GetFirst() (T, bool) {
	if b.list.Len() == 0 {
		var zero T
		return zero, false
	}
	return b.list.Front().Value.(T), true
}

func (b *Buffer[T]) GetLast() (T, bool) {
	if b.list.Len() == 0 {
		var zero T
		return zero, false
	}
	return b.list.Back().Value.(T), true
}

func (b *Buffer[T]) AddFirst(element T) {
	b.list.PushFront(element)
}

func (b *Buffer[T]) AddLast(element T) {
	b.list.PushBack(element)
}

func (b *Buffer[T]) RemoveFirst() {
	if b.list.Len() > 0 {
		b.list.Remove(b.list.Front())
	}
}

func (b *Buffer[T]) RemoveLast() {
	if b.list.Len() > 0 {
		b.list.Remove(b.list.Back())
	}
}

func (b *Buffer[T]) Len() int {
	return b.list.Len()
}
