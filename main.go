package main

import (
	"github.com/maargenton/go-cli"
	"github.com/maargenton/go-testreport/pkg/cmd/report"
)

func main() {
	cli.Run(&cli.Command{
		Handler:     &report.Cmd{},
		Description: "Generate test report from 'go test' output",
	})
}
