package template

// BuiltinTemplates contains all the template definitions for common fragments
// and supported builtin markdown-based output formats. All templates and
// fragments are named to be called explicitly by `Template.ExecuteTemplate()`.
var BuiltinTemplates = `

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
- {{  render "test-outcome" . }} {{ .Name }}
{{-   render "test-output" . | indent 2 -}}
{{-   range .Tests -}}
{{-     if .Failure -}}
{{-       render "test-failure" . | indent 2 -}}
{{-     end -}}
{{-   end -}}
{{- end -}}

{{/* -------------------------------------------------------------------- */}}

{{- define "package-details" -}}
{{-   if gt (len .Tests) 0 }}
{{ header 3 }} {{ .Name }}

Coverage: {{ .Coverage }}%
{{      range .Tests -}}
{{-       render "test-details" . -}}
{{-     end }}
{{    end -}}
{{- end -}}

{{- define "test-details" }}
- {{ render "test-outcome" . }} {{ .Name }}
{{-   render "test-output" . | indent 2 -}}
{{-   range .Tests -}}
{{-     render "test-details" . | indent 2 -}}
{{-   end -}}
{{- end -}}

{{- define "test-outcome" }}
{{-   if .Failure -}}❌{{- else -}}✅{{- end -}}
{{- end -}}

{{- define "test-output" }}
{{-   if gt (len .Output) 0 }}
{{      codeblock -}}
{{      range .Output }}
{{        .  }}
{{-     end }}
{{      codeblock }}
{{-   end -}}
{{- end -}}

{{/* -------------------------------------------------------------------- */}}

{{- define "markdown" }}
{{-   if ne Title "" -}}
{{      header 1 }} {{ Title }}

{{    end -}}
{{    header 2 }} Packages

| Package | Passed | Failed | Coverage |
|-|-|-|-|
{{-   range . -}}
{{-     template "package-summary" . -}}
{{-   end}}

{{    header 2 }} Full report
{{    range . -}}
{{-     template "package-details" . -}}
{{-   end -}}
{{- end -}}

{{- define "md" }}
{{-   template "markdown" . -}}
{{- end -}}


{{/* -------------------------------------------------------------------- */}}

{{- define "mds" }}
{{-   if ne Title "" -}}
{{      header 1 }} {{ Title }}

{{    end -}}
{{    header 2 }} Packages

| Package | Passed | Failed | Coverage |
|-|-|-|-|
{{-   range . -}}
{{-     template "package-summary" . -}}
{{-   end }}
{{- end -}}


{{/* -------------------------------------------------------------------- */}}

{{- define "mdsf" }}
{{-   template "mds" . -}}

{{-   $any_failure := false -}}
{{-   range . -}}
{{-     if gt .Failed 0 -}}{{- $any_failure = true -}}{{- end -}}
{{-   end -}}

{{-   if $any_failure }}

{{      header 2 }} Failures
{{      range . -}}
{{-       if gt .Failed 0 -}}
{{-         template "package-failures" . -}}
{{-       end -}}
{{-     end -}}
{{-   end }}
{{- end -}}


{{/* -------------------------------------------------------------------- */}}

{{- define "mdsfd" }}
{{-   template "mdsf" . -}}

{{    header 2 }} Full report

<details>
<summary>Expand</summary>

{{     range . -}}
{{-      template "package-details" . -}}
{{-     end -}}
</details>
{{- end -}}


{{- define "markdown-summary" }}
{{-   template "mdsfd" . -}}
{{- end -}}
`
