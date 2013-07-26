package errors

import (
	"testing"
)

func TestFmt(t *testing.T) {

	if Fmt("hello world").Error() != "hello world" {
		t.Fail()
	}

	if Fmt("hello %s", "world").Error() != "hello world" {
		t.Fail()
	}

}
