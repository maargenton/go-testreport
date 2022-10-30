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
			results, err := model.LoadFromYAMLFile(filename)
			require.That(t, err).IsError(nil)
			require.That(t, results).IsNotNil()
			require.That(t, results.Packages).Length().Eq(1)
			pkg := results.Packages[0]

			t.Then("package details are loaded correctly", func(t *bdd.T) {
				verify.That(t, pkg.Coverage).Eq(42)
				verify.That(t, pkg.Elapsed).Eq(258 * time.Millisecond)
			})
			t.Then("nested t are linked through Parent field", func(t *bdd.T) {
				tests := pkg.LeafTests()
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

		var results = &model.Results{
			Packages: []*model.Package{
				{Name: "foo", Coverage: 42},
				{Name: "bar", Coverage: 43},
			},
		}
		var tempDir = t.TempDir()
		var filename = filepath.Join(tempDir, "tests.yaml")

		t.When("when calling SaveToYAMLFile()", func(t *bdd.T) {
			err := model.SaveToYAMLFile(filename, results)
			require.That(t, err).IsError(nil)

			t.Then("content can be reloaded with LoadFromYAMLFile()", func(t *bdd.T) {
				loaded, err := model.LoadFromYAMLFile(filename)
				require.That(t, err).IsError(nil)
				require.That(t, loaded).IsNotNil()
				verify.That(t, loaded.Packages).Length().Eq(2)
				verify.That(t, loaded.Packages).Field("Name").Eq([]string{"foo", "bar"})
				verify.That(t, loaded.Packages).Field("Coverage").Eq([]float64{42, 43})
			})
		})
	})
}
