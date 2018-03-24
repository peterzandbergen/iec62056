package telegram

// import (
// 	"bytes"
// 	"io"
// )

// type Reader struct {
// 	r   io.Reader
// 	buf bytes.Buffer
// }

// func NewReader(r io.Reader) *Reader {
// 	return &Reader{
// 		r: r,
// 	}
// }

// func (r *Reader) Read(p []byte) (n int, err error) {
// 	if r.buf.Len() > 0 {
// 		return r.buf.Read(p)
// 	}
// 	return r.r.Read(p)
// }

// func (r *Reader) UnreadByte(b byte) {
// 	r.buf.UnreadByte(b)
// }
