package util

import (
	"errors"
	"math/rand"
)

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

func SortBy[K any](ks []K, fn func(this K, other K) bool) []K {
	if len(ks) < 2 {
		return ks
	}

	left, right := 0, len(ks)-1

	// Pick a pivot
	pivotIndex := rand.Int() % len(ks)

	// Move the pivot to the right
	ks[pivotIndex], ks[right] = ks[right], ks[pivotIndex]

	// Pile elements smaller than the pivot on the left
	for i := range ks {
		if fn(ks[i], ks[right]) {
			ks[i], ks[left] = ks[left], ks[i]
			left++
		}
	}

	// Place the pivot after the last smaller element
	ks[left], ks[right] = ks[right], ks[left]

	// Go down the rabbit hole
	SortBy(ks[:left], fn)
	SortBy(ks[left+1:], fn)

	return ks
}
