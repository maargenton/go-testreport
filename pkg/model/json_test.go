package model_test

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-testreport/pkg/model"
)

func TestLoad(t *testing.T) {
	bdd.Given(t, "a file containing json formated test output", func(t *bdd.T) {
		filename := "./testdata/test-output.json"

		t.When("calling LoadFromGoTestJsonFile()", func(t *bdd.T) {
			pkgs, err := model.LoadFromGoTestJsonFile(filename)

			t.Then("content is loaded correctly", func(t *bdd.T) {
				require.That(t, err).IsError(nil)
				require.That(t, pkgs).Length().Eq(6)

				verify.That(t, pkgs[1].Tests).Length().Eq(10)
			})

			// tmpl := template.New("report")
			// tmpl.Parse(template.MarkdownTemplate)

			// require.That(t, err).IsError(nil)
			// require.That(t, tmpl).IsNotNil()

			// err = fileutils.WriteFile("./testdata/test-output.md", func(w io.Writer) error {
			// 	return tmpl.Execute(w, pkgs)
			// })
			// verify.That(t, err).IsError(nil)

		})
	})
}
