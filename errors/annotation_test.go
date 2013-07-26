package errors

import (
	"fmt"
)

func ExampleAnnotate() {
	err := New("%s err", "foo")
	err.AddContext("hello=%s", "world")

	err = Annotate(err, "%s err", "bar")
	err.AddContext("c=%d", 7)
	err.AddContext("a=%d", 42)

	fmt.Println(err.Error())

	// Output:
	// hello
}
