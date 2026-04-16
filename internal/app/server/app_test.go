package server_test

import (
	"testing"

	"github.com/oleshko-g/gophkeeper/internal/app/server"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestApp(t *testing.T) {
	app := server.App{}
	t.Run("Configure", func(t *testing.T) {
		err := app.Configure()
		if err != nil {
			t.Fatal(err)
		}
	})
}
