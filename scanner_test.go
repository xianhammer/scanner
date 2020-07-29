package scanner

import (
	"bytes"
	"io"
	"testing"
)

func TestScanner001(t *testing.T) {
	keyword := []byte("fox")
	scannerExpect := []byte("fox jumps ")

	tests := []struct {
		bufferSize int
		address    []int
		input      string
		expectErr  error
	}{
		{0, []int{16}, "The quick brown fox jumps over the lazy dog.", nil},
		{17, []int{16}, "The quick brown fox jumps over the lazy dog.", nil},
		{12, []int{16}, "The quick brown fo jumps over the lazy dog.", nil},

		{0, []int{16, 61}, "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.", nil},
		// {2, "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.", nil},
		{17, []int{16, 61}, "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.", nil},
		{12, []int{16, 61}, "The quick brown fo jumps over the lazy dog. The quick brown fox ", io.EOF},
	}

	for testID, test := range tests {
		t.Logf("Test %2d\n", testID)

		inputReader := bytes.NewReader([]byte(test.input))
		s := New(inputReader, keyword)
		if test.bufferSize > 0 {
			s.buffer = make([]byte, test.bufferSize)
		}

		idxAddress := 0
		scannerBuffer := make([]byte, len(scannerExpect))
		err := s.Scan(func(address uint64, r io.Reader) (err error) {
			if address != test.address[idxAddress] {
				t.Errorf("%2d: scanner: Unexpected address argument (%d), expected %d\n", testID, address, test.address[idxAddress])
			}

			n, err := fillBuffer(t, r, scannerBuffer)
			if err != test.expectErr {
				t.Errorf("%2d: scanner: Expected error %v, got %v\n", testID, test.expectErr, err)
				return
			}
			if err == io.EOF {
				return
			}
			if n != len(scannerExpect) {
				t.Errorf("%2d: scanner: Unexpected number of bytes returned (%d), expected %d\n", testID, n, len(scannerExpect))
			}
			if !bytes.Equal(scannerBuffer, scannerExpect) {
				t.Errorf("%2d: scanner: Expected %c got %c\n", testID, scannerExpect, scannerBuffer)
			}
			return
		})

		if err == io.EOF {
			err = nil
		}
		if err != nil {
			t.Errorf("%2d: Unexpected error %v\n", testID, err)
		}
	}

	// t.Errorf("STOP\n")
}

func fillBuffer(t *testing.T, r io.Reader, buffer []byte) (n int, err error) {
	// t.Logf("fillBuffer: len(buffer) = %d\n", len(buffer))
	for err == nil && n < len(buffer) {
		n0, e0 := r.Read(buffer[n:])
		// t.Logf("fillBuffer-loop: n0=%d, e0=%v\n", n0, e0)
		if e0 != nil {
			err = e0
		} else if n0 == 0 {
			err = io.EOF
		}
		n += n0
	}

	return
}
