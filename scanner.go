package scanner

import (
	"bytes"
	"io"
)

const DefaultBufferSize = 10 * 1024

type scanner struct {
	r       io.Reader
	keyword []byte
	buffer  []byte
}

// New construct a new scanner object
func New(r io.Reader, keyword []byte) (s *scanner) {
	s = new(Scanner)
	s.r = r
	s.keyword = make([]byte, len(keyword))
	copy(s.keyword, keyword)
	return
}

func (s *scanner) Buffer(size int) {
	if size <= 1024 {
		size = DefaultBufferSize
	}
	s.buffer = make([]byte, size)
}

// Scan return io.EOF if the underlying reader is at EOF.
// If nil error is returned r contain a reader from the start of the keyword given.
func (s *scanner) Scan(keyword []byte, callback func(r io.Reader) (cerr error)) (err error) {
	if s.buffer == nil {
		s.Buffer(0)
	}

	offset := 0
	bufferEnd := 0
	fillBuffer := func() (err error) {
		bufferEnd, err = s.r.Read(s.buffer)
		offset = 0
		return
	}

	doCallback := func() {
		readBuffer := func(b []byte) (n int, e error) {
			n = copy(b, s.buffer[offset:bufferEnd])
			offset += n
			if offset >= bufferEnd {
				e = fillBuffer()
			}
			return
		}

		keywordOffset := 0 // Allow for mulitple (short receiver) reads from keyword
		rp := new(readerproxy)
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

	err = fillBuffer()

	for err == nil {
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
