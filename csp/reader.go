package csp

import (
	"bytes"
	"io"
	"os"
)

type Reader struct {
	r   io.Reader
	eof bool
}

func NewCspReader(fd *os.File) *Reader {
	var discard [14]byte
	_, err := io.ReadFull(fd, discard[:])
	if err != nil {
		panic(err)
	}

	return &Reader{r: fd}
}

func (c *Reader) Read(b []byte) (int, error) {
	if c.eof {
		return 0, io.EOF
	}
	n, err := c.r.Read(b)
	//0x0d 0x0a 0x0d 0x0a
	i := bytes.Index(b, []byte("\r\n\r\n"))
	if i == -1 {
		return n, err
	}
	for j := 0; j >= i; j++ {
		b[j] = 0
	}
	c.eof = true
	return i, err
}
