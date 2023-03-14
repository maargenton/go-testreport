package template

var Builtin = map[string][]string{
	"md":       {MarkdownShared, MarkdownTemplate},
	"markdown": {MarkdownShared, MarkdownTemplate},

	"mds":              {MarkdownShared, MarkdownSTemplate},
	"mdsf":             {MarkdownShared, MarkdownSFTemplate},
	"mdsfd":            {MarkdownShared, MarkdownSFDTemplate},
	"markdown-summary": {MarkdownShared, MarkdownSFDTemplate},
}

// MarkdownShared contains common built-in definitions shared by one or more
// markdown output templates.
var MarkdownShared = `
{{- define "package-summary" -}}
{{-   if gt (len .Tests) 0 }}
| {{ render "package-outcome" . }} {{ .Name }} | {{ .Passed }} | {{ .Failed }} | {{ .Coverage }}% |
{{-    end -}}
{{- end -}}

{{- define "package-outcome" }}
{{-   if gt .Failed 0 -}}❌{{- else -}}✅{{- end -}}
{{- end -}}

{{/* -------------------------------------------------------------------- */}}

{{- define "package-failures" }}
{{ header 3 }} {{ .Name }}
{{    range .Tests -}}
{{-     if .Failure -}}
{{-       render "test-failure" . -}}
{{-     end -}}
{{-   end }}
{{- end -}}

{{- define "test-failure" }}
- {{  render "outcome" . }} {{ .Name }}
{{-   render "output" . | indent 2 -}}
{{-   range .Tests -}}
{{-     if .Failure -}}
{{-       render "test-failure" . | indent 2 -}}
{{-     end -}}
{{-   end -}}
{{- end -}}

{{/* -------------------------------------------------------------------- */}}

{{- define "package" -}}
{{-   if gt (len .Tests) 0 }}
{{ header 3 }} {{ .Name }}

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
{{      codeblock -}}
{{      range .Output }}
{{        .  }}
{{-     end }}
{{      codeblock }}
{{-   end -}}
{{- end -}}
`

// MarkdownTemplate contains the definition for the built-in template used to
// generate Markdown report
var MarkdownTemplate = `
{{- if ne Title "" -}}
{{ header 1 }} {{ Title }}

{{ end -}}
{{ header 2 }} Packages

| Package | Passed | Failed | Coverage |
|-|-|-|-|
{{- range . -}}
{{-   template "package-summary" . -}}
{{- end}}

{{ header 2 }} Full report
{{ range . -}}
{{- template "package" . -}}
{{- end -}}
`

// MarkdownSFDTemplate contains the definition for the built-in template
// used to generate an alternate Markdown report focussed on summary and
// surfacing failures
var MarkdownSFDTemplate = `
{{- if ne Title "" -}}
{{ header 1 }} {{ Title }}

{{ end -}}
{{ header 2 }} Packages

| Package | Passed | Failed | Coverage |
|-|-|-|-|
{{- range . -}}
{{-   template "package-summary" . -}}
{{- end -}}

{{- $any_failure := false -}}
{{- range . -}}
{{-   if gt .Failed 0 -}}{{- $any_failure = true -}}{{- end -}}
{{- end -}}
{{- if $any_failure }}

{{ header 2 }} Failures
{{ range . -}}
{{-   if gt .Failed 0 -}}
{{-     template "package-failures" . -}}
{{-   end -}}
{{- end -}}
{{- end }}

{{ header 2 }} Full report

<details>
<summary>Expand</summary>

{{ range . -}}
{{- template "package" . -}}
{{- end -}}

</details>
`

// MarkdownSTemplate follows the same format as the SFD (summary, failures,
// details) template, but includes only the package summary.
var MarkdownSTemplate = `
{{- if ne Title "" -}}
{{ header 1 }} {{ Title }}

{{ end -}}
{{ header 2 }} Packages

| Package | Passed | Failed | Coverage |
|-|-|-|-|
{{- range . -}}
{{-   template "package-summary" . -}}
{{- end }}
`

// MarkdownSFTemplate follows the same format as the SFD (summary, failures,
// details) template, but includes only the package summary and test failures if
// any.
var MarkdownSFTemplate = `
{{- if ne Title "" -}}
{{ header 1 }} {{ Title }}

{{ end -}}
{{ header 2 }} Packages

| Package | Passed | Failed | Coverage |
|-|-|-|-|
{{- range . -}}
{{-   template "package-summary" . -}}
{{- end -}}

{{- $any_failure := false -}}
{{- range . -}}
{{-   if gt .Failed 0 -}}{{- $any_failure = true -}}{{- end -}}
{{- end -}}
{{- if $any_failure }}

{{ header 2 }} Failures
{{ range . -}}
{{-   if gt .Failed 0 -}}
{{-     template "package-failures" . -}}
{{-   end -}}
{{- end -}}
{{- end }}
`
