package converter

import (
	"math"
	"reflect"
	"testing"
	"unsafe"
)

func Test_sizeOf(t *testing.T) {
	if got := sizeOf(int(0)); got != unsafe.Sizeof(int(0)) {
		t.Errorf("sizeOf() = %v, want %v", got, unsafe.Sizeof(int(0)))
	}
	if got := sizeOf(float64(0)); got != unsafe.Sizeof(float64(0)) {
		t.Errorf("sizeOf() = %v, want %v", got, unsafe.Sizeof(float64(0)))
	}
	if got := sizeOf(byte(0)); got != unsafe.Sizeof(byte(0)) {
		t.Errorf("sizeOf() = %v, want %v", got, unsafe.Sizeof(byte(0)))
	}
}

func Test_convertSlice(t *testing.T) {
	src := []byte{'0', '1', '2', '3'}

	dst := convertSlice[byte, int](src, 1, 1)

	hsSrc := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	hsDst := (*reflect.SliceHeader)(unsafe.Pointer(&dst))

	if hsSrc.Data != hsDst.Data {
		t.Errorf("convertSlice() = SliceHeader.Data=%v, want %v", hsSrc.Data, hsDst.Data)
	}
	if hsDst.Len != 1 {
		t.Errorf("convertSlice() = SliceHeader.Len=%v, want 1", hsDst.Len)
	}
	if hsDst.Cap != 1 {
		t.Errorf("convertSlice() = SliceHeader.Cap=%v, want 1", hsDst.Cap)
	}
}

func expectPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Convert() = panic was expected and did not happen")
		}
	}()

	f()
}

func TestConvert(t *testing.T) {
	type type3bytes [3]byte

	expectPanic(t, func() {
		Convert[type3bytes, int]([]type3bytes{})
	})
	expectPanic(t, func() {
		Convert[byte, int](make([]byte, 8, 9))
	})

	dst := Convert[uint64, byte]([]uint64{math.MaxUint64})
	if len(dst) != 8 {
		t.Errorf("Convert() = len(dst)=%v, want 8", len(dst))
	}
	if cap(dst) != 8 {
		t.Errorf("Convert() = cap(dst)=%v, want 8", len(dst))
	}
	for i := 0; i < 8; i++ {
		if dst[i] != 255 {
			t.Errorf("Convert() = dst[%v]=%v, want 255", i, dst[i])
		}
	}
}
