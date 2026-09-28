package gotest_test

import (
	"fmt"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"
	"github.com/maargenton/go-testreport/pkg/gotest"
)

func TestLoadJSONOutput(t *testing.T) {
	var tcs = []struct {
		filename            string
		expectedImportPaths []string
		expectedPackages    []string
		expectedTestNames   []string
	}{
		// Build errors
		{
			filename: "testdata/sample/builderror.json",
			expectedImportPaths: []string{
				"github.com/maargenton/go-testreport/pkg/sample/builderror",
			},
			expectedPackages: []string{
				"github.com/maargenton/go-testreport/pkg/sample/builderror",
			},
		},

		{
			filename: "testdata/sample/builderrorindirect.json",
			expectedImportPaths: []string{
				"github.com/maargenton/go-testreport/pkg/sample/builderrorindirect/subpkg",
			},
			expectedPackages: []string{
				"github.com/maargenton/go-testreport/pkg/sample/builderrorindirect",
				"github.com/maargenton/go-testreport/pkg/sample/builderrorindirect/subpkg",
			},
		},

		{
			filename: "testdata/sample/buildtesterror.json",
			expectedImportPaths: []string{
				"github.com/maargenton/go-testreport/pkg/sample/buildtesterror_test",
			},
			expectedPackages: []string{
				"github.com/maargenton/go-testreport/pkg/sample/buildtesterror",
			},
		},

		// Test records
		{
			filename: "testdata/sample/testfailure.json",
			expectedPackages: []string{
				"github.com/maargenton/go-testreport/pkg/sample/testfailure",
			},
			expectedTestNames: []string{
				"TestError",
				"TestError/TestFoo",
				"TestError/TestFoo/GET_/api/v1/foo",
				"TestError/TestFoo#01",
				"TestError/TestFoo#01/GET_/api/v1/baz",
				"TestError/TestFoo#02",
				"TestError/TestFoo#02/GET_/api/v1/bar",
				"TestError/TestBar",
				"TestError/TestBar/GET_/api/v1/foo",
				"TestError/TestBar#01",
				"TestError/TestBar#01/GET_/api/v1/baz",
				"TestError/TestBar#02",
				"TestError/TestBar#02/GET_/api/v1/bar",
			},
		},

		{
			filename: "testdata/sample/testpanic.json",
			expectedPackages: []string{
				"github.com/maargenton/go-testreport/pkg/sample/testpanic",
			},
			expectedTestNames: []string{
				"TestSuccess",
				"TestSuccess/TestFoo",
				"TestSuccess/TestFoo/GET_/api/v1/foo",
				"TestSuccess/TestFoo#01",
				"TestSuccess/TestFoo#01/GET_/api/v1/baz",
				"TestSuccess/TestFoo#02",
				"TestSuccess/TestFoo#02/GET_/api/v1/bar",
				"TestSuccess/TestBar",
				"TestSuccess/TestBar/GET_/api/v1/foo",
				"TestSuccess/TestBar#01",
				"TestSuccess/TestBar#01/GET_/api/v1/baz",
			},
		},

		{
			filename: "testdata/sample/testskip.json",
			expectedPackages: []string{
				"github.com/maargenton/go-testreport/pkg/sample/testskip",
			},
			expectedTestNames: []string{
				"TestError",
			},
		},

		{
			filename: "testdata/sample/testsuccess.json",
			expectedPackages: []string{
				"github.com/maargenton/go-testreport/pkg/sample/testsuccess",
			},
			expectedTestNames: []string{
				"TestSuccess",
				"TestSuccess/TestFoo",
				"TestSuccess/TestFoo/GET_/api/v1/foo",
				"TestSuccess/TestFoo#01",
				"TestSuccess/TestFoo#01/GET_/api/v1/baz",
				"TestSuccess/TestFoo#02",
				"TestSuccess/TestFoo#02/GET_/api/v1/bar",
				"TestSuccess/TestBar",
				"TestSuccess/TestBar/GET_/api/v1/foo",
				"TestSuccess/TestBar#01",
				"TestSuccess/TestBar#01/GET_/api/v1/baz",
				"TestSuccess/TestBar#02",
				"TestSuccess/TestBar#02/GET_/api/v1/bar",
			},
		},
	}
	for _, tc := range tcs {
		bdd.Given(t, "a file containing captured JSON test output", func(t *bdd.T) {
			t.When(fmt.Sprintf("calling LoadJSONOutput( %q )", tc.filename), func(t *bdd.T) {

				var output, err = gotest.LoadJSONOutput(tc.filename)
				require.That(t, err).IsError(nil)

				t.Then("it loads and groups the collected outputs correctly", func(t *bdd.T) {
					verify.That(t, output.ImportPaths).MapKeys().IsEqualSet(tc.expectedImportPaths)
					verify.That(t, output.PackageOrder).IsEqualSet(tc.expectedPackages)

					var testNames []string
					for _, pkgName := range output.PackageOrder {
						var pkg = output.Packages[pkgName]
						for _, test := range pkg.Tests {
							_ = test
							testNames = append(testNames, test.TestName)
						}
					}
					verify.That(t, testNames).Eq(tc.expectedTestNames)
				})
			})
		})
	}
}
