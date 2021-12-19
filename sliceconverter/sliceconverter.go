package main

import (
	"reflect"
	"unsafe"
)

// sizeOf computes the size of a type using reflection.
//
// TODO(bga): Change this if we are ever allowed to use unsafe.Sizeof with
//            generic types.
func sizeOf[T any](t T) uintptr {
	return reflect.TypeOf(t).Size()
}

// convertSlice does a completelly unsafe conversion assuming that any
// validation has already been done.
func convertSlice[T1, T2 any](src []T1, dstCap, dstLen int) []T2 {
	var dst []T2

	// Adjust destination slice header.
	hs := (*reflect.SliceHeader)(unsafe.Pointer(&dst))
	hs.Data = uintptr(unsafe.Pointer(&src[0]))
	hs.Len = dstLen
	hs.Cap = dstCap

	return dst
}

// Convert does an in-place conversion of a slice of type T1 to a slice of type
// T2. The source and destination slices are validated before the conversion is
// done and, in case of issues, the code will panic.
func Convert[T1, T2 any](src []T1) []T2 {
	srcTypeSize := sizeOf(src[0])

	var dstTypeVar T2
	dstTypeSize := sizeOf(dstTypeVar)

	// The idea is to coerce the entire underlying array into the new slice so
	// we ignore the length and look at capacity instead. Note this wi;; result
	// in some extraneus entries in the slice when it is not fully used.
	srcCapInBytes := srcTypeSize * uintptr(cap(src))
	srcLenInBytes := srcTypeSize * uintptr(len(src))

	if srcCapInBytes%dstTypeSize != 0 {
		panic("Convert: src cap in bytes must be a multiple of the dst type size")
	}
	if srcLenInBytes%dstTypeSize != 0 {
		panic("Convert: src len in bytes must be a multiple of the dst type size")
	}

	dstCap := srcCapInBytes / dstTypeSize
	dstLen := srcLenInBytes / dstTypeSize

	return convertSlice[T1, T2](src, int(dstCap), int(dstLen))
}
