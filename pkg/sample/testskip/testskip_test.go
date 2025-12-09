//go:build sample

package testskip_test

import (
	"fmt"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"

	sample "github.com/maargenton/go-testreport/pkg/sample/testskip"
)

func TestError(t *testing.T) {
	t.Skip("Skipping TestBar due to known issue")

	var tcs = []struct {
		name string
	}{
		{"/api/v1/foo"},
		{"/api/v1/baz"},
		{"/api/v1/bar"},
	}

	bdd.Wrap(t, "TestFoo", func(t *bdd.T) {
		for _, tc := range tcs {
			name := fmt.Sprintf("GET %v", tc.name)
			t.Run(name, func(t *bdd.T) {
				require.That(t, sample.Foo(name)).Eq(name)
			})
		}
	})
	bdd.Wrap(t, "TestBar", func(t *bdd.T) {
		for _, tc := range tcs {
			name := fmt.Sprintf("GET %v", tc.name)
			t.Run(name, func(t *bdd.T) {
				require.That(t, sample.Bar(name)).Ne(name)
			})
		}
	})
}
