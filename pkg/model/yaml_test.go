package model_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"
	"github.com/maargenton/go-testreport/pkg/model"
)

func TestLoadFromYAMLFile(t *testing.T) {
	bdd.Given(t, "a YAML file containing tests", func(t *bdd.T) {
		filename := "./testdata/sample-failure.yaml"

		t.When("calling LoadFromYAMLFile()", func(t *bdd.T) {
			pkgs, err := model.LoadFromYAMLFile(filename)
			require.That(t, err).IsError(nil)
			require.That(t, pkgs).Length().Eq(1)

			t.Then("package details are loaded correctly", func(t *bdd.T) {
				verify.That(t, pkgs[0].Coverage).Eq(42)
				verify.That(t, pkgs[0].Elapsed).Eq(258 * time.Millisecond)
			})
			t.Then("nested tests are linked through Parent field", func(t *bdd.T) {
				tests := pkgs[0].LeafTests()
				require.That(t, tests).Length().Eq(2)
				verify.That(t, tests[0].Parent).IsNotNil()
				verify.That(t, tests[1].Parent).IsNotNil()
			})
		})
	})

	bdd.Given(t, "an invalid YAML", func(t *bdd.T) {
		filename := "./testdata/sample-invalid.yaml"

		t.When("calling LoadFromYAMLFile()", func(t *bdd.T) {
			_, err := model.LoadFromYAMLFile(filename)

			t.Then("an error is returned", func(t *bdd.T) {
				require.That(t, err).IsNotNil()
			})
		})
	})
}

func TestSaveToYAMLFile(t *testing.T) {
	bdd.Given(t, "a valid list of packages", func(t *bdd.T) {
		var pkgs = []model.Package{
			{Name: "foo", Coverage: 42},
			{Name: "bar", Coverage: 43},
		}
		var tempDir = t.TempDir()
		var filename = filepath.Join(tempDir, "tests.yaml")

		t.When("when calling SaveToYAMLFile()", func(t *bdd.T) {
			err := model.SaveToYAMLFile(filename, pkgs)
			require.That(t, err).IsError(nil)

			t.Then("content can be reloaded with LoadFromYAMLFile()", func(t *bdd.T) {
				loaded, err := model.LoadFromYAMLFile(filename)
				require.That(t, err).IsError(nil)
				verify.That(t, loaded).Length().Eq(2)
				verify.That(t, loaded).Field("Name").Eq([]string{"foo", "bar"})
				verify.That(t, loaded).Field("Coverage").Eq([]float64{42, 43})
			})
		})
	})
}
