package assert

import (
	"fmt"
	"io"
	"io/fs"
	"testing"
)

func TestErrorContains(t *testing.T) {
	err := fmt.Errorf(
		"closed socket: %w %w",
		io.EOF,
		&fs.PathError{Op: "read", Path: "socket", Err: io.ErrClosedPipe},
	)
	ErrorContains(t, err, "closed socket")
	ErrorContains(t, err, io.EOF)
	ErrorContains(t, err, io.ErrClosedPipe)
	var pathError *fs.PathError
	ErrorContains(t, err, &pathError)
}
