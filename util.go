package scanner

import (
	"io"
	"os"
)

// File scans a named file and call the callback when keyword is found.
func File(pathname string, keyword []byte, callback func(address uint64, r io.Reader) error) (err error) {
	return FileAt(pathname, 0, keyword, callback)
}

// FileAt scans a named file from the given offset and call the callback when keyword is found.
func FileAt(pathname string, offset int64, keyword []byte, callback func(address uint64, r io.Reader) error) (err error) {
	file, err := os.OpenFile(pathname, os.O_RDONLY, 0600)
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

// Read from reader r and call the callback for each new block. Recycle mark how much of the old buffer end should be
// copied to front of the buffer for next read. This is usefull for e.g. simplyfing byte search in files.
// Smallest recycle value is 128 and bufferSize must be at least twice that - if not, it is set to 8K.
func Read(r io.Reader, bufferSize, recycle int, callback func(address uint64, b []byte) int) (err error) {
	if recycle < 128 {
		recycle = 128
	}
	if bufferSize <= 2*recycle {
		bufferSize = 8 * 1024
	}
	buffer := make([]byte, bufferSize)

	var n int
	n, err = r.Read(buffer)
	if err != nil {
		return
	}

	var offset uint64
	for err == nil {
		callback(offset, buffer[:n])

		offset += uint64(n)
		if n > recycle {
			copy(buffer, buffer[n-recycle+1:n])
		}
		n, err = r.Read(buffer[recycle:])
	}

	if err == io.EOF {
		err = nil
	}
	return
}
