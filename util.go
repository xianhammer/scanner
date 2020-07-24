package scanner

import (
	"io"
	"os"
)

// File scans a named file and call the callback when keyword is found.
func File(filepath string, keyword []byte, callback func(address uint64, r io.Reader) error) (err error) {
	return FileAf(filepath, 0, keyword, callback)
}

// FileAt scans a named file from the given offset and call the callback when keyword is found.
func FileAt(filepath string, offset int64, keyword []byte, callback func(address uint64, r io.Reader) error) (err error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0600)
	if err != nil {
		return
	}
	defer file.Close()

	if offset > 0 {
		if _, err = file.Seek(offset, os.SEEK_SET); err != nil {
			return
		}
	}

	return New(file, keyword).Scan(callback)
}
