# v0.1.2

## Resolved issues

- Implement new algorithm for splitting test names and rebuilding test
  hierarchy, that supports '/' in test names.
- Implement better testing of JSON test result parsing, covering explicitly
  extraction of failures, coverage and elapsed time.
- Proceed with output generation when `go test` find test failures and produces
  a non-zero exit code.

# v0.1.1

First usable version of the command with full support for yaml, json and running
`go test` as input, yaml, markdown and custom template as output.
