package scanner

import (
	"bytes"
	"io"
)

// DefaultBufferSize used by the scanner.
// Default buffer size is set to 4K matching current standard block device block size.
const DefaultBufferSize = 4 * 1024

type Offset interface {
	Offset() uint64
}

type scanner struct {
	r       io.Reader
	keyword []byte
	buffer  []byte
}

// New construct a new scanner object
func New(r io.Reader, keyword []byte) (s *scanner) {
	s = new(scanner)
	s.r = r
	s.keyword = make([]byte, len(keyword))
	copy(s.keyword, keyword)
	return
}

// Buffer reset buffer size. If size is below 1K default size is used.
func (s *scanner) Buffer(size int) {
	if size <= 1024 {
		size = DefaultBufferSize
	}
	s.buffer = make([]byte, size)
}

// Scan a source reader for keyword until non-nil is returned.
//
// For each found keyword call the callback returning a reader into the source reader.
// If callback return an error, the scanner will stop and return that error.
//
// It is possible to read any number of bytes (including the keyword bytes) - even until
// source reader return non-nil.
// Keyword scan will resume after the last byte read by callback.
//
// E.g. source = [1, 2, 3, 4, 5, 1, 2, 3] and keyword = [3]
// Callback receives:
//    reader on [3, 4, 5, 1, 2, 3] -> callback read [3, 4] -> scanner resume on [5, 1, 2, 3]
//    reader on [3]                -> callback read [3] and EOF -> scanner return EOF
//
// Notice that the reader handed to callback rely directly on the source reader and can
// therefore not be used asynchronously.
func (s *scanner) Scan(callback func(r io.Reader) (cerr error)) (err error) {
	if s.buffer == nil {
		s.Buffer(0)
	}

	var address uint64 // Offset since start of source reader.

	offset := 0    // Offset into buffer
	bufferEnd := 0 // Size of last read
	rp := new(readerproxy)
	rp.offset = func() (n uint64) {
		return address + uint64(offset-len(s.keyword))
	}

	fillBuffer := func() (err error) {
		address += uint64(bufferEnd)
		bufferEnd, err = s.r.Read(s.buffer)
		offset = 0
		return
	}

	readBuffer := func(b []byte) (n int, e error) {
		n = copy(b, s.buffer[offset:bufferEnd])
		offset += n
		if offset >= bufferEnd {
			e = fillBuffer()
		}
		return
	}

	doCallback := func() {
		keywordOffset := 0 // Allow for multiple (short receiver) reads from keyword

		rp.reader = func(b []byte) (n int, e error) {
			n = copy(b, s.keyword[keywordOffset:])
			keywordOffset += n
			if keywordOffset >= len(s.keyword) {
				rp.reader = readBuffer
			}
			return
		}

		err = callback(rp)
	}

	for err = fillBuffer(); err == nil; {
		idx := bytes.Index(s.buffer[offset:bufferEnd], s.keyword)
		if idx >= 0 {
			offset += idx + len(s.keyword)
			doCallback()
			continue
		}

		offset = bufferEnd - len(s.keyword)
		idx = bytes.IndexByte(s.buffer[offset:bufferEnd], s.keyword[0])
		if idx < 0 {
			err = fillBuffer()
			continue
		}

		hasPartial := bytes.HasPrefix(s.keyword, s.buffer[offset+idx:bufferEnd])
		err = fillBuffer()
		if err == nil && hasPartial && bytes.HasPrefix(s.buffer, s.keyword[len(s.keyword)-idx:]) {
			offset = len(s.keyword) - idx // Skip remaining part of compared keyword bytes
		}
	}
	return
}

// func Scan(r io.Reader, keyword []byte, callback func(r io.Reader) (cerr error)) (err error) {
// 	return New(r, keyword).Scan(callback)
// }
