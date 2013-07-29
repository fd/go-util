package errors

import (
	"encoding/json"
	"fmt"
)

func ExampleAnnotate() {
	err := New("%s err", "foo")
	err.AddContext("hello=%s", "world")

	err = Annotate(err, "%s err", "bar")
	err.AddContext("c=%d", 7)
	err.AddContext("a=%d", 42)

	fmt.Println(err)

	// Output:
	// error: bar err
	//   context:
	//     a = 42
	//     c = 7
	//   location:
	//     /Users/fd/src/go/src/github.com/fd/go-util/errors/annotation_test.go:11 (0x33cd6)
	//       com/fd/go-util/errors.ExampleAnnotate: err = Annotate(err, "%s err", "bar")
	//     /opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/example.go:98 (0x2d1e5)
	//       runExample: eg.F()
	//     /opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/example.go:36 (0x2ce27)
	//       RunExamples: if !runExample(eg) {
	//   error: foo err
	//     context:
	//       hello = world
	//     location:
	//       /Users/fd/src/go/src/github.com/fd/go-util/errors/annotation_test.go:8 (0x33b2b)
	//         com/fd/go-util/errors.ExampleAnnotate: err := New("%s err", "foo")
	//       /opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/example.go:98 (0x2d1e5)
	//         runExample: eg.F()
	//       /opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/example.go:36 (0x2ce27)
	//         RunExamples: if !runExample(eg) {
}

func ExampleAnnotateJSON() {
	err := New("%s err", "foo")
	err.AddContext("hello=%s", "world")

	err = Annotate(err, "%s err", "bar")
	err.AddContext("c=%d", 7)
	err.AddContext("a=%d", 42)

	data, _ := json.MarshalIndent(err, "", "  ")
	fmt.Println(string(data))

	// Output:
	// {
	//   "message": "bar err",
	//   "context": {
	//     "a": "42",
	//     "c": "7"
	//   },
	//   "stack": "/Users/fd/src/go/src/github.com/fd/go-util/errors/error_test.go:46 (0x34445)\n  com/fd/go-util/errors.ExampleAnnotateJSON: err = Annotate(err, \"%s err\", \"bar\")\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/example.go:98 (0x2d1d5)\n  runExample: eg.F()\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/example.go:36 (0x2ce17)\n  RunExamples: if !runExample(eg) {\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/testing.go:366 (0x2df5c)\n  Main: exampleOk := RunExamples(matchString, examples)\ngithub.com/fd/go-util/errors/_test/_testmain.go:45 (0x21ca)\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/runtime/proc.c:182 (0x15592)\n  main: main·main();\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/runtime/proc.c:1223 (0x17420)\n  goexit: runtime·goexit(void)\n",
	//   "error": {
	//     "message": "foo err",
	//     "context": {
	//       "hello": "world"
	//     },
	//     "stack": "/Users/fd/src/go/src/github.com/fd/go-util/errors/error_test.go:43 (0x342a0)\n  com/fd/go-util/errors.ExampleAnnotateJSON: err := New(\"%s err\", \"foo\")\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/example.go:98 (0x2d1d5)\n  runExample: eg.F()\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/example.go:36 (0x2ce17)\n  RunExamples: if !runExample(eg) {\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/testing/testing.go:366 (0x2df5c)\n  Main: exampleOk := RunExamples(matchString, examples)\ngithub.com/fd/go-util/errors/_test/_testmain.go:45 (0x21ca)\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/runtime/proc.c:182 (0x15592)\n  main: main·main();\n/opt/boxen/homebrew/Cellar/golang/1.1-boxen1/src/pkg/runtime/proc.c:1223 (0x17420)\n  goexit: runtime·goexit(void)\n"
	//   }
	// }
}
