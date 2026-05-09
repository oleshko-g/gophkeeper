package transform

import (
	"bytes"
	"io"
	"iter"
)

// StringFromReader returns a single use iterator over an [io.Reader].
// If the reader can be read into memory, it yields the contents as a single string.
// If the iterator encounters an error, it yields the contents read so far.
func StringFromReader(r io.Reader) iter.Seq[string] {
	return func(yield func(string) bool) {
		b := bytes.Buffer{}
		n, err := b.ReadFrom(r)
		if err == nil {
			yield(b.String())
			return
		}
		yield(b.String()[:n])
	}
}

func StringPtrFromBytes(b []byte) *string {
	s := string(b)
	return &s
}
