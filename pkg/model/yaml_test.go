package model_test

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testreport/pkg/model"
)

func TestLoadFromYAMLFile(t *testing.T) {
	bdd.Given(t, "a YAML file containing tests", func(t *bdd.T) {
		filename := "./testdata/sample-failure.yaml"

		t.When("calling LoadFromYAMLFile()", func(t *bdd.T) {
			pkgs, err := model.LoadFromYAMLFile(filename)

			t.Then("something happens", func(t *bdd.T) {
				require.That(t, err).IsError(nil)
				require.That(t, pkgs).Length().Eq(1)
			})
		})
	})
}
