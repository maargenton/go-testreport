package json_test

import (
	"io"
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"
	"gopkg.in/yaml.v3"

	"github.com/maargenton/go-testreport/pkg/json"
	"github.com/maargenton/go-testreport/pkg/model"
	"github.com/maargenton/go-testreport/pkg/template"
)

func TestLoad(t *testing.T) {
	bdd.Given(t, "a file containing json formated test output", func(t *bdd.T) {
		filename := "./testdata/test-output.json"

		t.When("doing something", func(t *bdd.T) {
			var pkgs []model.Package
			err := fileutils.ReadFile(filename, func(r io.Reader) (err error) {
				pkgs, err = json.Load(r)
				return err
			})

			fileutils.WriteFile("./testdata/test-output.yaml", func(w io.Writer) error {
				e := yaml.NewEncoder(w)
				return e.Encode(pkgs)
			})

			tmpl := template.New("report")
			tmpl.Parse(template.MarkdownTemplate)

			require.That(t, err).IsError(nil)
			require.That(t, tmpl).IsNotNil()

			err = fileutils.WriteFile("./testdata/test-output.md", func(w io.Writer) error {
				return tmpl.Execute(w, pkgs)
			})
			verify.That(t, err).IsError(nil)

			t.Then("something happens", func(t *bdd.T) {
				verify.That(t, err).IsError(nil)
				verify.That(t, pkgs).Length().Eq(6)
			})
		})
	})
}
