package errtrace

import (
	"errors"
	"strings"
	"testing"
)

func ErrAtDepth(n int) error {
	if n == 0 {
		return Wrap(errors.New("test error"))
	}
	return Wrap(ErrAtDepth(n - 1))
}

func TestMessage(t *testing.T) {
	err := ErrAtDepth(1)

	msg := err.Error()

	if !strings.Contains(msg, "go-errtrace/err_test.go:") {
		t.Errorf("expected error text, got: %s", msg)
	}
}

func TestStackTrace(t *testing.T) {
	err := ErrAtDepth(3)

	var tracedErr *TracedError
	if !errors.As(err, &tracedErr) {
		t.Errorf("expected error to be of type *TracedError, got: %T", err)
	}

	for i := range 3 {
		file := tracedErr.stacktrace[i].file
		if !strings.Contains(file, "go-errtrace/err_test.go") {
			t.Errorf("expected stack trace to contain 'go-errtrace/err_test.go, got: %s'", file)
		}
	}
}
