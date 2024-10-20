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

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	head *ListItem
	tail *ListItem
}

func (l *list)Len() int {
	return l.length
}

func (l *list)Front() *ListItem {
	return l.head
}

func (l *list)Back() *ListItem {
	return l.tail
}

func (l *list)PushFront(v interface{}) *ListItem {
	li := &ListItem{Value: v}
	l.length++

	if l.head == nil {
		l.head = li
		l.tail = li
	} else {
		li.Next = l.head
		l.head.Prev = li
		l.head = li
	}

	return li
}

func (l *list)PushBack(v interface{}) *ListItem {
	li := &ListItem{Value: v}
	l.length++

	if l.head == nil {
		l.head = li
		l.tail = li
	} else {
		li.Prev = l.tail
		l.tail.Next = li
		l.tail = li
	}

	return li
}

func (l *list)Remove(li *ListItem) {
	if l.length == 0 {
		return
	}

	prev, next := li.Prev, li.Next
	l.length--

	if prev != nil {
		prev.Next = next
	} else {
		l.head = next
	}

	if next != nil {
		next.Prev = prev
	} else {
		l.tail = prev
	}
}

func (l *list)MoveToFront(li *ListItem) {
	if l.length == 0 {
		return
	}

	val := li.Value
	l.Remove(li)
	l.PushFront(val)
}

func NewList() List {
	return new(list)
}
