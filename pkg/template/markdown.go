package template

// MarkdownTemplate contains the definition for the built-in template used to
// generate Markdown report
var MarkdownTemplate = `
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
{{      codeblock -}}
{{      range .Output }}
{{        .  }}
{{-     end }}
{{      codeblock }}
{{-   end -}}
{{- end -}}
`

// MarkdownSummaryTemplate contains the definition for the built-in template
// used to generate an alternate Markdown report focussed on summary and
// surfacing failures
var MarkdownSummaryTemplate = `
# Test report

## Summary

{{ define "package-summary" -}}
{{-   if gt (len .Tests) 0 }}
| {{ render "package-outcome" . }} {{ .Name }} | {{ .Passed }} | {{ .Failed }} | {{ .Coverage }}% |
{{-    end -}}
{{- end -}}

{{- define "package-outcome" }}
{{-   if gt .Failed 0 -}}❌{{- else -}}✅{{- end -}}
{{- end -}}

| Package | Passed | Failed | Coverage |
|-|-|-|-|
{{- range . -}}
{{- template "package-summary" . -}}
{{- end }}

## Failures

{{ define "package-failures" -}}
### {{ .Name }}
{{    range .Tests -}}
{{-     if .Failure -}}
{{-       render "test" . -}}
{{-     end -}}
{{-   end }}
{{- end -}}

{{ range . -}}
{{-   if gt .Failed 0 -}}
{{-     template "package-failures" . -}}
{{-   end -}}
{{- end }}


## Details

<details>
<summary>Full report</summary>

{{ range . -}}
{{- template "package" . -}}
{{- end }}

</details>


{{- define "package" -}}
{{-   if gt (len .Tests) 0 }}
### {{ .Name }}

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
