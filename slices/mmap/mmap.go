//go:build linux
// +build linux

package mmap

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"syscall"

	"github.com/brunoga/generics/slice/converter"
)

// MMapSlice is a slice of type T that is backed by a memory mapped file. You
// can read and write to it normally but using append() will have unexpected
// results if the slice must be resized (i.e. the data will not be the mmapped
// one anymore which basically turns it into a normal slice
type MMapSlice[T any] []T

// Map maps the file represented by the given filename to a MMapSlice and
// returns it. The file will be closed on return.
func Map[T any](filename string) (MMapSlice[T], error) {
	f, err := os.OpenFile(filename, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := fi.Size()
	if size == 0 {
		return nil, fmt.Errorf("file %q is empty", filename)
	}
	if size < 0 {
		return nil, fmt.Errorf("file %q has negative size", filename)
	}
	if size != int64(int(size)) {
		return nil, fmt.Errorf("file %q is too large", filename)
	}

	// Actually map the file to some of our address space.
	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ|
		syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	// Create a slice that will interpret the data in the mapped memory as
	// having the type we want. This does an in-place conversion so it works as
	// expected
	ms := sliceconverter.Convert[byte, T](data)

	// Make sure we unmap the memory when the slice goes out of scope.
	runtime.SetFinalizer(ms, (*MMapSlice[T]).UnMap)

	return ms, nil
}

// Save is a convenience function that saves the data associated with the given
// slice to the file represented by the given filename. This file can later be
// mapped (with the Map function above) to a MMapSlice.
func Save[T any](slice []T, filename string) error {
	return ioutil.WriteFile(filename, sliceconverter.Convert[T, byte](slice), 0600)
}

// UnMap unmaps the memory associated with the given MMapSlice. Any changes to
// the given MMapSlice that were not persisted to the file will now be.
func (ms *MMapSlice[T]) UnMap() error {
	runtime.SetFinalizer(ms, nil)

	// Unmap associated memory.
	err := syscall.Munmap(sliceconverter.Convert[T, byte](*ms))
	if err != nil {
		return err
	}

	// Point out slice to an brand new empty slice. We need this as setting ms
	// to nil wouild not quite work.
	*ms = make(MMapSlice[T], 0)

	return nil
}
