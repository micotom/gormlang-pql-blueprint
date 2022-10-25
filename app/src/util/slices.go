package util

import (
	"errors"

	"golang.org/x/exp/constraints"
)

func Map[K any, T any](ks []K, m func(k K) T) []T {
	s := []T{}
	for _, k := range ks {
		s = append(s, m(k))
	}
	return s
}

func GroupBy[K any, T comparable](ks []K, fn func(k K) T) map[T][]K {
	m := make(map[T][]K)
	for _, k := range ks {
		t := fn(k)
		if val, present := m[t]; present {
			val = append(val, k)
			m[t] = val
		} else {
			m[t] = []K{}
			m[t] = append(m[t], k)
		}
	}
	return m
}

func SumByInt[K any](ks []K, sum func(k K) int) int {
	var s = 0
	for _, k := range ks {
		s += sum(k)
	}
	return s
}

func FindBy[K any](ks []K, fn func(k K) bool) (*K, error) {
	for _, k := range ks {
		if fn(k) {
			return &k, nil
		}
	}
	return nil, errors.New("No such element")
}

// Sorts the slice by determine the order based on the passede function for that specific element of the slice. Sorting algorithm is quicksort.
func SortBy2[T any, V constraints.Ordered](slice []T, fn func(t T) V) []T {

	var sort func([]T, int, int)
	sort = func(arr []T, start int, end int) {
		if (end - start) < 1 {
			return
		}

		pivot := arr[end]
		splitIndex := start
		for i := start; i < end; i++ {
			if fn(arr[i]) < (fn(pivot)) {
				temp := arr[splitIndex]
				arr[splitIndex] = arr[i]
				arr[i] = temp
				splitIndex++
			}
		}
		arr[end] = arr[splitIndex]
		arr[splitIndex] = pivot

		sort(arr, start, splitIndex-1)
		sort(arr, splitIndex+1, end)
	}

	r := make([]T, len(slice))

	for i, t := range slice {
		r[i] = t
	}

	sort(r, 0, len(r)-1)

	return r
}

func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func Fold[T any, V any](slice []T, initial V, fn func(acc V, t T) V) V {
	for _, t := range slice {
		initial = fn(initial, t)
	}
	return initial
}

func Filter[T any](slice []T, fn func(t T) bool) []T {
	foldFn := func(acc []T, t T) []T {
		if fn(t) {
			return append(acc, t)
		} else {
			return acc
		}
	}
	return Fold(slice, []T{}, foldFn)
}
