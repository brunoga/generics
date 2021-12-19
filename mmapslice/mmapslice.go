package mmapslice

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"syscall"

	"github.com/brunoga/generics/sliceconverter"
)

type MMapSlice[T1 any] []T1

func (ms *MMapSlice[T]) Load(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	size := fi.Size()
	if size == 0 {
		return fmt.Errorf("tried to map an empty file")
	}

	if size < 0 {
		return fmt.Errorf("file %q has negative size", filename)
	}
	if size != int64(int(size)) {
		return fmt.Errorf("file %q is too large", filename)
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ|
		syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return err
	}

	*ms = sliceconverter.Convert[byte, T](data)

	runtime.SetFinalizer(ms, (*MMapSlice[T]).Close)

	return nil
}

func (ms MMapSlice[T]) Save(filename string) error {
	return ioutil.WriteFile(filename, sliceconverter.Convert[T, byte](ms), 0600)
}

func (ms *MMapSlice[T]) Close() error {
	runtime.SetFinalizer(ms, nil)

	err := syscall.Munmap(sliceconverter.Convert[T, byte](*ms))
	if err != nil {
		return err
	}

	return nil
}
