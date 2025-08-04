package template

// BuiltinTemplates contains all the template definitions for common fragments
// and supported builtin markdown-based output formats. All templates and
// fragments are named to be called explicitly by `Template.ExecuteTemplate()`.
var BuiltinTemplates = `

{{/* -------------------------------------------------------------------- */}}

{{- define "package-summary" -}}
{{-   if or (gt .Passed 0) (gt .Failed 0) }}
| {{ render "package-outcome" . }} {{ .Name }} | {{ .Passed }} | {{ .Failed }} | {{ .Coverage }}% |
{{-    end -}}
{{- end -}}

{{- define "package-outcome" }}
{{-   if gt .Failed 0 -}}❌{{- else -}}✅{{- end -}}
{{- end -}}

{{- define "package-build-error" }}
{{   if .BuildError -}}
❌ Build Errors:
{{      codeblock }}
{{      .BuildError }}
{{      codeblock }}
{{    end -}}
{{- end -}}

{{- define "package-coverage" }}
{{    if gt .Coverage 0.0 -}}
Coverage: {{ .Coverage }}%
{{    end -}}
{{- end -}}

{{/* -------------------------------------------------------------------- */}}

{{- define "package-failures" }}
{{ header 3 }} {{ .Name }}
{{/* ---newline--- */}}

{{- render "package-build-error" . -}}
{{- render "package-coverage" . -}}

{{-    range .Tests -}}
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
{{-   if or (gt .Passed 0) (gt .Failed 0) }}
{{ header 3 }} {{ .Name }}
{{/* ---newline--- */}}

{{- render "package-build-error" . -}}
{{- render "package-coverage" . -}}

{{     range .Tests -}}
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
|---------|--------|--------|----------|
{{-   range . -}}
{{-     render "package-summary" . -}}
{{-   end}}

{{    header 2 }} Full report
{{    range . -}}
{{-     render "package-details" . -}}
{{-   end -}}
{{- end -}}

{{- define "md" }}
{{-   render "markdown" . -}}
{{- end -}}


{{/* -------------------------------------------------------------------- */}}

{{- define "mds" }}
{{-   if ne Title "" -}}
{{      header 1 }} {{ Title }}

{{    end -}}
{{    header 2 }} Packages

| Package | Passed | Failed | Coverage |
|---------|--------|--------|----------|
{{-   range . -}}
{{-     render "package-summary" . -}}
{{-   end }}
{{- end -}}


{{/* -------------------------------------------------------------------- */}}

{{- define "mdsf" }}
{{-   render "mds" . -}}

{{-   $any_failure := false -}}
{{-   range . -}}
{{-     if gt .Failed 0 -}}{{- $any_failure = true -}}{{- end -}}
{{-   end -}}

{{-   if $any_failure }}

{{      header 2 }} Failures
{{      range . -}}
{{-       if gt .Failed 0 -}}
{{-         render "package-failures" . -}}
{{-       end -}}
{{-     end -}}
{{-   end }}
{{- end -}}


{{/* -------------------------------------------------------------------- */}}

{{- define "mdsfd" }}
{{-   render "mdsf" . }}

{{    header 2 }} Full report

<details>
<summary>Expand</summary>

{{     range . -}}
{{-      render "package-details" . -}}
{{-     end }}
</details>
{{- end -}}


{{- define "markdown-summary" }}
{{-   render "mdsfd" . -}}
{{- end -}}
`
