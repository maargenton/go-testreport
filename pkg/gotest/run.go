package gotest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"

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

	// Catch exec.ExitError and ignore the error. `go test` can exit with status
	// 1 on test failure and still produce complete meaningful output.
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		if exitError.ProcessState.ExitCode() == 1 {
			err = nil
		}
	}

	// But generate an error if the output of `go test` did not report any
	// package.
	if len(results.Packages) == 0 {
		return nil, fmt.Errorf("unresolved package reference '%v'", input)
	}
	return
}
