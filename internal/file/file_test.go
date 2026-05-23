package file_test

import (
	"os"
	"path"
	"testing"

	"github.com/oleshko-g/gophkeeper/internal/file"
)

func TestRWOpenCreatePrivate(t *testing.T) {
	tmp := t.TempDir()

	t.Run("empty file name", func(t *testing.T) {
		t.Parallel()
		_, err := file.RWOpenCreatePrivate(path.Join(tmp, ""))
		if err == nil {
			t.Errorf("expected: an error. got: nil")
		}
	})

	t.Run("non-empty file name, successful close", func(t *testing.T) {
		t.Parallel()
		f, err := file.RWOpenCreatePrivate(path.Join(tmp, "non-empty"))
		if err != nil {
			t.Errorf("expected: no error. got: %v", err)
		}
		defer func() {
			err := f.Close()
			if err != nil {
				t.Error(err)
			}
		}()
	})

	t.Run("non-existent directory", func(t *testing.T) {
		t.Parallel()
		_, err := file.RWOpenCreatePrivate(path.Join(tmp, "non-existent", "dir"))
		if err == nil {
			t.Errorf("expected: an error. got: nil")
		}
	})

	t.Run("RW, successful close", func(t *testing.T) {
		t.Parallel()
		f, err := file.RWOpenCreatePrivate(path.Join(tmp, "private-rw"))
		if err != nil {
			t.Errorf("expected: no error. got: %v", err)
		}
		defer func() {
			err := f.Close()
			if err != nil {
				t.Error(err)
			}
		}()

		s, err := f.Stat()
		if err != nil {
			t.Errorf("expected: no error. got: %v", err)
		}

		var expected os.FileMode = 0o600
		if mode := s.Mode(); mode != expected.Perm() {
			t.Errorf("expected: os.FileMode == %v. got: %v", expected, mode)
		}
	})
}
