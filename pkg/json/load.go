package json

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

	"github.com/maargenton/go-testreport/pkg/model"
)

func Load(r io.Reader) ([]model.Package, error) {
	pkgMap, err := parseTestOutput(r)
	if err != nil {
		return nil, err
	}

	packageNames := make([]string, 0, len(pkgMap))
	for pkg := range pkgMap {
		packageNames = append(packageNames, pkg)
	}
	sort.Strings(packageNames)

	var packages = make([]model.Package, 0, len(pkgMap))
	for _, name := range packageNames {
		records := rebuildTestHierarchy(pkgMap[name])
		tests := records.toTests()

		var skipped = false
		var elapsed = 0 * time.Second
		var coverage = 0.0

		pkgRecords := records.m[""]
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

		for _, t := range tests {
			t.LinkSubTests()
		}

		packages = append(packages, model.Package{
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
	var records = newTestRecord()
	for _, l := range lines {
		if isDiscardable(l) {
			continue
		}
		nameParts := splitTestName(l.Test)
		records.recordTest(nameParts, l)
	}
	return records
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

func splitTestName(name string) []string {
	parts := strings.Split(name, "/")
	for i, p := range parts {
		j := strings.LastIndexByte(p, '#')
		if j >= 0 {
			n := p[j+1:]
			if _, err := strconv.ParseInt(n, 10, 64); err == nil {
				parts[i] = p[:j]
			}
		}
	}
	return parts
}

type testRecord struct {
	name    string
	l       []*testRecord
	m       map[string]*testRecord
	details []jsonInputLine
}

func newTestRecord() *testRecord {
	return &testRecord{m: make(map[string]*testRecord)}
}

func (r *testRecord) recordTest(name []string, details jsonInputLine) {
	if len(name) > 0 {
		n := name[0]
		if rr, ok := r.m[n]; ok {
			rr.recordTest(name[1:], details)
		} else {
			rr := newTestRecord()
			rr.name = n
			r.m[n] = rr
			r.l = append(r.l, rr)
			rr.recordTest(name[1:], details)
		}
	} else {
		r.details = append(r.details, details)
	}
}

func (r *testRecord) toTests() (tests []*model.Test) {
	for _, t := range r.l {
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

		tests = append(tests, &model.Test{
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
