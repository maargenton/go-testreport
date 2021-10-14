package model

import (
	"sort"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

func TestSplitTestNames(t *testing.T) {
	bdd.Given(t, "a set of test names", func(t *bdd.T) {
		var names = []string{
			"Test",
			"Test/Test1",
			"Test/Test1/GET_/api/v1/foo",
			"Test/Test1#01",
			"Test/Test1#01/GET_/api/v1/bar",
			"Test/Test1#02",
			"Test/Test1#02/GET_/api/v1/baz",
			"Test/Test2",
			"Test/Test2/GET_/api/v1/foo",
			"Test/Test2#01",
			"Test/Test2#01/GET_/api/v1/bar",
			"Test/Test2#02",
			"Test/Test2#02/GET_/api/v1/baz",
		}
		sort.Strings(names)

		t.When("calling generateTestNameSplits()", func(t *bdd.T) {
			splits := generateTestNameSplits(names)

			t.Then("longest names are fully split", func(t *bdd.T) {
				var tcs = []struct {
					name  string
					split []string
				}{
					{"Test/Test1/GET_/api/v1/foo", []string{"Test", "Test1", "GET_/api/v1/foo"}},
					{"Test/Test1#01/GET_/api/v1/bar", []string{"Test", "Test1#01", "GET_/api/v1/bar"}},
					{"Test/Test1#02/GET_/api/v1/baz", []string{"Test", "Test1#02", "GET_/api/v1/baz"}},
					{"Test/Test2/GET_/api/v1/foo", []string{"Test", "Test2", "GET_/api/v1/foo"}},
					{"Test/Test2#01/GET_/api/v1/bar", []string{"Test", "Test2#01", "GET_/api/v1/bar"}},
					{"Test/Test2#02/GET_/api/v1/baz", []string{"Test", "Test2#02", "GET_/api/v1/baz"}},
				}
				for _, tc := range tcs {
					require.That(t, splits[tc.name]).Eq(tc.split)
				}
			})
		})
	})
}
