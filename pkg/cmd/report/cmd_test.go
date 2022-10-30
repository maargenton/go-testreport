package report_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-testreport/pkg/cmd/report"
)

func TestCmdArgumentsValidation(t *testing.T) {
	bdd.Given(t, "a report command", func(t *bdd.T) {
		var cmd = report.Cmd{}

		t.When("called with no input", func(t *bdd.T) {
			var err = cmd.Run()

			t.Then("it returns an error", func(t *bdd.T) {
				require.That(t, err).ToString().Contains("no input specified")
			})
		})

		t.When("called with no output", func(t *bdd.T) {
			cmd.Inputs = []string{"testdata/sample-failure.yaml"}
			var err = cmd.Run()

			t.Then("it returns an error", func(t *bdd.T) {
				require.That(t, err).ToString().Contains("no output specified")
			})
		})

		t.When("called with input and output", func(t *bdd.T) {
			var dir = t.TempDir()
			var output = filepath.Join(dir, "summary.md")
			cmd.Inputs = []string{"testdata/sample-failure.yaml"}
			cmd.Outputs = []string{fmt.Sprintf("markdown-summary=%v", output)}
			var err = cmd.Run()

			t.Then("it reports test failure", func(t *bdd.T) {
				require.That(t, err).IsError(report.ErrTestFailure)
			})
		})
	})
}

func TestCmdInputType(t *testing.T) {
	bdd.Given(t, "a report command with setup output", func(t *bdd.T) {
		var cmd = report.Cmd{}
		var dir = t.TempDir()
		var output = filepath.Join(dir, "summary.md")
		cmd.Outputs = []string{fmt.Sprintf("markdown-summary=%v", output)}

		t.When("called with no input", func(t *bdd.T) {
			var err = cmd.Run()

			t.Then("it returns an error", func(t *bdd.T) {
				require.That(t, err).ToString().Contains("no input specified")
			})
		})

		t.When("called with valid yaml input", func(t *bdd.T) {
			cmd.Inputs = []string{"testdata/sample-failure.yaml"}
			var err = cmd.Run()

			t.Then("it succeeds but reports test failures", func(t *bdd.T) {
				require.That(t, err).IsError(report.ErrTestFailure)
			})
		})

		t.When("called with json input", func(t *bdd.T) {
			cmd.Inputs = []string{"testdata/sample-output.json"}
			var err = cmd.Run()

			t.Then("it succeeds but reports test failures", func(t *bdd.T) {
				require.That(t, err).IsError(report.ErrTestFailure)
			})
		})

		t.When("called with package name", func(t *bdd.T) {
			cmd.Inputs = []string{"../../sample"}
			cmd.Race = true
			var err = cmd.Run()

			t.Then("it succeeds with no error", func(t *bdd.T) {
				require.That(t, err).IsError(nil)
			})
		})
	})
}
