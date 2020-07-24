package scanner

import (
	"bytes"
	"io"
	"testing"
)

func TestReaderProxyEmpty(t *testing.T) {
	rp := new(readerproxy)

	b := make([]byte, 10)
	n, err := rp.Read(b)
	if n != 0 {
		t.Errorf("Read: Expected %d bytes, got %d\n", 0, n)
	}
	if err != io.EOF {
		t.Errorf("Read: Expected error %v, got %v\n", io.EOF, err)
	}
}

func TestReaderProxyReceiver(t *testing.T) {
	expect := "tester"
	rp := new(readerproxy)
	rp.reader = bytes.NewBufferString(expect).Read

	b := make([]byte, 10)
	n, err := rp.Read(b)
	if n != len(expect) {
		t.Errorf("Read: Expected %d bytes, got %d\n", len(expect), n)
	}

	if false == bytes.Equal(b[:n], []byte(expect)) {
		t.Errorf("Read: Expected %c, got %c\n", []byte(expect), b[:n])
	}

	if err != nil {
		t.Errorf("Read: Expected error %v, got %v\n", io.EOF, err)
	}
}
