package gotest

import (
	"context"
	"io"

	"github.com/maargenton/go-fileutils/pkg/popen"

	"github.com/maargenton/go-testreport/pkg/model"
)

type RunOpts func() []string

func Race() RunOpts {
	return func() []string {
		return []string{"-race"}
	}
}

// Run invokes `go test -json` on the specified package and parses the results
// like `Parse()`.
func Run(input string, opts ...RunOpts) (results *model.Results, err error) {
	testCmd := popen.Command{
		Command: "go",
		Arguments: []string{
			"test",
			"-cover",
			"-json",
		},
	}

	for _, opt := range opts {
		testCmd.Arguments = append(testCmd.Arguments, opt()...)
	}
	testCmd.Arguments = append(testCmd.Arguments, input)
	testCmd.StdoutReader = func(r io.Reader) error {
		results, err = Parse(r)
		return err
	}
	_, _, err = testCmd.Run(context.Background())
	return
}
