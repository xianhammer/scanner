package scanner

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Blob struct {
	Start, Size, End uint64
	Create           func(b *Blob) (f *os.File, err error)
	Filename         func(b *Blob) (name string)
}

func NewBlob(address uint64) (b *Blob) {
	b = new(Blob)
	b.Start = address
	return
}

func (b *Blob) GetStart() uint64 {
	return b.Start
}

func (b *Blob) GetEnd() uint64 {
	if b.End == 0 {
		return b.Start + b.Size
	}
	return b.End
}

func (b *Blob) GetSize() (size uint64) {
	if b.End == 0 || b.Size != 0 {
		return b.Size
	}

	size = b.End - b.Start
	if size > b.Size {
		size = b.Size
	}
	return
}

func (b *Blob) GetFilename() (fn string) {
	if b.Filename == nil {
		return fmt.Sprintf("%d.bin", b.Start)
	}
	return b.Filename(b)
}

func (b *Blob) CreateFile() (f *os.File, err error) {
	if b.Create == nil {
		return os.Create(b.GetFilename())
	}
	return b.Create(b)
}

func (b *Blob) WriteBlob(r io.ReadSeeker, bufferSize int) (err error) {
	// Support reading blobs from block devices, require mulitpla of block size.
	remainder := b.Start % 0x0200
	fileoffset := b.Start - remainder

	_, err = r.Seek(int64(fileoffset), io.SeekStart)
	if err != nil {
		log.Printf("Skipping [%x], seek error = %v\n", b.GetStart(), err)
		return
	}

	w, err := b.CreateFile()
	if err != nil {
		log.Printf("Skipping [%x], create file error = %v\n", b.GetStart(), err)
		return
	}
	defer w.Close()

	size := int64(b.GetSize())
	log.Printf("Extracting from %x to %x (%x bytes) to %s\n", b.GetStart(), b.GetEnd(), b.GetSize(), b.GetFilename())
	writeBuffer := make([]byte, int(bufferSize))

	var n, n0 int
	for err == nil && size > 0 {
		n, err = r.Read(writeBuffer[:])
		if err != nil && err != io.EOF {
			break
		}

		size -= int64(n)
		log.Printf("* writing %d bytes, %d bytes remaining, err=%v", n, size, err)
		if size < 0 {
			size += int64(n) // Go positive again
			n = int(size)    // Limit n for output
			size = 0
		}

		var errW error
		for offset := int(remainder); errW == nil && offset < n; {
			n0, errW = w.Write(writeBuffer[offset:n])
			log.Printf("- wrote %d bytes of %d, from %d to %d, err=%v", n0, n, offset, offset+n0, errW)
			if errW == nil {
				offset += n0
			}
		}
		if errW != nil && errW != io.EOF {
			err = errW
		}

		remainder = 0
	}
	w.Sync()

	log.Printf("remaining size=%d, err=%v", size, err)
	return
}
