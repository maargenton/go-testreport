package report

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
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
	Inputs    []string `opts:"args, name:input" desc:"package or packages to run tests from, or filename containing test results"`
	Templates []string `opts:"-t, --template"   desc:"one or more template files to load as part of the available templates"`
	Outputs   []string `opts:"-o, --output"     desc:"one of more output to generate, formatted as <template>=<output-filename>.\ntemplate can be either 'yaml' or the name of a builtin or custom template"`

	Race        bool   `opts:"--race"                                  desc:"run the tests with race detector on"`
	ShiftHeader int    `opts:"--md-shift-headers, default:0"           desc:"shift the level of markdown headers"`
	Title       string `opts:"--md-title,         default:Test report" desc:"shift the level of markdown headers"`

	tmpl *template.Template // Parsed template from builtin and custom template files
}

func (cmd *Cmd) Version() string {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		return buildInfo.Main.Version
	}
	return "v0.0.0-unknown"
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

	if err := cmd.loadTemplates(results); err != nil {
		return err
	}

	for _, output := range cmd.Outputs {
		parts := strings.SplitN(output, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid output specifier: '%v'; expected <type>=<filename>", output)
		}
		var format, outputFilename = parts[0], parts[1]

		if format == "yaml" {
			err := cmd.saveYAMLOutput(outputFilename, results)
			if err != nil {
				return err
			}
		} else {
			err := cmd.saveTemplateOutput(outputFilename, format, results)
			if err != nil {
				return err
			}
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

func (cmd *Cmd) loadTemplates(results *model.Results) error {

	var values = map[string]interface{}{
		"Title":       cmd.Title,
		"Results":     results,
		"Packages":    results.Packages,
		"HeaderShift": cmd.ShiftHeader,
	}

	cmd.tmpl = template.New("default", values)
	if len(cmd.Templates) != 0 {
		if _, err := cmd.tmpl.ParseFiles(cmd.Templates...); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *Cmd) saveYAMLOutput(output string, results *model.Results) error {
	if output == "-" {
		return model.SaveToYAML(os.Stdout, results)
	}
	return model.SaveToYAMLFile(output, results)
}

func (cmd *Cmd) saveTemplateOutput(output string, format string, results *model.Results) error {
	if output == "-" {
		return cmd.tmpl.ExecuteTemplate(os.Stdout, format, results.Packages)
	} else {
		return fileutils.WriteFile(output, func(w io.Writer) error {
			return cmd.tmpl.ExecuteTemplate(w, format, results.Packages)
		})
	}
}
