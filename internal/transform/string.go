package transform

import (
	"bytes"
	"io"
	"iter"
)

// StringFromReader returns a single use iterator over an [io.Reader].
// If the reader can be read into memory, it yields the contents as a single string.
// If the reader cannot be read into memory, it yields an empty string.
func StringFromReader(r io.Reader) iter.Seq[string] {
	var b bytes.Buffer
	return func(yield func(string) bool) {
		_, err := b.ReadFrom(r)
		if err == nil {
			yield(b.String())
			return
		}
		yield("")
	}
}
