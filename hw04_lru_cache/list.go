package hw04lrucache

import "fmt"

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
	key   Key
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front    *ListItem
	back     *ListItem
	elements map[Key]*ListItem
}

func (l *list) Len() int {
	return len(l.elements)
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	var value interface{}
	var key Key

	switch i := v.(type) {
	case Pair:
		value = i.value
		key = i.key
	default:
		value = v
	}

	k := genKey(value, key)
	val, ok := l.elements[k]
	if !ok {
		if l.front == nil {
			l.front = &ListItem{key: k, Value: value, Next: nil, Prev: nil}
		} else {
			temp := l.front
			l.front = &ListItem{key: k, Value: value, Next: temp, Prev: nil}
			temp.Prev = l.front
			if temp.Next == nil {
				l.back = temp
			}
		}
		l.elements[k] = l.front
		return l.front
	}
	return val
}

func (l *list) PushBack(v interface{}) *ListItem {
	var value interface{}
	var key Key

	switch i := v.(type) {
	case Pair:
		value = i.value
		key = i.key
	default:
		value = v
	}

	k := genKey(value, key)
	val, ok := l.elements[k]
	if !ok {
		if l.front == nil {
			return l.PushFront(v)
		}
		if l.back == nil {
			l.back = &ListItem{key: k, Value: value, Next: nil, Prev: l.front}
			l.front.Next = l.back
		} else {
			temp := l.back
			l.back = &ListItem{key: k, Value: value, Next: nil, Prev: temp}
			temp.Next = l.back
		}
		l.elements[k] = l.back
		return l.back
	}

	return val
}

func (l *list) Remove(i *ListItem) {
	val, ok := l.elements[i.key]
	if ok {
		switch val {
		case l.front:
			l.front = l.front.Next
			l.front.Prev = nil
		case l.back:
			l.back = l.back.Prev
			l.back.Next = nil
		default:
			changeLinks(val)
		}
		delete(l.elements, i.key)
	}
}

func (l *list) MoveToFront(i *ListItem) {
	val, ok := l.elements[i.key]
	if ok {
		if val == l.front {
			return
		}
		if val == l.back {
			newBack := val.Prev
			temp := l.front
			temp.Prev = val
			val.Next = temp
			l.front = val
			l.front.Prev = nil
			l.back = newBack
			l.back.Next = nil
		} else {
			changeLinks(val)
			temp := l.front
			temp.Prev = val
			val.Next = temp
			l.front = val
		}
	}
}

func NewList() List {
	return &list{front: nil, back: nil, elements: make(map[Key]*ListItem)}
}

func genKey(v interface{}, k Key) Key {
	if len(k) == 0 {
		return Key(fmt.Sprintf("%v", v))
	}
	return k
}

func changeLinks(val *ListItem) {
	prev := val.Prev
	next := val.Next
	prev.Next = next
	next.Prev = prev
}
