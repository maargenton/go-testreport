package gotest

import (
	"encoding/json"
	"io"
	"regexp"

	"github.com/maargenton/go-fileutils"
)

// JSONOutput represents the full parse output of a `go test` run in JSON
// format, separating build events from package and test level events. It is
// initially a flat list of tests in a flat list of packages, but with all the
// necessary fields to be reorganized into the original hierarchical structure.
type JSONOutput struct {
	PackageOrder []string                    // List of packages in the order they appear
	Packages     map[string]*JSONPackage     // Test events grouped by `Package`
	ImportPaths  map[string]*JSONBuildEvents // Build events grouped by `ImportPath`
}

// JSONBuildEvents represents all the reported build events for a specific import
// path.
type JSONBuildEvents struct {
	ImportPath string           // The package being built
	Lines      []JSONOutputLine // List of all JSON output lines for the build event
}

// JSONPackage represent the package and test level events for a package,
// grouped by test
type JSONPackage struct {
	PackageName string           // Name of the package
	Lines       []JSONOutputLine // List of all JSON output lines not assigned to any specific test
	Tests       []*JSONTest      // List of tests in the package

	testMap map[string]*JSONTest // Map of tests by test name
}

// JSONTest represents the events related to a single test, with additional
// fields to reconstruct the original test hierarchy.
type JSONTest struct {
	TestName string           // Name of the test
	Lines    []JSONOutputLine // List of all JSON output lines for the test
	Parent   *JSONTest        // Parent test if this is a nested test
	Nested   []*JSONTest      // List of nested tests within this test
}

// JSONOutputLine represents a single line of go test JSON output, and can be
// either a build action or a test action.
type JSONOutputLine struct {
	Action     string `json:"Action"`     // Either build action or test action
	Package    string `json:"Package"`    // The package being tested in test actions
	ImportPath string `json:"ImportPath"` // The package being built for build actions

	Output      string  `json:"Output"`      // The output text from the test or build action
	OutputType  string  `json:"OutputType"`  // The type of output, but does not seem to be included
	Test        string  `json:"Test"`        // The name of the test
	Elapsed     float64 `json:"Elapsed"`     // The time taken to execute the test
	FailedBuild string  `json:"FailedBuild"` // Ties back to BuildEvent.ImportPath
}

// addBuildEvent records a build event in the current JSONOutput
func (output *JSONOutput) addBuildEvent(line JSONOutputLine) {
	var importPath = cleanupImportPath(line.ImportPath)
	var buildEvents, exists = output.ImportPaths[importPath]
	if !exists {
		buildEvents = &JSONBuildEvents{ImportPath: importPath}
		output.ImportPaths[importPath] = buildEvents
	}
	line.ImportPath = importPath
	buildEvents.Lines = append(buildEvents.Lines, line)
}

// cleanupImportPath removes any trailing build/test information from the import
// path
func cleanupImportPath(importPath string) string {
	var re = regexp.MustCompile(`\s+\[`)
	var parts = re.Split(importPath, 2)
	return parts[0]
}

// addTestEvent records a test event in the current JSONOutput, grouped by
// package and test name.
func (output *JSONOutput) addTestEvent(line JSONOutputLine) {
	var packageName = line.Package
	var testName = line.Test

	var pkg, pkgExists = output.Packages[packageName]
	if !pkgExists {
		pkg = &JSONPackage{
			PackageName: packageName,
			testMap:     make(map[string]*JSONTest),
		}
		output.PackageOrder = append(output.PackageOrder, packageName)
		output.Packages[packageName] = pkg
	}

	if testName == "" {
		pkg.Lines = append(pkg.Lines, line)
	} else {
		var test, testExists = pkg.testMap[testName]
		if !testExists {
			test = &JSONTest{
				TestName: testName,
			}
			pkg.Tests = append(pkg.Tests, test)
			pkg.testMap[testName] = test
		}
		test.Lines = append(test.Lines, line)
	}
}

// LoadJSONOutput loads a go-test JSON output from a file using ParseJSONOutput
// and generating same output.
func LoadJSONOutput(filename string) (output JSONOutput, err error) {
	err = fileutils.ReadFile(filename, func(r io.Reader) error {
		output, err = ParseJSONOutput(r)
		return err
	})
	return
}

// ParseJSONOutput parses a go-test JSON output from an io.Reader into a flat
// list of tests, grouped by package, and build events groped by import path.
func ParseJSONOutput(r io.Reader) (output JSONOutput, err error) {

	output.ImportPaths = make(map[string]*JSONBuildEvents)
	output.Packages = make(map[string]*JSONPackage)

	var decoder = json.NewDecoder(r)
	for {
		var line JSONOutputLine
		err = decoder.Decode(&line)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}

		if line.ImportPath != "" {
			output.addBuildEvent(line)
		} else {
			output.addTestEvent(line)
		}
	}
	return
}
