package model_test

import (
	"strings"
	"testing"
	"time"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-testreport/pkg/model"
)

func TestLoad(t *testing.T) {
	bdd.Given(t, "a file containing json formated test output", func(t *bdd.T) {
		filename := "./testdata/sample-output.json"

		t.When("calling LoadFromGoTestJsonFile()", func(t *bdd.T) {
			pkgs, err := model.LoadFromGoTestJsonFile(filename)
			require.That(t, err).IsError(nil)
			require.That(t, pkgs).Length().Eq(1)
			pkg := pkgs[0]

			t.Then("all leaf tests are extracted", func(t *bdd.T) {
				leaves := pkg.LeafTests()
				require.That(t, leaves).Field("Name").Eq([]string{
					"GET /api/v1/foo", "GET /api/v1/baz", "GET /api/v1/bar",
					"GET /api/v1/foo", "GET /api/v1/baz", "GET /api/v1/bar",
				})
			})

			t.Then("failures are captured", func(t *bdd.T) {
				leaves := pkg.LeafTests()
				require.That(t, leaves).Field("Failure").Eq([]bool{
					false, false, false,
					true, true, true,
				})
			})

			t.Then("coverage is extracted", func(t *bdd.T) {
				require.That(t, pkg.Coverage).Eq(50.0)
			})
			t.Then("elapsed time is extracted", func(t *bdd.T) {
				require.That(t, pkg.Elapsed).Eq(131 * time.Millisecond)
			})
		})
	})

	bdd.Given(t, "a file with skipped package", func(t *bdd.T) {
		filename := "./testdata/sample-output-skipped.json"

		t.When("calling LoadFromGoTestJsonFile()", func(t *bdd.T) {
			pkgs, err := model.LoadFromGoTestJsonFile(filename)
			require.That(t, err).IsError(nil)
			require.That(t, pkgs).Length().Eq(1)
			pkg := pkgs[0]

			t.Then("package is marked as skipped", func(t *bdd.T) {
				require.That(t, pkg.Skipped).IsTrue()

			})
		})
	})
}
func TestLoadError(t *testing.T) {
	bdd.Given(t, "given an empty stream", func(t *bdd.T) {
		var r = strings.NewReader("   ")

		t.When("calling LoadFromGoTestJson()", func(t *bdd.T) {
			pkgs, err := model.LoadFromGoTestJson(r)

			t.Then("an error is returned with no packages", func(t *bdd.T) {
				require.That(t, err).IsNotNil()
				require.That(t, pkgs).IsNil()

			})
		})
	})
}
