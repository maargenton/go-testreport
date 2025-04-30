package template

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

func Test_regexMatch(t *testing.T) {
	bdd.Given(t, "a valid regex", func(t *bdd.T) {
		var regex = `^([a-z]+)$`
		t.When("calling regexMatch() with a matching string", func(t *bdd.T) {
			t.Then("it returns true", func(t *bdd.T) {
				var match, err = regexMatch(regex, "abc")
				require.That(t, err).IsError(nil)
				require.That(t, match).IsTrue()
			})
		})
		t.When("calling regexMatch() with a non-matching string", func(t *bdd.T) {
			t.Then("it returns false", func(t *bdd.T) {
				var match, err = regexMatch(regex, "ABC")
				require.That(t, err).IsError(nil)
				require.That(t, match).IsFalse()
			})
		})
	})
	bdd.Given(t, "an invalid regex", func(t *bdd.T) {
		var regex = `^([a-z+)$`
		t.When("calling regexMatch() with any string", func(t *bdd.T) {
			t.Then("it returns an error", func(t *bdd.T) {
				var _, err = regexMatch(regex, "abc")
				require.That(t, err).IsError("error parsing regex")
			})
		})
	})
}

func Test_regexFind(t *testing.T) {
	bdd.Given(t, "a valid regex", func(t *bdd.T) {
		var regex = `\d+`
		t.When("calling regexFind() with a matching string", func(t *bdd.T) {
			t.Then("it returns the matching string", func(t *bdd.T) {
				var match, err = regexFind(regex, "abc123def")
				require.That(t, err).IsError(nil)
				require.That(t, match).Eq("123")
			})
		})
		t.When("calling regexFind() with a non-matching string", func(t *bdd.T) {
			t.Then("it returns an empty string", func(t *bdd.T) {
				var match, err = regexFind(regex, "abcdef")
				require.That(t, err).IsError(nil)
				require.That(t, match).Eq("")
			})
		})
	})
	bdd.Given(t, "an invalid regex", func(t *bdd.T) {
		var regex = `^([a-z+)$`
		t.When("calling regexFind() with any string", func(t *bdd.T) {
			t.Then("it returns an error", func(t *bdd.T) {
				var _, err = regexFind(regex, "abc")
				require.That(t, err).IsError("error parsing regex")
			})
		})
	})
}

func Test_regexFindAll(t *testing.T) {
	var regex = `\[([^\]]+)\]`
	var matchingInput = "Given something, when something [gg-123][gg-124], then something"
	var notMatchingInput = "Given something, when something, then something"

	bdd.Given(t, "a valid regex", func(t *bdd.T) {
		t.When("calling regexFind() with a matching string", func(t *bdd.T) {
			t.Then("it returns the matching string", func(t *bdd.T) {
				var match, err = regexFindAll(regex, matchingInput)
				require.That(t, err).IsError(nil)
				require.That(t, match).Eq([]string{"[gg-123]", "[gg-124]"})
			})
		})
		t.When("calling regexFind() with a non-matching string", func(t *bdd.T) {
			t.Then("it returns an empty string", func(t *bdd.T) {
				var match, err = regexFindAll(regex, notMatchingInput)
				require.That(t, err).IsError(nil)
				require.That(t, match).Eq("")
			})
		})
	})
	bdd.Given(t, "an invalid regex", func(t *bdd.T) {
		var regex = `^([a-z+)$`
		t.When("calling regexFind() with any string", func(t *bdd.T) {
			t.Then("it returns an error", func(t *bdd.T) {
				var _, err = regexFindAll(regex, "abc")
				require.That(t, err).IsError("error parsing regex")
			})
		})
	})
}

func Test_regexFindSubmatch(t *testing.T) {
	var regex = `\[([^\]]+)\]`
	var matchingInput = "Given something, when something [gg-123][gg-124], then something"
	var notMatchingInput = "Given something, when something, then something"

	bdd.Given(t, "a valid regex", func(t *bdd.T) {
		t.When("calling regexFindSubmatch() with a matching string", func(t *bdd.T) {
			t.Then("it returns the requested capture group from the match", func(t *bdd.T) {
				var match, err = regexFindSubmatch(regex, 1, matchingInput)
				require.That(t, err).IsError(nil)
				require.That(t, match).Eq("gg-123")
			})

			t.With("0 index", func(t *bdd.T) {
				t.Then("it returns the full match", func(t *bdd.T) {
					var match, err = regexFindSubmatch(regex, 0, matchingInput)
					require.That(t, err).IsError(nil)
					require.That(t, match).Eq("[gg-123]")
				})
			})

			t.With("an out-of-bound capture group index", func(t *bdd.T) {
				t.Then("it returns an error", func(t *bdd.T) {
					var match, err = regexFindSubmatch(regex, 2, matchingInput)
					require.That(t, err).IsError("invalid submatch index")
					require.That(t, match).Eq("")
				})
			})
		})
		t.When("calling regexFindSubmatch() with a non-matching string", func(t *bdd.T) {
			t.Then("it returns an empty string", func(t *bdd.T) {
				var match, err = regexFindSubmatch(regex, 1, notMatchingInput)
				require.That(t, err).IsError(nil)
				require.That(t, match).Eq("")
			})
		})
	})

	bdd.Given(t, "an invalid regex", func(t *bdd.T) {
		var regex = `^([a-z+)$`
		t.When("calling regexFindSubmatch() with any string", func(t *bdd.T) {
			t.Then("it returns an error", func(t *bdd.T) {
				var _, err = regexFindSubmatch(regex, 1, "abc")
				require.That(t, err).IsError("error parsing regex")
			})
		})
	})
}

func Test_regexFindAllSubmatch(t *testing.T) {
	var regex = `\[([^\]]+)\]`
	var matchingInput = "Given something, when something [gg-123][gg-124], then something"
	var notMatchingInput = "Given something, when something, then something"

	bdd.Given(t, "a valid regex", func(t *bdd.T) {
		t.When("calling regexFindAllSubmatch() with a matching string", func(t *bdd.T) {
			t.Then("it returns the requested capture group from each match", func(t *bdd.T) {
				var match, err = regexFindAllSubmatch(regex, 1, matchingInput)
				require.That(t, err).IsError(nil)
				require.That(t, match).Eq([]string{"gg-123", "gg-124"})
			})

			t.With("0 index", func(t *bdd.T) {
				t.Then("it returns each full match", func(t *bdd.T) {
					var match, err = regexFindAllSubmatch(regex, 0, matchingInput)
					require.That(t, err).IsError(nil)
					require.That(t, match).Eq([]string{"[gg-123]", "[gg-124]"})
				})
			})

			t.With("an out-of-bound capture group index", func(t *bdd.T) {
				t.Then("it returns an error", func(t *bdd.T) {
					var match, err = regexFindAllSubmatch(regex, 2, matchingInput)
					require.That(t, err).IsError("invalid submatch index")
					require.That(t, match).Eq("")
				})
			})
		})
		t.When("calling regexFindAllSubmatch() with a non-matching string", func(t *bdd.T) {
			t.Then("it returns an empty string", func(t *bdd.T) {
				var match, err = regexFindAllSubmatch(regex, 1, notMatchingInput)
				require.That(t, err).IsError(nil)
				require.That(t, match).Eq([]string{})
			})
		})
	})

	bdd.Given(t, "an invalid regex", func(t *bdd.T) {
		var regex = `^([a-z+)$`
		t.When("calling regexFindAllSubmatch() with any string", func(t *bdd.T) {
			t.Then("it returns an error", func(t *bdd.T) {
				var _, err = regexFindAllSubmatch(regex, 1, "abc")
				require.That(t, err).IsError("error parsing regex")
			})
		})
	})
}

func Test_regexReplaceAll(t *testing.T) {
	var regex = `\s*\[([^\]]+)\]`
	var matchingInput = "Given something, when something [gg-123][gg-124], then something"
	var notMatchingInput = "Given something, when something, then something"

	bdd.Given(t, "a valid regex", func(t *bdd.T) {
		t.When("calling regexReplaceAll() with an empty replacement and a matching string", func(t *bdd.T) {
			t.Then("it returns the input with matches removed", func(t *bdd.T) {
				var result, err = regexReplaceAll(regex, "", matchingInput)
				require.That(t, err).IsError(nil)
				require.That(t, result).Eq(notMatchingInput)
			})
		})
		t.When("calling regexReplaceAll()", func(t *bdd.T) {
			var regex = `\[([^\]]+)\]`
			t.With("a callout replacement and a matching string", func(t *bdd.T) {
				t.Then("it returns the input with matches replaced", func(t *bdd.T) {
					var result, err = regexReplaceAll(regex, "($1)", matchingInput)
					require.That(t, err).IsError(nil)
					require.That(t, result).Eq(
						"Given something, when something (gg-123)(gg-124), then something",
					)
				})
			})
			t.With("an out of bound callout replacement", func(t *bdd.T) {
				t.Then("it returns the input with matches replaced with empty submatch", func(t *bdd.T) {
					var result, err = regexReplaceAll(regex, "($2)", matchingInput)
					require.That(t, err).IsError(nil)
					require.That(t, result).Eq(
						"Given something, when something ()(), then something",
					)
				})
			})
		})
	})

	bdd.Given(t, "an invalid regex", func(t *bdd.T) {
		var regex = `^([a-z+)$`
		t.When("calling regexReplaceAll() with any string", func(t *bdd.T) {
			t.Then("it returns an error", func(t *bdd.T) {
				var _, err = regexReplaceAll(regex, "", "abc")
				require.That(t, err).IsError("error parsing regex")
			})
		})
	})
}

func Test_regexSplit(t *testing.T) {
	var regex = `\s*,\s*`
	var matchingInput = "foo, bar, baz, foobar"
	// var notMatchingInput = "Given something, when something, then something"

	bdd.Given(t, "a valid regex", func(t *bdd.T) {
		t.When("calling regexSplit() with a matching string", func(t *bdd.T) {
			t.Then("it returns the input split around matches", func(t *bdd.T) {
				var result, err = regexSplit(regex, matchingInput)
				require.That(t, err).IsError(nil)
				require.That(t, result).Eq([]string{"foo", "bar", "baz", "foobar"})
			})
		})
	})

	bdd.Given(t, "an invalid regex", func(t *bdd.T) {
		var regex = `^([a-z+)$`
		t.When("calling regexSplit() with any string", func(t *bdd.T) {
			t.Then("it returns an error", func(t *bdd.T) {
				var _, err = regexSplit(regex, "abc")
				require.That(t, err).IsError("error parsing regex")
			})
		})
	})
}

func Test_regexSplitN(t *testing.T) {
	var regex = `\s*,\s*`
	var matchingInput = "foo, bar, baz, foobar"
	// var notMatchingInput = "Given something, when something, then something"

	bdd.Given(t, "a valid regex", func(t *bdd.T) {
		t.When("calling regexSplitN() with a matching string", func(t *bdd.T) {
			t.Then("it limits the split to the requested count", func(t *bdd.T) {
				var result, err = regexSplitN(regex, 2, matchingInput)
				require.That(t, err).IsError(nil)
				require.That(t, result).Length().Eq(2)
				require.That(t, result).Eq([]string{"foo", "bar, baz, foobar"})
			})
		})
	})

	bdd.Given(t, "an invalid regex", func(t *bdd.T) {
		var regex = `^([a-z+)$`
		t.When("calling regexSplitN() with any string", func(t *bdd.T) {
			t.Then("it returns an error", func(t *bdd.T) {
				var _, err = regexSplitN(regex, 2, "abc")
				require.That(t, err).IsError("error parsing regex")
			})
		})
	})
}
