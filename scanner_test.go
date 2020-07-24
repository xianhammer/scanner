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
		input      string
		expectErr  error
	}{
		{0, "The quick brown fox jumps over the lazy dog.", nil},
		{17, "The quick brown fox jumps over the lazy dog.", nil},
		{12, "The quick brown fo jumps over the lazy dog.", nil},

		{0, "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.", nil},
		// {2, "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.", nil},
		{17, "The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.", nil},
		{12, "The quick brown fo jumps over the lazy dog. The quick brown fox ", io.EOF},
	}

	for testID, test := range tests {
		t.Logf("Test %2d\n", testID)

		inputReader := bytes.NewReader([]byte(test.input))
		scanner := New(inputReader, keyword)
		if test.bufferSize > 0 {
			scanner.buffer = make([]byte, test.bufferSize)
		}

		scannerBuffer := make([]byte, len(scannerExpect))
		err := scanner.Scan(func(r io.Reader) (err error) {
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
