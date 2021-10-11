package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/maargenton/go-cli"
	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-fileutils/pkg/popen"
	"github.com/maargenton/go-testreport/pkg/model"
	"github.com/maargenton/go-testreport/pkg/template"
)

func main() {
	cli.Run(&cli.Command{
		Handler:     &reportCmd{},
		Description: "Generate test report from 'go test' output",
	})
}

type reportCmd struct {
	Inputs  []string `opts:"args, name:input" desc:"package or packages to run tests from, or filename containign test results"`
	Race    bool     `opts:"--race"             desc:"run the tests with race detector on"`
	Outputs []string `opts:"-o, --output"       desc:"one of more output to generate, formated as <template>=<output-filename>.\ntemplate can be either 'yaml', 'markdown' or a the name of a file containing a custom template"`
}

func (cmd *reportCmd) Run() error {

	if len(cmd.Inputs) == 0 {
		return fmt.Errorf("nothing to do; no inputs specified")
	}

	if len(cmd.Outputs) == 0 {
		return fmt.Errorf("nothing to do; no outputs specified")
	}

	var pkgs []model.Package
	for _, input := range cmd.Inputs {
		content, err := cmd.loadInput(input)
		if err != nil {
			return err
		}
		pkgs = append(pkgs, content...)
	}

	for _, output := range cmd.Outputs {
		err := cmd.saveOutput(output, pkgs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *reportCmd) loadInput(input string) (pkgs []model.Package, err error) {
	ext := filepath.Ext(input)
	if ext == ".yaml" {
		return model.LoadFromYAMLFile(input)
	}
	if ext == ".json" {
		return model.LoadFromGoTestJsonFile(input)
	}

	testCmd := popen.Command{
		Command: "go",
		Arguments: []string{
			"test",
			"-cover",
			"-json",
		},
	}
	_ = testCmd
	if cmd.Race {
		testCmd.Arguments = append(testCmd.Arguments, "-race")
	}
	testCmd.Arguments = append(testCmd.Arguments, input)
	testCmd.StdoutReader = func(r io.Reader) error {
		pkgs, err = model.LoadFromGoTestJson(r)
		return err
	}
	_, _, err = testCmd.Run(context.Background())
	return
}

func (cmd *reportCmd) saveOutput(output string, pkgs []model.Package) error {
	parts := strings.SplitN(output, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid output specifier: '%v'; expected <type>=<filename>", output)
	}
	if parts[0] == "yaml" {
		if parts[1] == "-" {
			return model.SaveToYAML(os.Stdout, pkgs)
		} else {
			return model.SaveToYAMLFile(parts[1], pkgs)
		}
	}

	var tmpl *template.Template

	if parts[0] == "md" || parts[0] == "markdown" {
		tmpl = template.New("report")
		_, err := tmpl.Parse(template.MarkdownTemplate)
		if err != nil {
			return err
		}
	} else {
		tmpl = template.New(parts[0])
		_, err := tmpl.ParseFiles(parts[0])
		if err != nil {
			return err
		}
	}

	if parts[1] == "-" {
		return tmpl.Execute(os.Stdout, pkgs)
	} else {
		return fileutils.WriteFile(parts[1], func(w io.Writer) error {
			return tmpl.Execute(w, pkgs)
		})
	}
}
