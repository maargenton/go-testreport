package gotest_test

import (
	"strings"
	"testing"
	"time"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/subexpr"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-testreport/pkg/gotest"
	"github.com/maargenton/go-testreport/pkg/model"
)

func TestLoad(t *testing.T) {
	bdd.Given(t, "a file containing json formatted test output", func(t *bdd.T) {
		var filename = "./testdata/sample/testfailure.json"

		t.When("calling ParseFile()", func(t *bdd.T) {
			results, err := gotest.ParseFile(filename)
			require.That(t, err).IsError(nil)
			require.That(t, results).IsNotNil()
			require.That(t, results.Packages).Length().Eq(1)
			pkg := results.Packages[0]

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
				require.That(t, pkg.Coverage).Eq(100.0)
			})
			t.Then("elapsed time is extracted", func(t *bdd.T) {
				require.That(t, pkg.Elapsed).Eq(323 * time.Millisecond)
			})
		})
	})

	bdd.Given(t, "a file with skipped package", func(t *bdd.T) {
		var filename = "./testdata/sample/testsuccess.json"

		t.When("calling ParseFile()", func(t *bdd.T) {
			results, err := gotest.ParseFile(filename)
			require.That(t, err).IsError(nil)
			require.That(t, results).IsNotNil()
			require.That(t, results.Packages).Length().Eq(1)
			pkg := results.Packages[0]

			t.Then("package is marked as skipped", func(t *bdd.T) {
				require.That(t, pkg.Skipped).IsTrue()

			})
		})
	})
}
func TestLoadError(t *testing.T) {
	bdd.Given(t, "given an empty stream", func(t *bdd.T) {
		var r = strings.NewReader("   ")

		t.When("calling Parse()", func(t *bdd.T) {
			results, err := gotest.Parse(r)

			t.Then("an error is returned with no packages", func(t *bdd.T) {
				require.That(t, err).IsNotNil()
				require.That(t, results).IsNil()

			})
		})
	})
}

func TestBuildError(t *testing.T) {
	bdd.Given(t, "a file containing build failures in json output", func(t *bdd.T) {
		var filename = "./testdata/sample/builderror.json"

		t.When("calling ParseFile()", func(t *bdd.T) {
			results, err := gotest.ParseFile(filename)
			require.That(t, err).IsError(nil)
			require.That(t, results).IsNotNil()
			require.That(t, results.Packages).Length().Eq(1)

			var failedPkg []*model.Package
			for _, p := range results.Packages {
				if p.PackageError != "" {
					failedPkg = append(failedPkg, p)
				}
			}

			t.Then("one package has build errors", func(t *bdd.T) {
				verify.That(t, failedPkg).Length().Eq(1)
				verify.That(t, failedPkg[0].Name).Eq("github.com/maargenton/go-testreport/pkg/sample/builderror")
			})

			t.Then("packages with build errors are marked with at least 1 failure", func(t *bdd.T) {
				verify.That(t, failedPkg).All(
					subexpr.Value().Field("Failed").Ge(1),
				)
			})
		})
	})
}
