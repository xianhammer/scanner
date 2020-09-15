package scanner

import (
	"bufio"
	"io"
	"os"
)

// FileLines return all lines from the given file.
func FileLines(filename string, callback func(line string) error) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	return Lines(file, callback)
}

// Lines return all lines from the given reader.
func Lines(r io.Reader, callback func(line string) error) (err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		err = callback(scanner.Text())
		if err != nil {
			break
		}
	}
	return scanner.Err()
}
