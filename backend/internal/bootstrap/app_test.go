package bootstrap

import (
	"context"
	"errors"
	"testing"
)

func TestAppShutdownClosesResourcesInReverseOrder(t *testing.T) {
	var closed []string
	app := &App{
		resources: []resource{
			shutdownResource{name: "database", closed: &closed},
			shutdownResource{name: "cache", closed: &closed},
			shutdownResource{name: "queue", closed: &closed},
		},
	}

	if err := app.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	want := []string{"queue", "cache", "database"}
	for i := range want {
		if closed[i] != want[i] {
			t.Fatalf("closed[%d] = %q, want %q", i, closed[i], want[i])
		}
	}
}

func TestAppShutdownJoinsResourceErrors(t *testing.T) {
	errDatabase := errors.New("database close failed")
	errCache := errors.New("cache close failed")
	app := &App{
		resources: []resource{
			shutdownResource{err: errDatabase},
			shutdownResource{err: errCache},
		},
	}

	err := app.Shutdown(context.Background())
	if !errors.Is(err, errDatabase) {
		t.Fatalf("Shutdown() error does not include database error: %v", err)
	}
	if !errors.Is(err, errCache) {
		t.Fatalf("Shutdown() error does not include cache error: %v", err)
	}
}

type shutdownResource struct {
	name   string
	err    error
	closed *[]string
}

func (r shutdownResource) Shutdown(ctx context.Context) error {
	if r.closed != nil {
		*r.closed = append(*r.closed, r.name)
	}
	return r.err
}
