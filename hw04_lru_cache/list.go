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
	frontItem *ListItem
	backItem  *ListItem
	len       int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.Front(),
	}
	if l.len == 0 {
		l.backItem = item
	} else {
		l.frontItem.Prev = item
	}
	l.frontItem = item
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  l.Back(),
		Next:  nil,
	}
	if l.len == 0 {
		l.frontItem = item
	} else {
		l.backItem.Next = item
	}
	l.backItem = item
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.frontItem = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.backItem = i.Prev
	}
	i.Prev = nil
	i.Next = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.frontItem {
		return
	}
	if i == l.backItem {
		l.backItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	i.Prev.Next = i.Next
	l.Front().Prev = i
	i.Next = l.Front()
	i.Prev = nil
	l.frontItem = i
}

func NewList() List {
	return new(list)
}
