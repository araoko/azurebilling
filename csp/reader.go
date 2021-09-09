package csp

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type Reader struct {
	r       io.Reader
	summary bool
	eof     bool
}

func NewCspReader(fd *os.File, summary bool) (*Reader, error) {
	var err error
	switch summary {
	case true:
		err = skipBOMSummaryHeading(fd)

	case false:
		err = skipSummary(fd)
	}
	if err != nil {
		return nil, err
	}

	return &Reader{r: fd, summary: summary}, nil
}

//func skipBOMSummaryHeading(fd *os.File) {
//	//Byte Order Mark (BOM) and summary headiing is 14 bytes
//	var discard [14]byte
//	_, err := io.ReadFull(fd, discard[:])
//	if err != nil {
//		panic(err)
//	}
//}

const maxSearchByte = 15

func skipBOMSummaryHeading(fd *os.File) error {
	//locate the first newline byte and skip to after it
	bb := make([]byte, 1)
	var offset int
	var err error
	for {
		if offset > maxSearchByte {
			return fmt.Errorf("could not skipSummary beginig of csv data after %d bytes. Fine invalid", maxSearchByte)
		}
		_, err = fd.Read(bb)
		if err != nil {
			return err
		}
		offset++
		if bb[0] == 10 {
			break
		}
	}

	return nil
}

const chunkSize = 4096

func skipSummary(fd *os.File) error {
	var offset int64
	search := []byte("Daily Usage")
	chunk := make([]byte, chunkSize+len(search))
	for {
		n, err := fd.ReadAt(chunk, offset)
		idx := bytes.Index(chunk[:n], search)
		if idx >= 0 {
			_, err = fd.Seek(offset+int64(idx)+2, 0)
			return err
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		offset += chunkSize
	}
	return fmt.Errorf("not found")
}

func (c *Reader) Read(b []byte) (int, error) {
	if c.eof {
		return 0, io.EOF
	}
	n, err := c.r.Read(b)
	//log.Println(string(b[:n]))
	if !c.summary {
		return n, err
	}
	//0x0d 0x0a 0x0d 0x0a
	i := bytes.Index(b, []byte("\r\n\r\n"))
	if i == -1 {
		return n, err
	}
	// for j := 0; j >= i; j++ {
	// 	b[j] = 0
	// }

	for j := i; j < len(b); j++ {
		b[j] = 0
	}

	c.eof = true
	return i, err
}
