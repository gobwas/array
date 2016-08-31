// Package array contains immutable sorted array implementation.
// This data structure is useful when you do many read operations and much less writes.
// See https://en.wikipedia.org/wiki/Sorted_array
package array

import (
	"sort"
)

// Immutable sorted array.
type Array struct {
	data []Item
}

// Item represents single element in array.
type Item interface {
	// Less tests whether the current item is less than the given argument.
	//
	// This must provide a strict weak ordering.
	// If !a.Less(b) && !b.Less(a), then this means a == b.
	Less(than Item) bool
}

// Int implements Item interface for ints.
type Int int

// Less tests that a is less than b.
func (a Int) Less(b Item) bool {
	return a < b.(Int)
}

// Has checks whether array contains x.
func (a Array) Has(x Item) (ok bool) {
	ok, _ = a.search(x)
	return
}

// Get returns item x if it exists.
func (a Array) Get(x Item) Item {
	ok, i := a.search(x)
	if !ok {
		return nil
	}
	return a.data[i]
}

// Upsert inserts x into new copy of Array if it not exists yet.
// If such item is already exists, it replaces it with x in returned copy.
// If item was replaced, previous item is returned. Otherwise prev is nil.
func (a Array) Upsert(x Item) (cp Array, prev Item) {
	var with []Item
	has, n := a.search(x)
	if has {
		with = make([]Item, len(a.data))
		copy(with, a.data)
		a.data[n], prev = x, a.data[n]
	} else {
		with = make([]Item, len(a.data)+1)
		copy(with[:n], a.data[:n])
		copy(with[n+1:], a.data[n:])
		with[n] = x
	}
	return Array{with}, prev
}

// Delete deletes x from new copy of array.
// It returns item that was deleted from array.
// If item was not found, prev is nil.
func (a Array) Delete(x Item) (cp Array, prev Item) {
	has, n := a.search(x)
	if !has {
		return a, nil
	}

	without := make([]Item, len(a.data)-1)
	copy(without[:n], a.data[:n])
	copy(without[n:], a.data[n+1:])

	return Array{without}, a.data[n]
}

// Ascend calls the cb for every array's item until it returns false.
func (a Array) Ascend(cb func(x Item) bool) {
	for _, x := range a.data {
		if !cb(x) {
			break
		}
	}
}

// Len returns length of the underlying data.
func (a Array) Len() int {
	return len(a.data)
}

// search searches x in Array's data using binary search.
// It returns true if it found exact same element.
// It also returns index in l.data, where x could lay.
func (a Array) search(x Item) (ok bool, i int) {
	i = sort.Search(len(a.data), func(j int) bool {
		return !a.data[j].Less(x)
	})
	if i == len(a.data) {
		return
	}
	f := a.data[i]
	ok = !f.Less(x) && !x.Less(f)
	return
}
