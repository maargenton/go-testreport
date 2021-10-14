package sample_test

import (
	"fmt"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-testreport/pkg/sample"
)

func TestFoo(t *testing.T) {

	// Package sample and sample_test contains sample failing tests and poor
	// coverage for the purpose of generating sample output for go-testreport.
	// To regenerate the output, comment out the t.Skip() line below and run:
	//   go test -race -cover -json ./pkg/sample >pkg/model/testdata/sample-output.json

	var tcs = []struct {
		name string
	}{
		{"/api/v1/foo"},
		{"/api/v1/baz"},
		{"/api/v1/bar"},
	}

	bdd.Wrap(t, "Test1", func(t *bdd.T) {
		for _, tc := range tcs {
			name := fmt.Sprintf("GET %v", tc.name)
			t.Run(name, func(t *bdd.T) {
				require.That(t, sample.Foo(name)).Eq(name)
			})
		}
	})
	bdd.Wrap(t, "Test2", func(t *bdd.T) {

		t.Skip() // Comment out to generate test failures

		for _, tc := range tcs {
			name := fmt.Sprintf("GET %v", tc.name)
			t.Run(name, func(t *bdd.T) {
				require.That(t, sample.Foo(name)).Ne(name)
			})
		}
	})
}
