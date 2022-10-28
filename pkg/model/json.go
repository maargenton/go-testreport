package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/maargenton/go-fileutils"
)

// LoadFromGoTestJson loads output of `go test` using JSON format into the
// internal representation of test results. The process involves identifying
// packages and tests, grouping and linking nested tests, cleaning up test names
// and test output, and capturing coverage, elapsed time, and success / failure.
func LoadFromGoTestJson(r io.Reader) (pkgs []Package, err error) {
	pkgs, err = load(r)
	if err != nil {
		return nil, err
	}
	for i := range pkgs {
		pkgs[i].linkTests()
	}
	return
}

// LoadFromGoTestJsonFile is similar to LoadFromGoTestJson, but loading from a
// file instead of an `io.Reader`.
func LoadFromGoTestJsonFile(filename string) (pkgs []Package, err error) {
	err = fileutils.ReadFile(filename, func(r io.Reader) error {
		pkgs, err = LoadFromGoTestJson(r)
		return err
	})
	return
}

func load(r io.Reader) ([]Package, error) {
	pkgMap, err := parseTestOutput(r)
	if err != nil {
		return nil, err
	}

	packageNames := make([]string, 0, len(pkgMap))
	for pkg := range pkgMap {
		packageNames = append(packageNames, pkg)
	}
	sort.Strings(packageNames)

	var packages = make([]Package, 0, len(pkgMap))
	for _, name := range packageNames {
		records := rebuildTestHierarchy(pkgMap[name])
		tests := records.toTests()

		var skipped = false
		var elapsed = 0 * time.Second
		var coverage = 0.0

		pkgRecords := records.nestedMap[""]
		if pkgRecords != nil {
			for _, l := range pkgRecords.details {
				if l.Action == "skip" {
					skipped = true
				}
				if l.Elapsed != 0 {
					elapsed = time.Duration(l.Elapsed * float64(time.Second))
				}
				if l.Action == "output" {
					if cov, ok := parseCoverageOutput(l.Output); ok {
						coverage = cov
					}
				}
			}
		}

		packages = append(packages, Package{
			Name:     name,
			Tests:    tests,
			Elapsed:  elapsed,
			Coverage: coverage,
			Skipped:  skipped,
		})

	}
	return packages, nil
}

func parseTestOutput(r io.Reader) (map[string][]jsonInputLine, error) {
	var result = make(map[string][]jsonInputLine)
	scanner := bufio.NewScanner(r)
	lineno := 0
	for scanner.Scan() {
		lineno++
		l := scanner.Bytes()
		input := jsonInputLine{}
		err := json.Unmarshal(l, &input)
		if err != nil {
			return nil, fmt.Errorf("parse error on line %v: %w", lineno, err)
		}
		pkg := input.Package
		result[pkg] = append(result[pkg], input)
	}
	return result, nil
}

func rebuildTestHierarchy(lines []jsonInputLine) *testRecord {
	var nameSet = make(map[string]struct{})
	for _, l := range lines {
		name := l.Test
		if name != "" {
			nameSet[name] = struct{}{}
		}
	}
	var names []string
	for name := range nameSet {
		if name != "" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	var splitNames = generateTestNameSplits(names)

	var records = newTestRecord()
	for _, l := range lines {
		if isDiscardable(l) {
			continue
		}
		nameParts := splitNames[l.Test]
		records.recordTest(nameParts, l)
	}
	return records
}

// generateTestNameSplits takes a complete sorted list of test and sub-test
// names and produces a map between a test name and a split version of that
// name where the prefix is another test.
func generateTestNameSplits(names []string) map[string][]string {
	var splits = make(map[string][]string)
	splits[""] = []string{""}
	for _, name := range names {
		splits[name] = []string{name}
	}

	// Look for longest prefix in names list, spliting on '/' only
	for _, name := range names {
		var i = len(name)
		for {
			i = strings.LastIndexByte(name[:i], '/')
			if i < 0 {
				break
			}
			if _, ok := splits[name[:i]]; ok {
				splits[name] = []string{name[:i], name[i+1:]}
				break
			}
		}
	}

	// Replace prefix with split prefix if any. NOTE: A single through sorted
	// list is enough, as prefixes will appear first, even if not directly
	// adjacent.
	for _, name := range names {
		var split = splits[name]
		if len(split) > 1 && len(splits[split[0]]) > 1 {
			var expanded = append([]string{}, splits[split[0]]...)
			expanded = append(expanded, split[1:]...)
			splits[name] = expanded
		}
	}

	return splits
}

// ---------------------------------------------------------------------------
// Private helpers to load test output into a structure test hierarchy
// ---------------------------------------------------------------------------

type jsonInputLine struct {
	Action  string
	Package string
	Test    string
	Output  string
	Elapsed float64
}

func isDiscardable(l jsonInputLine) bool {
	o := strings.TrimSpace(l.Output)
	return l.Action == "output" &&
		(strings.HasPrefix(o, "===") || strings.HasPrefix(o, "---"))
}

func stripIndexSuffix(name string) string {
	j := strings.LastIndexByte(name, '#')
	if j >= 0 {
		n := name[j+1:]
		if _, err := strconv.ParseInt(n, 10, 64); err == nil {
			return name[:j]
		}
	}
	return name
}

type testRecord struct {
	name       string
	nestedList []*testRecord
	nestedMap  map[string]*testRecord
	details    []jsonInputLine
}

func newTestRecord() *testRecord {
	return &testRecord{nestedMap: make(map[string]*testRecord)}
}

func (r *testRecord) recordTest(name []string, details jsonInputLine) {
	if len(name) > 0 {
		n := stripIndexSuffix(name[0])
		if rr, ok := r.nestedMap[n]; ok {
			rr.recordTest(name[1:], details)
		} else {
			rr := newTestRecord()
			rr.name = n
			r.nestedMap[n] = rr
			r.nestedList = append(r.nestedList, rr)
			rr.recordTest(name[1:], details)
		}
	} else {
		r.details = append(r.details, details)
	}
}

func (r *testRecord) toTests() (tests []*Test) {
	for _, t := range r.nestedList {
		if t.name == "" {
			continue
		}

		var success = true
		for _, l := range t.details {
			if l.Action == "fail" {
				success = false
			}
		}

		var output []string
		for _, l := range t.details {
			if l.Action == "output" {
				output = append(output, l.Output)
			}
		}

		tests = append(tests, &Test{
			Name:    cleanupTestName(t.name),
			Tests:   t.toTests(),
			Failure: !success,
			Output:  cleanupOutputs(output),
		})
	}
	return
}

func parseCoverageOutput(s string) (float64, bool) {
	var i = strings.Index(s, "coverage:")
	if i >= 0 && i+10 <= len(s) {
		s = s[i+10:]
		i = strings.IndexByte(s, '%')
		if i >= 0 {
			s = s[:i]

			if cov, err := strconv.ParseFloat(s, 64); err == nil {
				return cov, true
			}
		}
	}
	return 0, false
}

func cleanupTestName(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "_", " "))
}

func cleanupOutputs(lines []string) []string {
	var result = make([]string, 0, len(lines))
	var prefix = ""
	if len(lines) > 0 {
		l0 := lines[0]
		i := strings.IndexFunc(l0, func(r rune) bool { return !unicode.IsSpace(r) })
		if i >= 0 {
			prefix = l0[:i]
		}
	}

	for _, l := range lines {
		l = strings.Replace(l, prefix, "", 1)
		l = strings.TrimRightFunc(l, unicode.IsSpace)
		result = append(result, l)
	}

	return result
}
