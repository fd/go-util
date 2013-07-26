package errors

import (
	"testing"
)

func TestFwd(t *testing.T) {

	err := New("foo bar")

	if Fwd(err, "hello world").Error() != "hello world (foo bar)" {
		t.Fatal("bad")
	}

	if Fwd(err, "hello %s", "world").Error() != "hello world (foo bar)" {
		t.Fatal("bad")
	}

}
