package gotest_test

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"
	"github.com/maargenton/go-testreport/pkg/gotest"
)

func TestRun(t *testing.T) {

	bdd.Given(t, "a package reference", func(t *bdd.T) {
		var pkgName = "../model"

		t.When("calling Run()", func(t *bdd.T) {
			results, err := gotest.Run(pkgName, gotest.Race())

			t.Then("it runs go test and captures the results", func(t *bdd.T) {
				require.That(t, err).IsError(nil)
				require.That(t, results).IsNotNil()
				verify.That(t, results.Passed).Eq(8)
				verify.That(t, results.Failed).Eq(0)
				verify.That(t, results.Success).IsTrue()

				require.That(t, results.Packages).Length().Eq(1)
				var pkg = results.Packages[0]
				verify.That(t, pkg.LeafTests()).Length().Eq(8)
			})
		})
	})

	// An invalid package reference as a `go test` argument is no longer
	// reported as an output with no package reference; it instead report a
	// failure on the bogus package name. Since the report is focused on tests
	// success / failure, this is no longer caught, except there will be no test
	// report.
	//
	// TODO: find a way to report that properly

	// bdd.Given(t, "an invalid package reference", func(t *bdd.T) {
	// 	var pkgName = "../bad_sample"

	// 	t.When("calling Run()", func(t *bdd.T) {
	// 		results, err := gotest.Run(pkgName, gotest.Race())

	// 		t.Then("it runs go test and captures the results", func(t *bdd.T) {
	// 			verify.That(t, results).IsNil()
	// 			verify.That(t, err).IsNotNil()
	// 			verify.That(t, err).ToString().Contains("unresolved package reference")
	// 		})
	// 	})
	// })
}
