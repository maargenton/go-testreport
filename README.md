# go-testreport

Test report documentation generator for Go.

[![Latest](
  https://img.shields.io/github/v/tag/maargenton/go-testreport?color=blue&label=latest&logo=go&logoColor=white&sort=semver)](
  https://pkg.go.dev/github.com/maargenton/go-testreport)
[![Build](
  https://img.shields.io/github/workflow/status/maargenton/go-testreport/build?label=build&logo=github&logoColor=aaaaaa)](
  https://github.com/maargenton/go-testreport/actions?query=branch%3Amaster)
[![Codecov](
  https://img.shields.io/codecov/c/github/maargenton/go-testreport?label=codecov&logo=codecov&logoColor=aaaaaa&token=fVZ3ZMAgfo)](
  https://codecov.io/gh/maargenton/go-testreport)
[![Go Report Card](
  https://goreportcard.com/badge/github.com/maargenton/go-testreport)](
  https://goreportcard.com/report/github.com/maargenton/go-testreport)


---------------------------

## Usage

```
go install github.com/margenton/go-testreport

go-testreport --race ./... \
    -oyaml=build/build-report.yaml \
    -omarkdown=build/build-report.yaml

go run github.com/maargenton/go-testreport@latest --race ./... -omd=test-report.md
```

The `yaml` output is a straight representation of the internal model, with a
list of packages; each package as a list of tests and each test can have nested
tests. The `markdown` output is generated from a built-in template applied to
this internal model. Custom outputs can be produces with custom templates (see
[Using Custom Template](#using-custom-template))

### Disclaimer

The focus of this tool is to help document the test that have been implemented
on a codebase, and therefore are all passing on a successful build. It does not
bother counting tests and generating ratio of passing / failing counts; ***all
tests should always pass***. In that context, the only relevant pieces of
information are:
- a coverage value that should be or approach 100%.
- a human readable list of the tests being run to provide external observers
  with an idea of the scope of the testing implemented.

If you work on a project where some tests are allowed to fail, and you would
like to produce reports of how many tests pass / fail over time, this is outside
the scope of this project. You could however consider using the YAML output from
this tool as the input for some further analysis.

### Sample output

github.com/maargenton/go-cli/pkg/cli

Coverage: 100%

- ❌ TestGetCompletion
  - ❌ Given an configured OptionSet
    - ✅ when calling GetCompletion() with no arguments
      - ✅ then no specific option is being completed
      - ✅ then the first argument is being completed
      - ✅ then available options are listed
    - ✅ when calling GetCompletion() with flag expecting a value
      - ✅ then only yhe option is returned
    - ❌ when calling GetCompletion() with flag and value
      - ❌ then remaining options are listed
        ```
        completion_test.go:71:
            expected:       set(value.Option) == []string{ "--format <value>", "-v" }
            value:          []option.Description{
                            	option.Description{
                            		Option:      "--forfmat <value>",
                            		Description: "communcation format, e.g. 8N1",
                            	},
                            	option.Description{
                            		Option:      "-v",
                            		Description: "display additional information on startup",
                            	},
                            }
            $.Option:       { "--forfmat <value>", "-v" }
            extra values:   "--forfmat <value>"
            missing values: "--format <value>"
        ```
      - ✅ then the first argument is being completed



## Using Custom Template

`go-testreport` can also produce output from a custom template file, based on
standard [Go template](https://pkg.go.dev/text/template). The template is
applied to the list of packages stored in the internal model, generated from
parsing the test results.

```
go-testreport ./... -otemplate/my_template.tmpl=build/my-build-report.yaml
```

### Package

- `{{ .Name }}`: The name of the package
- `{{ Elapsed }}`: The time taken to run the tests
- `{{ Coverage }}`: A percentage value (0..100) of how much of the code is
  exercised by the tests running within the package.
- `{{ Skipped }}`: A boolean value that is `true` if the tests for that package
  have been skipped; usually when there is no test in the package
- `{{ Tests }}`: A list of `Test` objects containing the top level tests of the
  package. Those tests might include nested tests.
- `{{ LeafTests }}`: A list of `Test` objects containing all the leaf tests in
  the package, i.e. all the tests in the test hierarchy that do not contain any
  more nested tests. Those tests are usually the meaningful ones -- the ones
  that evaluate and verify some assertions that might cause a failure.


### Test

- `{{ .Name }}`: The name of the test
- `{{ .FullName }}`: The full name of the test including its name and the name
  of all its `Parent` tests, joined by a ', ' (comma + space).
- `{{ .PartialName <n> }}`: Same as `FullName`, but excluding `n` test names
  from the top of the hierarchy. Using n = 1 skips the top level test name (the
  name of the test function), which is usually formatted differently from the
  rest because it has to be a valid function name, starting with `Test`.
- `{{ .Failure }}`: A boolean value that is `true` if the test was a failure.
- `{{ .Output }}`: A list of strings containing the output of the test, one
  string per line. Common white space prefix and all white spaces suffix are
  stripped.
- `{{ .Parent }}`: A reference back to the parent test if any, or nil for
  package level tests.
- `{{ .Tests }}`: The list of nested tests.
- `{{ .LeafTests }}`: The list of all the tests in the hierarchy of the current
  test that do not contain any more nested tests.

### Custom template function

- `render`: is similar to the built-in `template` function that applies a
  partial template to an object. `render` actually renders the template and
  produces a multi-line string containing the result.
  ```
  {{- render "test" . -}}
  ```
- `indent`: applies the specified number of indentation spaces to a multi-line
  string, shifting the whole content to the right. Note that `indent` does not
  work with the built-in `template` function; it has be used with `render`
  instead.
  ```
  {{- render "test" . | indent 2 -}}
  ```

### Sample template

For reference, the template below is equivalent to the built-in template used to
produce the default markdown output.

```tmpl
# Test report

{{ range . -}}
{{- template "package" . -}}
{{- end }}


{{- define "package" -}}
{{-   if gt (len .Tests) 0 }}

## {{ .Name }}

Coverage: {{ .Coverage }}%
{{      range .Tests -}}
{{-       render "test" . -}}
{{-     end }}
{{    end -}}
{{- end -}}


{{- define "test" }}
- {{ render "outcome" . }} {{ .Name }}
{{-   render "output" . | indent 2 -}}
{{-   range .Tests -}}
{{-     render "test" . | indent 2 -}}
{{-   end -}}
{{- end -}}


{{- define "outcome" }}
{{-   if .Failure -}}❌{{- else -}}✅{{- end -}}
{{- end -}}


{{- define "output" }}
{{-   if gt (len .Output) 0 }}
        ```
{{      range .Output }}
{{        .  }}
{{-     end }}
        ```
{{-   end -}}
{{- end -}}
```
