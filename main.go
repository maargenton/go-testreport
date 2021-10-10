package main

import (
	"encoding/json"
	"fmt"

	"github.com/maargenton/go-cli"
)

func main() {
	cli.Run(&cli.Command{
		Handler:     &reportCmd{},
		Description: "Generate test report from 'go test' output",
	})
}

type reportCmd struct {
	Package []string `opts:"args, name:package" desc:"one or more package to tests on"`
	Input   string   `opts:"-i, --input"        desc:"filename for pre-generated json test outputs"`
	Race    bool     `opts:"--race"             desc:"run the tests with race detector on"`
	Outputs []string `opts:"-o, --output"       desc:"one of more output to generate, formated as <template>=<output-filename>.\ntemplate can be either 'yaml', 'markdown' or a the name of a file containing a custom template"`
}

func (options *reportCmd) Run() error {
	d, err := json.Marshal(options)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(d))
	return nil
}
