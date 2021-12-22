# slices/mmap
Memory mapping files to arbitrary slice types.

## Usage:

src := []MyFancyStruct{ ... }

err := mmap.Save(src, "file.mmap")
...
err, mMappedSlice := mmap.Map("file.mmap")
...
err := mMappedSlice.Close()
