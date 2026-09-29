package progress

import (
	"fmt"
	"io"
)

type Reader struct {
	Reader  io.Reader
	Total   int64
	Current int64
}

func NewReader(reader io.Reader, total int64) *Reader {
	return &Reader{
		Reader:  reader,
		Total:   total,
		Current: 0,
	}
}

func (r *Reader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.Current += int64(n)
	percent := float64(r.Current) / float64(r.Total) * 100
	fmt.Println("Sending...", percent, "%")
	return n, err
}
