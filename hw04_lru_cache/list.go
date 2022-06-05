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
	elem := &ListItem{Prev: nil, Value: v}
	if l.len == 0 {
		elem.Next = nil
		l.head = elem
		l.tail = elem
	} else {
		elem.Next = l.head
		l.head.Prev = elem
		l.head = elem
	}
	l.len++
	return elem
}

func (l *list) PushBack(v interface{}) *ListItem {
	elem := &ListItem{Next: nil, Value: v}
	if l.len == 0 {
		elem.Prev = nil
		l.head = elem
		l.tail = elem
	} else {
		elem.Prev = l.tail
		l.tail.Next = elem
		l.tail = elem
	}
	l.len++
	return elem
}

func (l *list) Remove(i *ListItem) {
	elemPrev := i.Prev
	elemNext := i.Next
	if elemPrev == nil {
		l.head = l.head.Next
	} else if elemNext == nil {
		elemPrev.Next = nil
		l.tail = elemPrev
	} else {
		elemPrev.Next = elemNext
		elemNext.Prev = elemPrev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	elemPrev := i.Prev
	elemNext := i.Next
	if elemPrev == nil {
		return
	}
	i.Prev = nil
	i.Next = l.head
	l.head = i
	if elemNext == nil {
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
