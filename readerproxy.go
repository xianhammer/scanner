package scanner

import "io"

type ReaderFunc func(b []byte) (n int, err error)
type OffsetFunc func() (n uint64)

type readerproxy struct {
	reader ReaderFunc
	offset OffsetFunc
}

func (rc *readerproxy) Read(b []byte) (n int, err error) {
	if rc.reader == nil {
		return 0, io.EOF
	}
	return rc.reader(b)
}

func (rc *readerproxy) Offset() (n uint64) {
	if rc.offset == nil {
		return 0
	}
	return rc.offset()
}
