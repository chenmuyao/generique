// MIT License
//
// Copyright (c) 2024 chenmuyao
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gslice

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum returns the sum of a slice
func Sum[T Number](vals []T) (T, error) {
	if len(vals) == 0 {
		var zero T
		return zero, ErrEmptySlice
	}

	var res T
	for _, i := range vals {
		res += i
	}

	return res, nil
}

// Sum returns the maximum value of a slice
func Max[T Number](vals []T) (T, error) {
	n := len(vals)
	if n == 0 {
		var zero T
		return zero, ErrEmptySlice
	}
	t := vals[0]
	for i := 1; i < n; i++ {
		if t < vals[i] {
			t = vals[i]
		}
	}
	return t, nil
}

// Sum returns the minimum value of a slice
func Min[T Number](vals []T) (T, error) {
	n := len(vals)
	if n == 0 {
		var zero T
		return zero, ErrEmptySlice
	}
	t := vals[0]
	for i := 1; i < n; i++ {
		if t > vals[i] {
			t = vals[i]
		}
	}
	return t, nil
}

// Sum returns the first element that is found by the filter function.
func Find[T any](vals []T, filter func(t T) bool) (T, error) {
	for _, v := range vals {
		if filter(v) {
			return v, nil
		}
	}
	var zero T
	return zero, ErrNoValueFound
}

// Sum returns the all elements that are found by the filter function.
func FindAll[T any](vals []T, filter func(t T) bool) ([]T, error) {
	var results []T

	for _, v := range vals {
		if filter(v) {
			results = append(results, v)
		}
	}

	if len(results) == 0 {
		return nil, ErrNoValueFound
	}

	return results, nil
}

// Insert inserts an element at the given index. Return the same slice if
// not reallocated.
func Insert[T any](index int, val T, vals []T) ([]T, error) {
	if index < 0 || index > len(vals) {
		var zero []T
		return zero, ErrInvalidIndex[T]{Index: index, Size: len(vals)}
	}

	vals = append(vals, val)

	for i := len(vals) - 1; i > index; i-- {
		if i-1 >= 0 {
			vals[i] = vals[i-1]
		}
	}
	vals[index] = val

	return vals, nil
}

// DeleteV1 appends 2 sections of slices
func DeleteV1[T any](index int, vals []T) ([]T, error) {
	n := len(vals)
	if index < 0 || index >= n {
		var zero []T
		return zero, ErrInvalidIndex[T]{Index: index, Size: n}
	}

	vals = append(vals[:index], vals[index+1:]...)
	return vals, nil
}

// DeleteV1 moves memory
func DeleteV2[T any](index int, vals []T) ([]T, error) {
	n := len(vals)
	if index < 0 || index >= n {
		var zero []T
		return zero, ErrInvalidIndex[T]{Index: index, Size: n}
	}

	for i := index; i < n; i++ {
		if i+1 < n {
			vals[i] = vals[i+1]
		}
	}

	vals = vals[:n-1]
	return vals, nil
}

func DeleteUnordered[T any](index int, vals []T) ([]T, error) {
	n := len(vals)
	if index < 0 || index >= n {
		var zero []T
		return zero, ErrInvalidIndex[T]{Index: index, Size: n}
	}

	vals[index] = vals[n-1]
	return vals[:n-1], nil
}

func DeleteShrink[T any](
	index int,
	vals []T,
	deletFn func(idx int, vals []T) ([]T, error),
) ([]T, error) {
	// Shrinking means memory copy. So we shrink only when the size is a half
	// of the original slice. And we keep 25% capacity for future growth

	vals, err := deletFn(index, vals)
	if err != nil {
		return vals, nil
	}

	oldLen := len(vals)
	oldCap := cap(vals)

	if oldLen < oldCap/2 {
		newCap := int(float64(oldCap/2) * 1.25)
		newVals := make([]T, oldLen, newCap)
		copy(newVals, vals)
		return newVals, nil
	}

	return vals, nil
}

func Map[In any, Out any](s []In, mapper func(id int, src In) Out) []Out {
	res := make([]Out, len(s))
	for i, el := range s {
		o := mapper(i, el)
		res[i] = o
	}
	return res
}
