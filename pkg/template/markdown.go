package template

// MarkdownTemplate contains the definition for the built-in template used to
// generate Markdown report
var MarkdownTemplate = `
{{- if ne Title "" -}}
{{ header 1 }} {{ Title }}
{{- end }}
{{- range . -}}
{{- template "package" . -}}
{{- end }}

{{/* -------------------------------------------------------------------- */}}

{{- define "package" -}}
{{-   if gt (len .Tests) 0 }}

{{ header 2 }} {{ .Name }}

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

// {{- "" -}}

// MarkdownSummaryTemplate contains the definition for the built-in template
// used to generate an alternate Markdown report focussed on summary and
// surfacing failures
var MarkdownSummaryTemplate = `
{{- $any_failure := false -}}
{{- range . -}}
{{-   if gt .Failed 0 -}}{{- $any_failure = true -}}{{- end -}}
{{- end -}}

{{- if ne Title "" -}}
{{ header 1 }} {{ Title }}
{{- end }}

{{ header 2 }} Summary

| Package | Passed | Failed | Coverage |
|-|-|-|-|
{{- range . -}}
{{-   template "package-summary" . -}}
{{- end -}}

{{- if $any_failure }}

{{ header 2 }} Failures
{{ range . -}}
{{-   if gt .Failed 0 -}}
{{-     template "package-failures" . -}}
{{-   end -}}
{{- end -}}
{{- end }}

{{ header 2 }} Details

<details>
<summary>Full report</summary>

{{ range . -}}
{{- template "package" . -}}
{{- end -}}

</details>

{{/* -------------------------------------------------------------------- */}}

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
