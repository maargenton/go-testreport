package model

import (
	"slices"
	"strings"
	"time"
	"unicode"
)

// Results collects all the test record from all tested packages
type Results struct {
	Success  bool       `yaml:"success"`
	Passed   int        `yaml:"passed"`
	Failed   int        `yaml:"failed"`
	Packages []*Package `yaml:"packages"`
}

// UpdateCounts updates the internal references of collected results and the
// total counts of passed and failed leaf tests. It should be called for
// consistency once the model has been fully constructed or changed.
func (r *Results) UpdateCounts() {
	var passed = 0
	var failed = 0

	for _, pkg := range r.Packages {
		pkg.updateCounts()
		passed += pkg.Passed
		failed += pkg.Failed
	}
	r.Passed = passed
	r.Failed = failed
	r.Success = failed == 0
}

// Package collects all the test record for one go package within the project
type Package struct {
	Name     string        `yaml:"package"`
	Elapsed  time.Duration `yaml:"elapsed"`
	Passed   int           `yaml:"passed"`
	Failed   int           `yaml:"failed"`
	Coverage float64       `yaml:"coverage"`
	Skipped  bool          `yaml:"skipped"`
	Tests    []*Test       `yaml:"tests,omitempty"`
}

func (p *Package) updateCounts() {
	for _, t := range p.Tests {
		t.linkTests()
	}

	var passed = 0
	var failed = 0
	for _, t := range p.LeafTests() {
		if t.Failure {
			failed++
		} else {
			passed++
		}
	}
	p.Passed = passed
	p.Failed = failed
}

// LeafTests returns the set of tests from the package that don't have any
// sub-test.
func (p *Package) LeafTests() (r []*Test) {
	for _, tt := range p.Tests {
		r = append(r, tt.LeafTests()...)
	}
	return
}

// Test collects nested test records and outputs
type Test struct {
	Name    string   `yaml:"name"`
	Failure bool     `yaml:"failure,omitempty"`
	Output  []string `yaml:"output,omitempty"`
	Parent  *Test    `yaml:"-"`
	Tests   []*Test  `yaml:"tests,omitempty"`
}

func (t *Test) linkTests() {
	for _, tt := range t.Tests {
		tt.Parent = t
		tt.linkTests()
	}
}

// LeafTests returns the set of tests under the current test that don't have any
// sub-test. That may include the test itself if it has not sub-test.
func (t *Test) LeafTests() (r []*Test) {
	if len(t.Tests) == 0 {
		return []*Test{t}
	}

	for _, tt := range t.Tests {
		r = append(r, tt.LeafTests()...)
	}
	return
}

// FullName returns the full name of the test as the concatenation of all the
// name fragments in the test hierarchy, using either comma or colon as separator
// depending on the capitalisation  of the next name fragment.
func (t *Test) FullName() string {
	return t.PartialName(0)
}

// PartialName return the full concatenated name of the test just like
// `FullName()`, but skips the requested number of fragments from the beginning
// of the test hierarchy.
func (t *Test) PartialName(skip int) string {
	var parts []string
	for tt := t; tt != nil; tt = tt.Parent {
		var name = strings.TrimSpace(tt.Name)
		if len(name) > 0 {
			parts = append(parts, name)
		}
	}
	slices.Reverse(parts)

	if skip >= 0 && skip < len(parts) {
		parts = parts[skip:]
	} else {
		parts = nil
	}

	var result strings.Builder
	for i, s := range parts {
		result.WriteString(s)
		if i < len(parts)-1 {
			var next = parts[i+1]
			if len(next) > 0 && unicode.IsUpper(rune(next[0])) {
				result.WriteString(": ")
			} else {
				result.WriteString(", ")
			}
		}
	}
	return result.String()
}
