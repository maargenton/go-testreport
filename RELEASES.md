# v0.1.6

- Maintenance release, cleanup CI and build reporting
- Cleanup rakefile to be more usable as baseline
- Add package summary to default markdown report

## Code changes

- Cleanup CI and templates ([#7](https://github.com/maargenton/go-testreport/pull/7))


# v0.1.5

- Add more builtin markdown templates

## Code changes

- Add more builtin markdown templates ([#6](https://github.com/maargenton/go-testreport/pull/6))


# v0.1.4

- Maintenance release, cleanup CI build reporting


# v0.1.3

- Enhance package model with `Passed` and `Failed` as a count passed and failed
  leaf tests to help with summary and failure only outputs
- Enhance default markdown template with:
  - a summary table
  - a failures section reporting only failures
  - a detailed section containing the full report, but collapsed by default
- Fix linking

# v0.1.2

## Resolved issues

- Implement new algorithm for splitting test names and rebuilding test
  hierarchy, that supports '/' in test names.
- Implement better testing of JSON test result parsing, covering explicitly
  extraction of failures, coverage and elapsed time.
- Proceed with output generation when `go test` finds test failures and produces
  a non-zero exit code.

# v0.1.1

First usable version of the command with full support for yaml, json and running
`go test` as input, yaml, markdown and custom template as output.
