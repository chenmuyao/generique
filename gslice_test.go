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

package generique

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGSliceSum(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var empty []float32

		_, err := Sum(empty)
		assert.ErrorIs(t, err, ErrEmptySlice)
	})

	t.Run("sum of ints", func(t *testing.T) {
		integers := []int64{1, 277865, 88934893312345}
		want := int64(1 + 277865 + 88934893312345)

		got, err := Sum(integers)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestGSliceMax(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var empty []float32

		_, err := Max(empty)
		assert.ErrorIs(t, err, ErrEmptySlice)
	})

	t.Run("max of ints", func(t *testing.T) {
		integers := []int64{1, 277865, 88934893312345}
		want := int64(88934893312345)

		got, err := Max(integers)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestGSliceMin(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var empty []float32

		_, err := Max(empty)
		assert.ErrorIs(t, err, ErrEmptySlice)
	})

	t.Run("min of ints", func(t *testing.T) {
		integers := []int64{1, 277865, -88934893312345}
		want := int64(-88934893312345)

		got, err := Min(integers)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestGSliceFind(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var empty []float32

		_, err := Find(empty, func(t float32) bool {
			return true
		})
		assert.ErrorIs(t, err, ErrNoValueFound)
	})

	t.Run("not found", func(t *testing.T) {
		integers := []int64{1, 277865, -88934893312345}

		_, err := Find(integers, func(t int64) bool {
			return t > 400000
		})
		assert.ErrorIs(t, err, ErrNoValueFound)
	})

	t.Run("found", func(t *testing.T) {
		integers := []int64{1, 277865, -88934893312345}
		want := int64(1)

		got, err := Find(integers, func(t int64) bool {
			return t < 400000
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestGSliceFindAll(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var empty []float32

		_, err := FindAll(empty, func(t float32) bool {
			return true
		})
		assert.ErrorIs(t, err, ErrNoValueFound)
	})

	t.Run("not found", func(t *testing.T) {
		integers := []int64{1, 277865, -88934893312345}

		_, err := FindAll(integers, func(t int64) bool {
			return t > 400000
		})
		assert.ErrorIs(t, err, ErrNoValueFound)
	})

	t.Run("found one", func(t *testing.T) {
		integers := []int64{1, 277865, -88934893312345}
		want := []int64{1}

		got, err := FindAll(integers, func(t int64) bool {
			return t == 1
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("found many", func(t *testing.T) {
		integers := []int64{1, 277865, -88934893312345}
		want := []int64{1, 277865, -88934893312345}

		got, err := FindAll(integers, func(t int64) bool {
			return t < 400000
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestGSliceInsert(t *testing.T) {
	t.Run("invalid index", func(t *testing.T) {
		var empty []float32

		_, err := Insert(1, 3.4, empty)
		assert.ErrorIs(t, err, ErrInvalidIndex[float32]{1, 0})

		_, err = Insert(-1, 3.4, empty)
		assert.ErrorIs(t, err, ErrInvalidIndex[float32]{-1, 0})
	})

	t.Run("insert empty", func(t *testing.T) {
		var empty []float32
		want := []float32{3.4}

		got, err := Insert(0, 3.4, empty)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	type myType struct {
		integer int
		str     string
	}

	testCases := []struct {
		Name      string
		BaseSlice []myType
		ToInsert  myType
		Index     int
		Want      []myType
	}{
		{
			"insert head",
			[]myType{
				{3, "test"},
				{2, "texto"},
			},
			myType{1, "tête"},
			0,
			[]myType{
				{1, "tête"},
				{3, "test"},
				{2, "texto"},
			},
		},
		{
			"insert middle",
			[]myType{
				{3, "test"},
				{2, "texto"},
			},
			myType{1, "tête"},
			1,
			[]myType{
				{3, "test"},
				{1, "tête"},
				{2, "texto"},
			},
		},
		{
			"insert tail",
			[]myType{
				{3, "test"},
				{2, "texto"},
			},
			myType{1, "tête"},
			2,
			[]myType{
				{3, "test"},
				{2, "texto"},
				{1, "tête"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := Insert(tc.Index, tc.ToInsert, tc.BaseSlice)
			assert.NoError(t, err)
			assert.Equal(t, tc.Want, got)
		})
	}
}

func testGSliceDelete(
	t *testing.T,
	deleteFn func(index int, vals []int) ([]int, error),
	funcName string,
) {
	t.Run(funcName+"invalid index", func(t *testing.T) {
		var empty []int

		_, err := deleteFn(1, empty)
		assert.ErrorIs(t, err, ErrInvalidIndex[int]{1, 0})

		_, err = deleteFn(-1, empty)
		assert.ErrorIs(t, err, ErrInvalidIndex[int]{-1, 0})
	})

	testCases := []struct {
		Name      string
		BaseSlice []int
		Index     int
		Want      []int
	}{
		{
			funcName + "delete head",
			[]int{3, 2, 1},
			0,
			[]int{2, 1},
		},
		{
			funcName + "delete middle",
			[]int{3, 2, 1},
			1,
			[]int{3, 1},
		},
		{
			funcName + "delete tail",
			[]int{3, 2, 1},
			2,
			[]int{3, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := deleteFn(tc.Index, tc.BaseSlice)
			assert.NoError(t, err)
			assert.Equal(t, tc.Want, got)
		})
	}
}

func TestGSliceDelete(t *testing.T) {
	testGSliceDelete(t, DeleteV1[int], "DeleteV1")
	testGSliceDelete(t, DeleteV2[int], "DeleteV2")
}

func benchmarkDelete(b *testing.B, size int, deleteFn func(index int, vals []int) ([]int, error)) {
	original := make([]int, size)
	for i := 0; i < size; i++ {
		original[i] = i
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		integers := make([]int, len(original))
		copy(integers, original)

		_, _ = deleteFn(b.N/2, integers)
	}
}

func BenchmarkGSliceDelete(b *testing.B) {
	sizes := []int{10, 100, 1000, 1000000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("DeleteV1-%d", size), func(b *testing.B) {
			benchmarkDelete(b, size, DeleteV1[int])
		})
		b.Run(fmt.Sprintf("DeleteV2-%d", size), func(b *testing.B) {
			benchmarkDelete(b, size, DeleteV2[int])
		})
		b.Run(fmt.Sprintf("DeleteUnordered-%d", size), func(b *testing.B) {
			benchmarkDelete(b, size, DeleteUnordered[int])
		})
	}
}

func TestGSliceDeleteShrinkV1(t *testing.T) {
	original := make([]int, 20)
	for i := 0; i < 20; i++ {
		original[i] = i
	}

	for i := 0; i < 11; i++ {
		original, _ = DeleteShrink(0, original, DeleteV1)
	}

	assert.Equal(t, 12, cap(original))
}
