// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package collections defines operations on collections of items.
package collections

// ============================================================================
// SPDX-SnippetBegin
// SPDX-FileCopyrightText: 2018 Christopher James <https://github.com/quii>
// SPDX-License-Identifier: MIT

// Find is a higher-order function that iterates through a collection to find
// a value. If the value is found, it returns the value and a boolean
// that signifies the value exists in the collection. Otherwise, the zero
// value of the value is returned along with false.
func Find[A any](collection []A, finder func(A, A) bool, target A) (A, bool) {
	for _, item := range collection {
		if finder(item, target) {
			return item, true
		}
	}

	var zero A
	return zero, false
}

// Reduce takes a collection of elements A and applies a reduction function
// on the elements, reducing the collection into a single value of type B.
func Reduce[A, B any](collection []A, reductionFunc func(B, A) B, initialValue B) B {
	result := initialValue
	for _, item := range collection {
		result = reductionFunc(result, item)
	}
	return result
}

// SPDX-SnippetEnd
// ============================================================================

// Filter takes a collection of elements and filters a value out according to
// a filter function. The collection returned is a new collection of type A.
func Filter[A any](collection []A, filter func(target, item A) bool, target A) []A {
	var results []A
	for _, item := range collection {
		if filter(target, item) {
			results = append(results, item)
		}
	}
	return results
}
