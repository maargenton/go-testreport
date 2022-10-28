package model

import (
	"strings"
	"time"
)

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

func (p *Package) linkTests() {
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

// FullName returns the name of all the tests in a branch, concatenated by a
// comma.
func (t *Test) FullName() string {
	if t.Parent == nil {
		return t.Name
	}

	return t.Parent.FullName() + ", " + t.Name
}

// PartialName like FullName returns the name of all the tests in a branch,
// concatenated by a comma, but skipping the requested upper level tests.
func (t *Test) PartialName(skip int) string {
	var parts []string
	for tt := t; tt != nil; tt = tt.Parent {
		parts = append([]string{tt.Name}, parts...)
	}

	if skip > 0 && skip < len(parts) {
		parts = parts[skip:]
	} else {
		parts = nil
	}
	return strings.Join(parts, ", ")
}
