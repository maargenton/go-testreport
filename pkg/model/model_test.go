package model_test

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/subexpr"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-testreport/pkg/model"
)

func samplePkg() model.Package {
	var pkg = model.Package{
		Name: "pkg1",
		Tests: []*model.Test{
			{
				Name: "foo",
				Tests: []*model.Test{
					{
						Name: "bar",
						Tests: []*model.Test{
							{
								Name: "baz",
							},
						},
					},
					{
						Name: "bar2",
						Tests: []*model.Test{
							{
								Name: "baz",
							},
						},
					},
				},
			},
		},
	}
	pkg.LinkTests()
	return pkg
}

func TestPackage(t *testing.T) {
	bdd.Given(t, "a package with tests", func(t *bdd.T) {
		pkg := samplePkg()
		t.When("enumerating leaf tests", func(t *bdd.T) {
			tests := pkg.LeafTests()

			t.Then("only tests with no sub-tests are reported", func(t *bdd.T) {
				verify.That(t, tests).Length().Eq(2)
				verify.That(t, tests).All(
					subexpr.Value().Field("Tests").IsEmpty(),
				)
			})
			t.Then("full names include full test hierarchy names", func(t *bdd.T) {
				var names []string
				for _, t := range tests {
					names = append(names, t.FullName())
				}
				verify.That(t, names).Eq([]string{
					"foo, bar, baz",
					"foo, bar2, baz",
				})
			})
			t.Then("partial names include partial test hierarchy names", func(t *bdd.T) {
				var names []string
				for _, t := range tests {
					names = append(names, t.PartialName(1))
				}
				verify.That(t, names).Eq([]string{
					"bar, baz",
					"bar2, baz",
				})
			})
			t.Then("partial names are empty when skipping too much", func(t *bdd.T) {
				var names []string
				for _, t := range tests {
					names = append(names, t.PartialName(10))
				}
				verify.That(t, names).Eq([]string{"", ""})
			})
		})
	})
}
