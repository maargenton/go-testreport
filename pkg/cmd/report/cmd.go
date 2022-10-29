package report

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/maargenton/go-errors"
	"github.com/maargenton/go-fileutils"

	"github.com/maargenton/go-testreport/pkg/gotest"
	"github.com/maargenton/go-testreport/pkg/model"
	"github.com/maargenton/go-testreport/pkg/template"
)

// ErrTestFailure is returned by the command when the command executes
// successfully and generates the requested reports, but some of the inputs
// contain test failures.
const ErrTestFailure = errors.Sentinel("ErrTestFailure")

type Cmd struct {
	Inputs  []string `opts:"args, name:input" desc:"package or packages to run tests from, or filename containing test results"`
	Outputs []string `opts:"-o, --output"     desc:"one of more output to generate, formatted as <template>=<output-filename>.\ntemplate can be either 'yaml', 'markdown' or a the name of a file containing a custom template"`

	Race        bool   `opts:"--race"                                  desc:"run the tests with race detector on"`
	ShiftHeader int    `opts:"--md-shift-headers, default:0"           desc:"shift the level of markdown headers"`
	Title       string `opts:"--md-title,         default:Test report" desc:"shift the level of markdown headers"`
}

func (cmd *Cmd) Run() error {
	if len(cmd.Inputs) == 0 {
		return fmt.Errorf("nothing to do; no input specified")
	}

	if len(cmd.Outputs) == 0 {
		return fmt.Errorf("nothing to do; no output specified")
	}

	var results = &model.Results{}
	for _, input := range cmd.Inputs {
		content, err := cmd.loadInput(input)
		if err != nil {
			return err
		}
		results.Packages = append(results.Packages, content.Packages...)
	}
	results.UpdateCounts()

	for _, output := range cmd.Outputs {
		err := cmd.saveOutput(output, results)
		if err != nil {
			return err
		}
	}

	if !results.Success {
		if results.Failed > 1 {
			return ErrTestFailure.Errorf("%v tests failed", results.Failed)
		} else {
			return ErrTestFailure.Errorf("%v test failed", results.Failed)
		}
	}
	return nil
}

func (cmd *Cmd) loadInput(input string) (results *model.Results, err error) {
	ext := filepath.Ext(input)
	if ext == ".yaml" {
		return model.LoadFromYAMLFile(input)
	}
	if ext == ".json" {
		return gotest.ParseFile(input)
	}

	var opts []gotest.RunOpts
	if cmd.Race {
		opts = append(opts, gotest.Race())
	}
	return gotest.Run(input, opts...)
}

func (cmd *Cmd) saveOutput(output string, results *model.Results) error {
	parts := strings.SplitN(output, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid output specifier: '%v'; expected <type>=<filename>", output)
	}
	if parts[0] == "yaml" {
		if parts[1] == "-" {
			return model.SaveToYAML(os.Stdout, results)
		}
		return model.SaveToYAMLFile(parts[1], results)
	}

	var tmpl *template.Template
	var values = map[string]interface{}{
		"Title":       cmd.Title,
		"Results":     results,
		"Packages":    results.Packages,
		"HeaderShift": cmd.ShiftHeader,
	}

	if builtin, ok := template.Builtin[parts[0]]; ok {
		tmpl = template.New("report", values)
		_, err := tmpl.Parse(builtin)
		if err != nil {
			return err
		}
	} else {
		tmpl = template.New(parts[0], values)
		_, err := tmpl.ParseFiles(parts[0])
		if err != nil {
			return err
		}
	}

	if parts[1] == "-" {
		return tmpl.Execute(os.Stdout, results.Packages)
	}
	return fileutils.WriteFile(parts[1], func(w io.Writer) error {
		return tmpl.Execute(w, results.Packages)
	})
}
