package uuidv7_test

import (
	"testing"
	"time"

	uuidv7 "github.com/oleshko-g/gophkeeper/internal/uuid-v7"
)

func TestTime(t *testing.T) {
	now := time.Now()
	time.Sleep(100 * time.Millisecond)
	id := uuidv7.New()
	if !id.Time().After(now) {
		t.Errorf("expected: id.Time() %q to be before after time.Now() + 100ms %v", now, id.Time())
	}
}
