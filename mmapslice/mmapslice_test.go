package mmapslice

import (
	"os"
	"testing"
)

func TestMMapSlice(t *testing.T) {
	mmapSlice := MMapSlice[int]([]int{1, 2, 3})
	err := mmapSlice.Save("slice.mmap")
	if err != nil {
		t.Error(err)
	}

	defer func() { os.Remove("slice.mmap") }()

	mmapSlice2 := MMapSlice[int]{}
	err = mmapSlice2.Load("mmapslice_test")
	if err != nil {
		t.Error(err)
	}

	if len(mmapSlice) != len(mmapSlice2) {
		t.Error("mmapslice and mmapSlice2 have different lengths")
	}
	for i, v := range mmapSlice {
		if v != mmapSlice2[i] {
			t.Errorf("mmapSlice[%d] = %d, mmapSlice2[%d] = %d", i, v, i, mmapSlice2[i])
		}
	}
}
