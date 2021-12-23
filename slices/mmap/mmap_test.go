package mmap

import (
	"os"
	"testing"
)

func TestMMapSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	err := Save(slice, "slice.mmap")
	if err != nil {
		t.Error(err)
	}

	defer func() { os.Remove("slice.mmap") }()

	mmapSlice, err := Map[int]("slice.mmap")
	if err != nil {
		t.Error(err)
	}
	if len(slice) != len(mmapSlice) {
		t.Error("original slice and mmapSlice have different lengths")
	}

	for i, v := range slice {
		if v != mmapSlice[i] {
			t.Errorf("slice[%d] = %d, mmapSlice[%d] = %d", i, v, i,
				mmapSlice[i])
		}
	}

	// Just make sure writing to the slice works.
	mmapSlice[0] = 4
	if mmapSlice[0] != 4 {
		t.Errorf("mmapSlice[0] = %d, expected 4", mmapSlice[0])
	}
}
