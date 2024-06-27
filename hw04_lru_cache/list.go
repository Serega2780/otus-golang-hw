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
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.front == nil {
		l.front = &ListItem{Value: v, Next: nil, Prev: nil}
		l.back = l.front
	} else {
		temp := l.front
		l.front = &ListItem{Value: v, Next: temp, Prev: nil}
		temp.Prev = l.front
	}
	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.back == nil {
		return l.PushFront(v)
	}
	if l.back == l.front {
		l.back = &ListItem{Value: v, Next: nil, Prev: l.front}
		l.front.Next = l.back
	} else {
		temp := l.back
		l.back = &ListItem{Value: v, Next: nil, Prev: temp}
		temp.Next = l.back
	}
	l.len++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	if l.len == 1 {
		l.front, l.back = nil, nil
		l.len = 0
		return
	}
	switch i {
	case l.front:
		l.front = l.front.Next
		l.front.Prev = nil
	case l.back:
		l.back = l.back.Prev
		l.back.Next = nil
	default:
		changeLinks(i)
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	if i == l.back {
		newBack := i.Prev
		temp := l.front
		temp.Prev = i
		i.Next = temp
		l.front = i
		l.front.Prev = nil
		l.back = newBack
		l.back.Next = nil
	} else {
		changeLinks(i)
		temp := l.front
		temp.Prev = i
		i.Next = temp
		l.front = i
	}
}

func NewList() List {
	return &list{front: nil, back: nil, len: 0}
}

func changeLinks(val *ListItem) {
	prev := val.Prev
	next := val.Next
	prev.Next = next
	next.Prev = prev
}
