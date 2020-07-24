package scanner

import "io"

type ReaderFunc func(b []byte) (n int, err error)

type readerproxy struct {
	reader ReaderFunc
}

func (rc *readerproxy) Read(b []byte) (n int, err error) {
	if rc.reader == nil {
		return 0, io.EOF
	}
	return rc.reader(b)
}
