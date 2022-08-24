// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"bytes"
	"sort"

	"github.com/ava-labs/avalanchego/utils/hashing"
	"golang.org/x/exp/constraints"
)

// TODO can we handle sorting where the Less function relies on a codec?

type Sortable[T any] interface {
	Less(T) bool
}

type Hashable interface {
	~[]byte
}

// TODO add tests
func SortSliceSortable[T Sortable[T]](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Less(s[j])
	})
}

// Sorts a slice of elements that satisfy constraints.Ordered.
// TODO add tests
func SortSliceOrdered[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

// TODO test
func SortByHash[T Hashable](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return bytes.Compare(hashing.ComputeHash256(s[i]), hashing.ComputeHash256(s[j])) == -1
	})
}

// Sort2DByteSlice sorts a 2D byte slice.
// Each byte slice is not sorted internally; the byte slices are sorted relative to another.
func Sort2DByteSlice[T Hashable](arr []T) {
	sort.Slice(
		arr,
		func(i, j int) bool {
			return bytes.Compare(arr[i], arr[j]) == -1
		})
}

// TODO add tests
func IsSortedAndUniqueSortable[T Sortable[T]](s []T) bool {
	for i := 0; i < len(s)-1; i++ {
		if !s[i].Less(s[i+1]) {
			return false
		}
	}
	return true
}

// Returns true iff the elements in [s] are unique and sorted.
func IsSortedAndUniqueOrdered[T constraints.Ordered](s []T) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] >= s[i+1] {
			return false
		}
	}
	return true
}

// TODO test
func IsSortedAndUniqueByHash[T Hashable](s []T) bool {
	for i := 0; i < len(s)-1; i++ {
		if bytes.Compare(hashing.ComputeHash256(s[i]), hashing.ComputeHash256(s[i+1])) != -1 {
			return false
		}
	}
	return true
}

// IsSorted2DByteSlice returns true iff [arr] is sorted
func IsSorted2DByteSlice[T Hashable](arr []T) bool {
	return sort.SliceIsSorted(
		arr,
		func(i, j int) bool {
			return bytes.Compare(arr[i], arr[j]) == -1
		})
}

// IsSortedAndUnique returns true if the elements in the data are unique and sorted.
func IsSortedAndUnique(data sort.Interface) bool {
	for i := 0; i < data.Len()-1; i++ {
		if !data.Less(i, i+1) {
			return false
		}
	}
	return true
}
