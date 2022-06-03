package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.head
}

func (l *list) Back() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	var elem *ListItem
	if l.len == 0 {
		elem = &ListItem{Next: nil, Prev: nil, Value: v}
		l.head = elem
		l.tail = elem
	} else {
		elem = &ListItem{Next: l.head, Prev: nil, Value: v}
		l.head.Prev = elem
		l.head = elem
	}
	l.len++
	return elem
}

func (l *list) PushBack(v interface{}) *ListItem {
	var elem *ListItem
	if l.len == 0 {
		elem = &ListItem{Next: nil, Prev: nil, Value: v}
		l.head = elem
		l.tail = elem
	} else {
		elem = &ListItem{Next: nil, Prev: l.tail, Value: v}
		l.tail.Next = elem
		l.tail = elem
	}
	l.len++
	return elem
}

//крайние случаи
func (l *list) Remove(i *ListItem) {
	if i == l.head {
		l.head = l.head.Next
		l.len--
		return
	}
	elemPrev := i.Prev
	if i == l.tail {
		elemPrev.Next = nil
		l.tail = elemPrev
		l.len--
		return
	}
	elemNext := i.Next
	elemPrev.Next = elemNext
	elemNext.Prev = elemPrev
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head {
		return
	}
	elemPrev := i.Prev
	var elemNext *ListItem
	i.Prev = nil
	i.Next = l.head
	l.head = i
	if i == l.tail {
		elemPrev.Next = nil
		l.tail = elemPrev
	} else {
		elemNext = i.Next
		elemPrev.Next = elemNext
		elemNext.Prev = elemPrev
	}
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func NewList() List {
	return new(list)
}
