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
