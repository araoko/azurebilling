package csp

import (
	"bytes"
	"io"
	"os"
)

type Reader struct {
	r       io.Reader
	summary bool
	eof     bool
}

func NewCspReader(fd *os.File, summary bool) *Reader {
	skipBOMSummaryHeading(fd)
	if !summary {
		err := skipSummary(fd)
		if err != nil {
			panic(err)
		}
	}
	return &Reader{r: fd, summary: summary}
}

func skipBOMSummaryHeading(fd *os.File) {
	//Byte Order Mark (BOM) and summary headiing is 14 bytes
	var discard [14]byte
	_, err := io.ReadFull(fd, discard[:])
	if err != nil {
		panic(err)
	}
}

func skipSummary(fd *os.File) error {
	var offset int64
	var retErr error
	//finfo, _ := fd.Stat()
	//log.Println("skipping Summary of file:", finfo.Name(), "with mode", finfo.Mode().String())
	discard := make([]byte, 1024)
	for {
		n, err := fd.Read(discard)
		if n == 0 {
			retErr = err
			//log.Println("n = 0 with eror", err)
			break
		}
		h := []byte("Daily Usage\"")
		i := bytes.Index(discard, h)
		if i != -1 {
			//log.Println("i = ", i)
			//log.Println(string(discard[i : i+len(h)]))
			offset += int64(i + len(h) + 14)
			break
		}
		offset += int64(n)
		//log.Println("i = -1")

	}
	if retErr != nil {
		return retErr
	}
	_, err := fd.Seek(offset, 0)
	//log.Println("calc offset", offset, "real offset", off)
	return err
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
